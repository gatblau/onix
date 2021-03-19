/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

// describes PGP keys required by functions
type File struct {
	// the unique reference for the file
	Name string `yaml:"name"`
	// a description of the intended use of this file
	Description string `yaml:"description"`
	// path to the file within the Artisan registry
	Path string
	// the file content
	Content string
}

type Files []*File
