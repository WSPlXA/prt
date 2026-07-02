package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPTransportConnect(t *testing.T) {
	tra := NewTCPTransport("127.0.0.1:0")
	defer tra.Close()
	trb := NewTCPTransport("127.0.0.1:0")
	defer trb.Close()

	tra.addr = NetAddr(tra.listener.Addr().String())
	trb.addr = NetAddr(trb.listener.Addr().String())

	err := tra.Connect(trb)
	require.Nil(t, err)

	<-trb.peerCh

	err = trb.Connect(tra)
	require.Nil(t, err)

	assert.Contains(t, tra.peers, trb.Addr())
	assert.Contains(t, trb.peers, tra.Addr())
}

func TestTCPTransportSendMessage(t *testing.T) {
	tra := NewTCPTransport("127.0.0.1:0")
	defer tra.Close()
	trb := NewTCPTransport("127.0.0.1:0")
	defer trb.Close()

	tra.addr = NetAddr(tra.listener.Addr().String())
	trb.addr = NetAddr(trb.listener.Addr().String())

	err := tra.Connect(trb)
	require.Nil(t, err)

	<-trb.peerCh

	err = trb.Connect(tra)
	require.Nil(t, err)

	msg := []byte("Hello TCP")
	err = tra.SendMessage(trb.Addr(), msg)
	require.Nil(t, err)

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestTCPTransportSendMessageReverse(t *testing.T) {
	tra := NewTCPTransport("127.0.0.1:0")
	defer tra.Close()
	trb := NewTCPTransport("127.0.0.1:0")
	defer trb.Close()

	tra.addr = NetAddr(tra.listener.Addr().String())
	trb.addr = NetAddr(trb.listener.Addr().String())

	err := tra.Connect(trb)
	require.Nil(t, err)

	<-trb.peerCh

	err = trb.Connect(tra)
	require.Nil(t, err)

	err = tra.SendMessage(trb.Addr(), []byte("A to B"))
	require.Nil(t, err)
	err = trb.SendMessage(tra.Addr(), []byte("B to A"))
	require.Nil(t, err)

	rpcA := <-trb.Consume()
	assert.Equal(t, []byte("A to B"), rpcA.Payload)
	assert.Equal(t, tra.Addr(), rpcA.From)

	rpcB := <-tra.Consume()
	assert.Equal(t, []byte("B to A"), rpcB.Payload)
	assert.Equal(t, trb.Addr(), rpcB.From)
}
