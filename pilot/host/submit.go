/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package host

import "hash/fnv"

// SubmitJob submits
type SubmitJob struct {
}

func (c *SubmitJob) Execute() {
}

func (c *SubmitJob) Description() string {
	return "submits a payload to the remote  service"
}

func (c *SubmitJob) Key() int {
	return hashCode(c.Description())
}

func hashCode(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32())
}
