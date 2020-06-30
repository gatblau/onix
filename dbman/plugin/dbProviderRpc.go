//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

import "net/rpc"

// Database Provider RPC client
type DatabaseProviderRPC struct {
	Client *rpc.Client
}

func (db *DatabaseProviderRPC) Setup(config string) string {
	var result string
	err := db.Client.Call("Plugin.Setup", config, &result)
	if err != nil {
		output := NewParameter()
		output.SetError(err)
		return output.ToString()
	}
	return result
}

func (db *DatabaseProviderRPC) GetVersion() string {
	var result string
	err := db.Client.Call("Plugin.GetVersion", "", &result)
	if err != nil {
		return db.errorToString(err)
	}
	return result
}

func (db *DatabaseProviderRPC) RunCommand(cmd string) string {
	var result string
	err := db.Client.Call("Plugin.RunCommand", cmd, &result)
	if err != nil {
		return db.errorToString(err)
	}
	return result
}

func (db *DatabaseProviderRPC) SetVersion(args string) string {
	var result string
	err := db.Client.Call("Plugin.SetVersion", args, &result)
	if err != nil {
		return db.errorToString(err)
	}
	return result
}

func (db *DatabaseProviderRPC) RunQuery(query string) string {
	var result string
	err := db.Client.Call("Plugin.RunQuery", query, &result)
	if err != nil {
		return db.errorToString(err)
	}
	return result
}

func (db *DatabaseProviderRPC) errorToString(err error) string {
	output := NewParameter()
	output.SetError(err)
	return output.ToString()
}
