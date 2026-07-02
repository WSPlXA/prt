package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

var _ Transport = (*TCPTransport)(nil)

type TCPPeer struct {
	conn net.Conn
	enc  *gob.Encoder
	mu   sync.Mutex
}

type TCPTransport struct {
	addr      NetAddr
	listener  net.Listener
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*TCPPeer
	peerCh    chan NetAddr
	quitCh    chan struct{}
}

func NewTCPTransport(addr NetAddr) *TCPTransport {
	listener, err := net.Listen("tcp", string(addr))
	if err != nil {
		panic(err)
	}
	t := &TCPTransport{
		addr:      addr,
		listener:  listener,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*TCPPeer),
		peerCh:    make(chan NetAddr, 1024),
		quitCh:    make(chan struct{}),
	}
	go t.acceptLoop()
	return t
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.quitCh:
				return
			default:
			}
			continue
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	dec := gob.NewDecoder(conn)

	var peerAddr NetAddr
	if err := dec.Decode(&peerAddr); err != nil {
		conn.Close()
		return
	}

	select {
	case t.peerCh <- peerAddr:
	default:
	}

	for {
		var rpc RPC
		if err := dec.Decode(&rpc); err != nil {
			break
		}
		t.consumeCh <- rpc
	}
	conn.Close()
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *TCPTransport) Connect(tr Transport) error {
	t.lock.RLock()
	_, exists := t.peers[tr.Addr()]
	t.lock.RUnlock()
	if exists {
		return nil
	}

	conn, err := net.Dial("tcp", string(tr.Addr()))
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(t.addr); err != nil {
		conn.Close()
		return err
	}

	peer := &TCPPeer{
		conn: conn,
		enc:  enc,
	}

	t.lock.Lock()
	t.peers[tr.Addr()] = peer
	t.lock.Unlock()

	return nil
}

func (t *TCPTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	peer, ok := t.peers[to]
	t.lock.RUnlock()

	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.mu.Lock()
	defer peer.mu.Unlock()

	return peer.enc.Encode(RPC{
		From:    t.addr,
		Payload: payload,
	})
}

func (t *TCPTransport) Addr() NetAddr {
	return t.addr
}

func (t *TCPTransport) PeerCh() <-chan NetAddr {
	return t.peerCh
}

func (t *TCPTransport) Close() error {
	close(t.quitCh)
	return t.listener.Close()
}
