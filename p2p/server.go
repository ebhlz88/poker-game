package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Message struct {
	From    net.Addr
	Payload io.Reader
}

type Server struct {
	mu sync.RWMutex

	handler  Handler
	listener net.Listener
	peers    map[net.Addr]*Peer
	ServerConfig
	addPeer chan *Peer
	delPeer chan *Peer
	msgChan chan *Message
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		handler:      &DefaultHandler{},
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
		delPeer:      make(chan *Peer),
	}
}

func (s *Server) Start() {
	go s.loop()
	err := s.Listen()
	if err != nil {
		panic(err)
	}
	fmt.Printf("New Server Started %s\n", s.ListenAddr)

	s.AcceptLoop()
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		p := &Peer{
			conn: conn,
		}
		s.addPeer <- p
		p.Send([]byte(s.Version))
		go s.handleConn(p)

	}
}

func (s *Server) handleConn(p *Peer) {
	defer func() {
		s.delPeer <- p
	}()
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}
		s.msgChan <- &Message{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}

		fmt.Println(string(buf[:n]))
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.listener = ln
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new Peer Connected %s\n", peer.conn.RemoteAddr())
		case msg := <-s.msgChan:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic("not working")
			}
		case peer := <-s.delPeer:
			delete(s.peers, peer.conn.RemoteAddr())
			fmt.Printf("Player disconnected %s", peer.conn.RemoteAddr())

		}
	}
}
