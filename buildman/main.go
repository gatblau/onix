package main

import "github.com/gatblau/onix/buildman/server"

func main() {
	s := server.NewServer()
	s.Serve()
}
