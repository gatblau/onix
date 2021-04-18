/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package host

// PingJob pings the remote service periodically
type PingJob struct {
}

func (c *PingJob) Execute() {
}

func (c *PingJob) Description() string {
	return "pings the remote service"
}

func (c *PingJob) Key() int {
	return hashCode(c.Description())
}
