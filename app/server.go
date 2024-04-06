package main

import (
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	// accept concurrent connections
	for {
		go func() {
			serv := http.Listen("4221")
			req := serv.ReadRequest()
			if req.Path == "/" {
				serv.Respond(http.StatusOk)
			} else if strings.HasPrefix(req.Path, "/echo") {
				content := strings.ReplaceAll(req.Path, "/echo/", "")
				serv.RespondWithContent(http.StatusOk, &content)
			} else if req.Path == "/user-agent" {
				content := req.Headers["User-Agent"]
				serv.RespondWithContent(http.StatusOk, &content)
			} else {
				serv.Respond(http.StatusNotFound)
			}
		}()
	}
}
