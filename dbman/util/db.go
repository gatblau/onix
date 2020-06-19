//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"errors"
	"fmt"
	"strings"
)

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

// creates a new db instance
func NewDbProvider(appCfg *Config) DatabaseProvider {
	switch strings.ToLower(appCfg.Get(DbProvider)) {
	case "pgsql":
		// load the default native postgres provider
		return &PgSQLProvider{
			cfg: appCfg,
		}
	default:
		// only supports connections to postgres at the moment
		// in time, a plugin approach for database providers could be implemented
		panic(errors.New(fmt.Sprintf("!!! the database provider '%v' is not supported.", appCfg.Get(DbProvider))))
	}
}
