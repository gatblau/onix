/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

type procStatus int

const (
	started procStatus = iota
	stopped
	stopRequested
)

func (s procStatus) String() string {
	switch s {
	case started:
		return "started"
	case stopped:
		return "stopped"
	case stopRequested:
		return "stop requested"
	}
	return "unknown"
}
