package main

import (
	"fmt"
	"time"

	"github.com/ebhlz88/poker-game/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "0.1",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(cfg)
	go server.Start()
	time.Sleep(time.Second * 1)

	remoteCfg := p2p.ServerConfig{
		Version:     "0.1",
		ListenAddr:  ":4000",
		GameVariant: p2p.Holdem,
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		fmt.Print(err)
	}
	// d := deck.New()
	// fmt.Print(d)
	select {}
}
