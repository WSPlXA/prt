package main

//server
//Transport => tcp, udp
//Block
//Tx

import (
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
		Transport: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()

}
