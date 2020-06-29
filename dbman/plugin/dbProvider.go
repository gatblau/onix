//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

// the interface implemented by database specific implementations of a database provider
type DatabaseProvider interface {
	// setup the provider
	// config: a map[string]interface{} serialised as a JSON string, containing DbMan's current config set
	// result: a map[string]interface{} serialised as a JSON string, containing log and error items
	Setup(config string) string

	// get db version
	GetVersion() string

	// execute the specified db scripts
	RunCommand(cmd string) string

	// set the version
	SetVersion(versionInfo string) string

	// execute a query
	RunQuery(query string) string
}
