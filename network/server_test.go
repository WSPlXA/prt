package network

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerOnRPC(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")
	trA.Connect(trB)
	trB.Connect(trA)

	var mu sync.Mutex
	var received []RPC

	opts := ServerOpts{
		Transports: []Transport{trA},
		OnRPC: func(rpc RPC) error {
			mu.Lock()
			received = append(received, rpc)
			mu.Unlock()
			return nil
		},
	}

	s := NewServer(opts)
	go s.Start()
	defer s.Stop()

	// give initTransport goroutines time to start
	time.Sleep(50 * time.Millisecond)

	err := trB.SendMessage(trA.Addr(), []byte("hello from B"))
	require.Nil(t, err)

	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(received) == 1
	}, time.Second, 10*time.Millisecond)

	mu.Lock()
	rpc := received[0]
	mu.Unlock()

	assert.Equal(t, []byte("hello from B"), rpc.Payload)
	assert.Equal(t, NetAddr("B"), rpc.From)
}

func TestServerOnPeer(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")

	var peerAddr NetAddr
	var mu sync.Mutex

	opts := ServerOpts{
		Transports: []Transport{trA},
		OnPeer: func(addr NetAddr) {
			mu.Lock()
			peerAddr = addr
			mu.Unlock()
		},
	}

	s := NewServer(opts)
	go s.Start()
	defer s.Stop()

	time.Sleep(50 * time.Millisecond)

	trA.Connect(trB)

	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return peerAddr == NetAddr("B")
	}, time.Second, 10*time.Millisecond)
}

func TestServerBroadcast(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")
	trC := NewLocalTransport("C")

	trA.Connect(trB)
	trA.Connect(trC)

	opts := ServerOpts{
		Transports: []Transport{trA},
	}

	s := NewServer(opts)
	go s.Start()
	defer s.Stop()

	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 2, s.PeerCount())

	err := s.Broadcast([]byte("broadcast msg"))
	require.Nil(t, err)

	rpcB := <-trB.Consume()
	assert.Equal(t, []byte("broadcast msg"), rpcB.Payload)
	assert.Equal(t, NetAddr("A"), rpcB.From)

	rpcC := <-trC.Consume()
	assert.Equal(t, []byte("broadcast msg"), rpcC.Payload)
	assert.Equal(t, NetAddr("A"), rpcC.From)
}

func TestServerSend(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")

	trA.Connect(trB)

	opts := ServerOpts{
		Transports: []Transport{trA},
	}

	s := NewServer(opts)
	go s.Start()
	defer s.Stop()

	time.Sleep(50 * time.Millisecond)

	err := s.Send(NetAddr("B"), []byte("direct msg"))
	require.Nil(t, err)

	rpc := <-trB.Consume()
	assert.Equal(t, []byte("direct msg"), rpc.Payload)
	assert.Equal(t, NetAddr("A"), rpc.From)
}

func TestServerStop(t *testing.T) {
	tr := NewLocalTransport("A")

	opts := ServerOpts{
		Transports: []Transport{tr},
	}

	s := NewServer(opts)

	done := make(chan struct{})
	go func() {
		s.Start()
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not stop")
	}
}
