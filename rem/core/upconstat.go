/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

// UpdateConnStatusJob updates the connection status based on ping age
type UpdateConnStatusJob struct {
}

func (c *UpdateConnStatusJob) Execute() {
}

func (c *UpdateConnStatusJob) Description() string {
	return "updates the connection status based on ping age"
}

func (c *UpdateConnStatusJob) Key() int {
	return hashCode(c.Description())
}
