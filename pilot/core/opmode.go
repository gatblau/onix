/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

// the pilot operating mode
type opMode int

const (
	// sits side by side with the application
	Sidecar opMode = iota
	// launch the application as a subprocess
	Controller
	// launch applications in containers
	Host
)

// string representation of the operating mode
func (m opMode) String() string {
	switch m {
	case Sidecar:
		return "Sidecar"
	case Controller:
		return "Controller"
	case Host:
		return "Host"
	}
	return "?unknown mode?"
}
