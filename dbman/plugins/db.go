//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

// implemented by database specific implementations
type DatabaseProvider interface {
	// setup the provider
	Setup(config *Config)
	// execute the specified db scripts
	RunCommand(cmd *Command) (string, error)
	// execute a query
	RunQuery(query *Query, params ...interface{}) (Table, error)
	// get db version
	GetVersion() (appVersion string, dbVersion string, err error)
	// set the version
	SetVersion(appVersion string, dbVersion string, description string, source string) error
}
