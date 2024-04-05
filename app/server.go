package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	serv := http.Listen("4221")
	req := serv.ReadRequest()
	if req.Path == "/" {
		serv.Respond(http.StatusOk)
	} else {

		serv.Respond(http.StatusNotFound)
	}
}
