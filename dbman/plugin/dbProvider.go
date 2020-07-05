//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

// the interface implemented by database specific implementations of a database provider
type DatabaseProvider interface {
	// setup the provider with the specified configuration information
	Setup(config string) string

	// get database server general information
	GetInfo() string

	// get database release version information
	GetVersion() string

	// set database release version information
	SetVersion(versionInfo string) string

	// execute the specified command
	RunCommand(cmd string) string

	// execute the specified query
	RunQuery(query string) string
}
