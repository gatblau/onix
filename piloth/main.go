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
	ctl "github.com/gatblau/onix/pilotctl/types"
	"github.com/gatblau/onix/piloth/core"
	"os"
)

func main() {
	// collects host information
	hostInfo, err := ctl.NewHostInfo()
	if err != nil {
		panic(err)
	}
	// check for and execute any command line arguments
	if handleCommands(hostInfo) {
		os.Exit(0)
	}
	// creates pilot instance
	p, err := core.NewPilot(hostInfo)
	if err != nil {
		fmt.Printf("cannot start pilot: %s\n", err)
		os.Exit(1)
	}
	// start the pilot
	p.Start()
}

// handleCommands handle any command line arguments and return true if a command has been handled
func handleCommands(i *ctl.HostInfo) bool {
	switch len(os.Args[1:]) {
	case 0:
		// do nothing
	case 1:
		if os.Args[1] == "info" {
			// prints the host information
			fmt.Printf("%s\n", i)
		} else if os.Args[1] == "uuid" {
			// prints the host UUID
			fmt.Printf("%s\n", i.HostUUID)
		} else if os.Args[1] == "version" {
			// prints the program version
			fmt.Printf("%s\n", core.Version)
		}
	default:
		// shows usage message
		fmt.Printf("unknown argument '%s', valid arguments are 'uuid, 'version', 'info' or nothing to launch pilot\n", os.Args[1])
	}
	return len(os.Args[1:]) > 0
}
