package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	server := http.CreateConnection("tcp", "0.0.0.0:4221")
	dir := flag.String("directory", "", "The name of the directory")
	flag.Parse()
	for {
		conn := server.AcceptConnection()
		go func() {
			defer conn.Close()
			request := server.GetRequest(conn)
			fmt.Println("URL is", request.GetUrl())
			err := processRequest(request, dir)
			if err != nil {
				fmt.Println("Error while writing: ", err.Error())
				os.Exit(1)
			}
		}()
	}
}

func processRequest(request *http.Request, dir *string) error {
	url := request.GetUrl()
	conn := request.GetConnection()
	var err error
	if url == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(url, "/echo/") {
		content := (url)[6:]
		fmt.Println("Content is ", content)
		_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%+v\r\n\r\n%v\r\n", len(content), content)))
	} else if url == "/user-agent" {
		content := request.Header["User-Agent"]
		_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%+v\r\n\r\n%v\r\n", len(content), content)))
	} else if strings.Contains(url, "/files/") {
		fileName := url[7:]
		fmt.Println("file name is", fileName)
		if dir != nil {
			file := fmt.Sprintf("%s%s", *dir, fileName)
			data, readFileError := os.ReadFile(file)
			if readFileError != nil {
				_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n\r\n"))
			} else {
				content := string(data[:])
				_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length:%+v\r\n\r\n%v\r\n", len(content), content)))
			}
		} else {
			_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	return err
}
