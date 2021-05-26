package main

/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/piloth/core"
	"os"
)

func main() {
	p, err := core.NewPilot()
	if err != nil {
		fmt.Printf("cannot start pilot: %s", err)
		os.Exit(1)
	}
	p.Start()
}
