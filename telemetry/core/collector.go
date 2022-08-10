/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"context"
)

// Collector interface for host metrics
type Collector interface {
	Run(context.Context) error
	Stop()
	Restart(context.Context) error
	Status() <-chan *Status
}

// Status is the status of a collector.
type Status struct {
	Running bool
	Err     error
}
