package main

import (
	"flag"
)

func main() {
	dir := flag.String("directory", "", "The name of the directory")
	flag.Parse()
	http.Listen("4221", dir)
}
