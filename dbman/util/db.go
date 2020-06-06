//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

// implemented by database specific implementations
type DatabaseProvider interface {
	// check it can connect to the database server
	CanConnectToServer() (bool, error)
	// check the database exists
	DbExists() (bool, error)
	// create the database
	InitialiseDb(init *DbInit) error
	// get the database version information
	GetVersion() (string, string, error)
	// deploy the schemas and functions
	DeployDb(release *Release) error
	// upgrade the database
	UpgradeDb() error
	// create database version tracking table
	CreateVersionTable() error
	// insert version information in the tracking table
	InsertVersion(appVersion string, dbVersion string, description string, origin string) error
}
