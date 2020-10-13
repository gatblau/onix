/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

// file manager
// monitor configuration file changes and trigger reloads
type fileman struct {
	files []*file
}

// create a new file manager instance
func NewFileManager(r *pilot) *fileman {
	return &fileman{
		files: make([]*file, 0),
	}
}

// add a new file
func (f *fileman) add(cf *appCfg) {
	f.files = append(f.files, NewFile(cf))
}

// get the file specified by the path
func (f *fileman) get(path string) *file {
	for _, file := range f.files {
		if file.meta.Path == path {
			return file
		}
	}
	return nil
}

// stop monitoring all files
func (f *fileman) stop() {
	for _, file := range f.files {
		file.stop()
	}
}

// check if the configuration is already managed by the file manager
func (f *fileman) isManaged(config *appCfg) bool {
	for _, file := range f.files {
		if file.meta.Path == config.meta.Path {
			return true
		}
	}
	return false
}
