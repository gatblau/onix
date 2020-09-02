/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

// https://billglover.me/2020/01/12/the-sidecar-pattern/

const (
	// pilot runs in a OS host launching applications in containers
	HostCtl PilotMode = iota
	// pilot controls the application process fetching, refreshing and reloading the configuration
	AppCtl
	// pilot runs in a sidecar container refreshing app configuration (using POSIX shared memory)
	SideCar
	// pilot runs as an init container fetching the initial application configuration (using POSIX shared memory)
	InitC
)

type PilotMode int

func (s PilotMode) String() string {
	return [...]string{"HostCtl", "AppCtl", "SideCar", "InitC"}[s]
}
