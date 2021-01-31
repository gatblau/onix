package main

import "github.com/gatblau/onix/artisan/artreg/server"

func main() {
	s := new(server.Server)
	s.Serve()
}
