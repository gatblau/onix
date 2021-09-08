package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import "time"

// JobBatch a representation of a batch in the database
type JobBatch struct {
	// the id of the job batch
	BatchId int64 `json:"batch_id"`
	// the name of the batch (not unique, a user-friendly name)
	Name string `json:"name"`
	// a description for the batch (not mandatory)
	Description string `json:"description,omitempty"`
	// creation time
	Created time.Time `json:"created"`
	// one or more search labels
	Label []string `json:"label,omitempty"`
	// owner
	Owner string `json:"owner"`
	// jobs
	Jobs int `json:"jobs"`
}
