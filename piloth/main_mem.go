// +build mem

/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import "github.com/pkg/profile"

func main() {
	// memory profiling
	defer profile.Start(profile.MemProfile).Stop()
	run()
}
