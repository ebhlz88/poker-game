package main

import (
	"github.com/ebhlz88/poker-game/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "0.1",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(cfg)
	server.Start()
	// d := deck.New()
	// fmt.Print(d)
}
