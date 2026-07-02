package main

import (
	"fmt"
	"prt/network"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("Hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
		OnRPC: func(rpc network.RPC) error {
			fmt.Printf("[%s] => %s: %s\n", rpc.From, "LOCAL", string(rpc.Payload))
			return nil
		},
		OnPeer: func(addr network.NetAddr) {
			fmt.Printf("new peer connected: %s\n", addr)
		},
	}

	s := network.NewServer(opts)
	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()

	s.Start()
}
