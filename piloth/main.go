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
	"strings"
)

func main() {
	printMachineId()
	p, err := core.NewPilot()
	if err != nil {
		fmt.Printf("cannot start pilot: %s\n", err)
		os.Exit(1)
	}
	p.Start()
}

func printMachineId() {
	i, err := core.NewHostInfo()
	if err != nil {
		panic(err)
	}
	// if machineid is passed as a command line argument
	switch len(os.Args[1:]) {
	case 0:
		// do nothing
	case 1:
		if os.Args[1] == "info" {
			// prints the machine id
			fmt.Printf("%s\n", i)
			// terminates programme
			os.Exit(0)
		} else if os.Args[1] == "machineid" {
			// prints the machine id
			fmt.Printf("%s\n", strings.Replace(i.MachineId, "-", "", -1))
			// terminates programme
			os.Exit(0)
		} else if os.Args[1] == "uuid" {
			// prints the machine id
			fmt.Printf("%s\n", strings.Replace(i.HostUUID, "-", "", -1))
			// terminates programme
			os.Exit(0)
		}
	default:
		// prints the machine id
		fmt.Printf("unknown argument '%s', valid argument is 'machineid' or 'info' or nothing\n", os.Args[1])
		// terminates programme
		os.Exit(0)
	}
}
