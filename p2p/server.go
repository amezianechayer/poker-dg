package p2p

import (
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Server struct {
	ServerConfig

	handler   Handler
	transport *TCPTransport
	mu        sync.RWMutex
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer := s.addPeer
	tr.DelPeer := S.addPeer


	return s
}

func (s *Server) Start() {
	go s.loop()

	fmt.Printf("serveur de jeu fonctionnant sur le port %s\n", s.ListenAddr)

	s.transport.ListenAndAccept()
}

// TODO(@amezianechayer): Maintenant on a un code redondant en inscrivant de nouveaux pairs au réseau du jeu
// peut-être construire un nouveau protocole d'homologue et de prise de contact après avoir enregistré une connexion simple?
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(s.Version))
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("nouveau joueur connecté")
			fmt.Printf("joueur déconnecté %s\n", peer.conn.RemoteAddr())
			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("nouveau joueur connecté")
			
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("Nouveau joueur connecté %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
