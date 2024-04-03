package main

import (
	"fmt"
	"log"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const (
	httpProtocol = "HTTP"
	httpVersion  = "1.1"
	httpStatusOK = "200 OK"
	crlf         = "\r\n\r\n"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	// Uncomment this block to pass the first stage
	buffer := make([]byte, 1024)
	_, err := c.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	c.Write([]byte(httpProtocol + "/" + httpVersion + " " + httpStatusOK + " " + httpVersion + crlf))
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Failed to bind to port 4221")
		os.Exit(1)

	}
	// Uncomment this block to pass the first stage
	defer func() {
		err := l.Close()
		if err != nil {
			log.Println("Error closing listener: ", err.Error())
			return
		}

	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}

}
