package network

import (
	"fmt"
	"sync"
)

type ServerOpts struct {
	Transports []Transport
	OnRPC      func(RPC) error
	OnPeer     func(NetAddr)
}

type Server struct {
	ServerOpts

	rpcCh  chan RPC
	quitCh chan struct{}

	peers    map[NetAddr]Transport
	peerLock sync.RWMutex
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC, 1024),
		quitCh:     make(chan struct{}),
		peers:      make(map[NetAddr]Transport),
	}
}

func (s *Server) Start() {
	s.initTransport()

	for {
		select {
		case rpc := <-s.rpcCh:
			if s.OnRPC != nil {
				s.OnRPC(rpc)
			}
		case <-s.quitCh:
			return
		}
	}
}

func (s *Server) Stop() {
	close(s.quitCh)
}

func (s *Server) Broadcast(payload []byte) error {
	s.peerLock.RLock()
	defer s.peerLock.RUnlock()

	for addr, tr := range s.peers {
		if err := tr.SendMessage(addr, payload); err != nil {
			return fmt.Errorf("broadcast to %s: %w", addr, err)
		}
	}
	return nil
}

func (s *Server) Send(to NetAddr, payload []byte) error {
	s.peerLock.RLock()
	tr, ok := s.peers[to]
	s.peerLock.RUnlock()

	if !ok {
		return fmt.Errorf("peer %s not found", to)
	}

	return tr.SendMessage(to, payload)
}

func (s *Server) PeerCount() int {
	s.peerLock.RLock()
	defer s.peerLock.RUnlock()
	return len(s.peers)
}

func (s *Server) initTransport() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)

		go func(tr Transport) {
			for peerAddr := range tr.PeerCh() {
				s.peerLock.Lock()
				s.peers[peerAddr] = tr
				s.peerLock.Unlock()

				if s.OnPeer != nil {
					s.OnPeer(peerAddr)
				}
			}
		}(tr)
	}
}
