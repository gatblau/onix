package types

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"time"
)

// JobResult the result of the execution of a job
// note: ensure it is aligned with the same struct in piloth
type JobResult struct {
	// the unique job id
	JobId int64
	// indicates of the job was successful
	Success bool
	// the execution log for the job
	Log string
	// the error if any
	Err string
	// the completion time
	Time time.Time
}
