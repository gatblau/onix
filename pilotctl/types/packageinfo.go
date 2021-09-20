package types

/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// PackageInfo describes a package and all its tags
type PackageInfo struct {
	Name string    `json:"name"`
	Tags []TagInfo `json:"tags,omitempty"`
}

type TagInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Ref     string `json:"ref"`
	Created string `json:"created"`
	Type    string `json:"type"`
	Size    string `json:"size"`
}
