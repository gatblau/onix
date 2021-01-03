/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

type Manifest struct {
	// the package type
	Type string `json:"type,omitempty"`
	// the license associated to the package
	License string `json:"license"`
	// the name of the package file
	Ref string `json:"ref"`
	// the build profile used
	Profile string `json:"profile"`
	// the labels assigned to the package
	Labels map[string]string `json:"labels,omitempty"`
	// the URI of the package source
	Source string `json:"source,omitempty"`
	// the path within the source where the project is (for uber repos)
	SourcePath string `json:"source_path,omitempty"`
	// the commit hash
	Commit string `json:"commit,omitempty"`
	// repo branch
	Branch string `json:"branch,omitempty"`
	// repo tag
	Tag string `json:"tag,omitempty"`
	// the name of the file or folder that has been packaged
	Target string `json:"target,omitempty"`
	// the timestamp
	Time string `json:"time"`
	// the size of the artefact
	Size string `json:"size"`
	// true if the target was zipped previous to packaging (e.g. jar files)
	Zip bool `json:"zip"`
	// what functions are available to call?
	Functions []*FxInfo `json:"functions,omitempty"`
}

// exported function list
type FxInfo struct {
	Name        string
	Description string
	Input       []*Input
}
