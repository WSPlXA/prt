package network

import (
	"fmt"
	"sync"
)

var _ Transport = (*LocalTransport)(nil)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
	peerCh    chan NetAddr
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
		peerCh:    make(chan NetAddr, 1024),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.peers[tr.Addr()] = tr.(*LocalTransport)

	select {
	case t.peerCh <- tr.Addr():
	default:
	}

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, Payload []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send messages to %s", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: Payload,
	}

	return nil

}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) PeerCh() <-chan NetAddr {
	return t.peerCh
}
