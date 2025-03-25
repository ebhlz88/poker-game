package p2p

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	Holdem GameVariant = iota
	Other
)

func (g GameVariant) String() string {
	switch g {
	case Holdem:
		return "Holdem"
	case Other:
		return "Other"
	default:
		return "Unknown"
	}
}

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	transport *TCPTransport
	handler   Handler
	peers     map[net.Addr]*Peer
	ServerConfig
	addPeer chan *Peer
	delPeer chan *Peer
	msgChan chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	gob.Register(&HandShake{}) // Fix: Register HandShake as a pointer
	gob.Register(GameVariant(0))

	s := &Server{
		ServerConfig: cfg,
		handler:      &DefaultHandler{},
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
		delPeer:      make(chan *Peer),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr
	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer
	return s
}

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"version":      s.Version,
		"port":         s.ListenAddr,
		"Game variant": s.GameVariant,
	}).Info("started new server")

	s.transport.ListenAndAccept()
}

func (s *Server) SendHandShake(p *Peer) error {
	hs := &HandShake{
		Version:     s.Version,
		GameVariant: s.GameVariant,
	}

	encoder := gob.NewEncoder(p.conn) // Fix: Use encoder directly on connection
	if err := encoder.Encode(hs); err != nil {
		return fmt.Errorf("gob encode error: %v", err)
	}

	return nil
}

func (s *Server) Connect(Addr string) error {
	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		return err
	}
	p := &Peer{conn: conn}
	s.addPeer <- p

	// Fix: Send proper handshake
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			s.SendHandShake(peer)
			if err := s.HandShake(peer); err != nil {
				logrus.Errorf("Handshake failed %s", err)
				continue
			}
			go peer.ReadLoop(s.msgChan)
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("handshake successful, New Player joined")
			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				panic("not working")
			}
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("player left")
			delete(s.peers, peer.conn.RemoteAddr())
			fmt.Printf("Player disconnected %s", peer.conn.RemoteAddr())
		}
	}
}

type HandShake struct {
	Version     string
	GameVariant GameVariant
}

func (s *Server) HandShake(p *Peer) error {
	hs := &HandShake{}
	decoder := gob.NewDecoder(p.conn) // Fix: Use decoder directly on connection

	if err := decoder.Decode(hs); err != nil {
		return fmt.Errorf("gob decode error: %v", err)
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("gamevariant does not match %s", hs.GameVariant)
	}
	if s.Version != hs.Version {
		return fmt.Errorf("invalid version %s", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"Version": hs.Version,
		"Variant": hs.GameVariant,
		"peer":    p.conn.RemoteAddr(),
	}).Info("received handshake")

	return nil
}

func (s *Server) handleMessage(m *Message) error {
	fmt.Printf("%+v\n", m)
	return nil
}
