package main

import (
	"fmt"
	"log"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const (
	httpProtocol = "http"
	httpVersion  = "1.1"
	httpStatusOK = "200 OK"
	crlf         = "\r\n\r\n"
)

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
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		log.Println("Error accepting connection: ", err.Error())

		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection: ", err.Error())
			return
		}

	}()
	okResponse := httpProtocol + "/" + httpVersion + " " + httpStatusOK + " " + httpVersion + crlf
	_, err = conn.Write([]byte(okResponse))
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		log.Println("Error writing response: ", err.Error())
		os.Exit(1)

	}
}
