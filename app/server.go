package main

import (
	"flag"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	dir := flag.String("directory", "", "The name of the directory")
	flag.Parse()
	http.Listen("4221", dir)
}
