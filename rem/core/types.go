package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

type Command struct {
	Package  string            `json:"package"`
	Function string            `json:"function"`
	Input    map[string]string `json:"input"`
}

type Host struct {
	Name      string `json:"name"`
	Customer  string `json:"customer"`
	Region    string `json:"region"`
	Location  string `json:"location"`
	Connected bool   `json:"connected"`
	Up        bool   `json:"up"`
}
