package http

import (
	"fmt"
	"net"
	"os"
)

type server struct {
	ctx      ServerContext
	listener net.Listener
}

func NewServer(l net.Listener) *server {
	return &server{
		listener: l,
	}
}

func CreateConnection(network string, address string) *server {

	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	server := NewServer(l)
	return server
}

func (s *server) GetRequest(conn net.Conn) *Request {
	request := NewRequest(conn)
	return request
}

func (s *server) GetContext() ServerContext {
	return s.ctx
}

func (s *server) AcceptConnection() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	return conn
}
