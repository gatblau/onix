package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import "fmt"

func NewUpdateConnStatusJob(rem *ReMan) (*UpdateConnStatusJob, error) {
	conf := NewConf()
	return &UpdateConnStatusJob{
		rem:          rem,
		pingInterval: conf.GetPingInterval(),
	}, nil
}

// UpdateConnStatusJob updates the connection status based on ping age
type UpdateConnStatusJob struct {
	rem          *ReMan
	pingInterval int
}

func (c *UpdateConnStatusJob) Execute() {
	err := c.rem.RecordConnStatus(c.pingInterval)
	if err != nil {
		fmt.Printf("ERROR: cannot check for disconnected events, %s\n", err)
	}
}

func (c *UpdateConnStatusJob) Description() string {
	return "updates the disconnected status of hosts based on last seen time"
}

func (c *UpdateConnStatusJob) Key() int {
	return hashCode(c.Description())
}
