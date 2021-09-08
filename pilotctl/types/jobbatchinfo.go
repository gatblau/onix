package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// JobBatchInfo information required to create a new job batch
type JobBatchInfo struct {
	// the name of the batch (not unique, a user-friendly name)
	Name string `json:"name"`
	// a description for the batch (not mandatory)
	Description string `json:"description,omitempty"`
	// one or more search labels
	Label []string `json:"label,omitempty"`
	// the universally unique host identifier created by pilot
	HostUUID []string `json:"host_uuid"`
	// the unique key of the function to run
	FxKey string `json:"fx_key"`
	// the version of the function to run
	FxVersion int64 `json:"fx_version"`
}
