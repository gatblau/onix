//  Onix Config Manager - Dbman
//  Copyright (c) 2018-2020 by www.gatblau.org
//  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//  Contributors to this project, hereby assign copyright in this code to the project,
//  to be licensed under the same terms as the rest of the code.
package plugin

import (
	"bytes"
)

// the interface implemented by database plugins
type DatabasePlugin interface {
	// setup the plugin
	Setup(config *Conf) error

	// get db version
	GetVersion() (*Version, error)

	// execute the specified db scripts
	RunCommand(cmd *Command) (bytes.Buffer, error)

	// set the version
	SetVersion(version *Version) error

	// execute a query
	RunQuery(query *Query) (*Table, error)

	// get database server information
	GetInfo() (*DbInfo, error)
}
