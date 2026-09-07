package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"redis-go/internal/command"
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

type Server struct {
	address  string
	store    *storage.Store
	commands *command.Registry
}

func New(address string) *Server {
	return &Server{
		address:  address,
		store:    storage.NewStore(),
		commands: command.NewRegistry(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	parser := protocol.NewParser(reader)
	response := protocol.NewWriter(writer)

	for {
		args, err := parser.ReadCommand()
		if err != nil {
			if err != io.EOF {
				log.Printf("connection error: %v", err)
			}

			return
		}

		if err := s.commands.Execute(args, s.store, response); err != nil {
			log.Printf("command error: %v", err)
			return
		}
	}
}
