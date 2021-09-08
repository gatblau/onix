package types

/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// Location host location
type Location struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Area host area within a Location
type Area struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Org host organisation
type Org struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
