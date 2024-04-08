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

func processRequest(req *http.Request, dir *string) error {
	url := req.GetUrl()
	conn := req.GetConnection()
	var err error
	res := http.NewResponse(conn)

	if url == "/" {
		res.Ok()
		return nil
	} else if strings.Contains(url, "/echo/") {
		content := (url)[6:]
		res.WriteHeader("Content-Type", "text/plain")
		res.SendWithBody(http.StatusOk, &content)
	} else if url == "/user-agent" {
		content := req.Header["User-Agent"]
		res.WriteHeader("Content-Type", "text/plain")
		res.SendWithBody(http.StatusOk, &content)
	} else if strings.Contains(url, "/files/") {
		fileName := url[7:]
		if dir != nil {
			file := fmt.Sprintf("%s%s", *dir, fileName)
			if req.Method == http.GET {
				data, readFileError := os.ReadFile(file)
				if readFileError != nil {
					res := http.NewResponse(conn)
					res.NotFound()
				} else {
					content := string(data[:])
					res.WriteHeader("Content-Type", "application/octet-stream")
					res.SendWithBody(http.StatusOk, &content)
				}
			} else if req.Method == http.POST {
				if req.Body == nil {
					res.BadRequest()
					return nil
				}
				f, err := os.Create(file)
				if err != nil {
					res.ServerError()
					return err
				}
				_, err = f.WriteString(*req.Body)
				if err != nil {
					res.ServerError()
					return err
				}
				res.Status = http.StatusCreated
				res.Send()
			}

		} else {
			res.NotFound()
		}

	} else {
		res.NotFound()
	}
	return err
}
