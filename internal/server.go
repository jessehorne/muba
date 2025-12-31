package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	Address string
	Port    string

	Users map[string]*User
	Teams map[string]*Team

	listener net.Listener
	closer   chan struct{}

	welcomeMessage string
}

func NewServer(address, port string) (*Server, error) {
	s := &Server{
		Address: address,
		Port:    port,
		Users:   make(map[string]*User),
		Teams: map[string]*Team{
			"red":  NewTeam("red"),
			"blue": NewTeam("blue"),
		},
		closer:         make(chan struct{}),
		welcomeMessage: string(readFromFile("./data/welcome.utf8")),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port))
	if err != nil {
		return nil, err
	}

	s.listener = listener
	return s, nil
}

func (s *Server) Run() error {
	log.Println("Starting server...")
	for {
		select {
		case <-s.closer:
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil && err != io.EOF {
				log.Println("[ERROR]", err)
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("[INFO] New connection", conn.RemoteAddr())
	s.handleUser(conn)
	log.Println("[INFO] Closing connection", conn.RemoteAddr())
}

func (s *Server) handleRawMessage(conn net.Conn, msg []byte) {
	splitted := strings.Split(string(msg), " ")
	if splitted[0] == "exit" {
		s.handleUserLeave(conn)
	}
}

func (s *Server) handleUser(conn net.Conn) {
	newUser := NewUser(conn, s)
	s.Users[conn.RemoteAddr().String()] = newUser
	err := s.Users[conn.RemoteAddr().String()].Handle()
	if err != nil {
		log.Println("[ERROR]", conn.RemoteAddr(), err)
	}
}

func (s *Server) handleUserLeave(conn net.Conn) {
	_, ok := s.Users[conn.RemoteAddr().String()]
	if ok {
		delete(s.Users, conn.RemoteAddr().String())
	}
}

func (s *Server) Close() error {
	s.closer <- struct{}{}
	return s.listener.Close()
}
