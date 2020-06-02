//   Onix Config Db - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

// implemented by database specific implementations
type Db interface {
	// check it can connect to the database server
	CanConnect() (bool, error)
	// check the database exists
	Exists() (bool, error)
	// create the database
	Initialise() error
	// get the database version information
	GetVersion() (string, string, error)
	// deploy the schemas and functions
	Deploy() error
	// upgrade the database
	Upgrade() error
}
