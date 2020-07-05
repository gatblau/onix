//  Onix Config Manager - Dbman
//  Copyright (c) 2018-2020 by www.gatblau.org
//  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//  Contributors to this project, hereby assign copyright in this code to the project,
//  to be licensed under the same terms as the rest of the code.
package plugin

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-plugin"
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

// launch the database plugin
func ServeDbPlugin(pluginName string, impl DatabasePlugin) {
	// launch the plugin as an rpc server
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "dbman-db-provider",
			MagicCookieValue: fmt.Sprintf("dbman-db-%s", pluginName),
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &DatabaseProviderPlugin{
				Impl: &DatabasePluginDecorator{
					Plugin: impl,
				},
			},
		},
	})
}

// the decorator wraps the DatabasePlugin interface and exposes it as a DatabaseProvider interface
// the DatabaseProvider is the underlying interface used by the net/rpc protocol to communicate with the plugin
// whereas the DatabasePlugin interface is a friendlier version used by plugin writers
type DatabasePluginDecorator struct {
	Plugin DatabasePlugin
}

func (db *DatabasePluginDecorator) Setup(config string) string {
	output := NewParameter()
	// parse the configuration
	c, err := NewConf(config)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	// allocate the parsed object to cfg
	db.Plugin.Setup(c)
	// return the output
	return output.ToString()
}

// RPC serialisation wrapper for getting database version information
func (db *DatabasePluginDecorator) GetVersion() string {
	// create the output struct
	output := NewParameter()
	// call the plugin operation
	version, err := db.Plugin.GetVersion()
	// if an error is found
	if err != nil {
		// return the error
		return output.ToError(err)
	}
	// set the result value
	output.Set("result", version)
	// return the serialised output back to the RPC client
	return output.ToString()
}

func (db *DatabasePluginDecorator) RunCommand(command string) string {
	output := NewParameter()
	cmd, err := NewCommand(command)
	if err != nil {
		return output.ToError(err)
	}
	log, err := db.Plugin.RunCommand(cmd)
	if log.Len() > 0 {
		output.Log(log.String())
	}
	if err != nil {
		return output.ToError(err)
	}
	return output.ToString()
}

func (db *DatabasePluginDecorator) RunQuery(queryInfo string) string {
	output := NewParameter()
	query, err := NewQuery(queryInfo)
	if err != nil {
		return output.ToError(err)
	}
	result, err := db.Plugin.RunQuery(query)
	if err != nil {
		return output.ToError(err)
	}
	output.Set("result", result)
	return output.ToString()
}

func (db *DatabasePluginDecorator) SetVersion(versionInfo string) string {
	output := NewParameter()
	v, err := NewVersion(versionInfo)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	err = db.Plugin.SetVersion(v)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	return v.ToString()
}

func (db *DatabasePluginDecorator) GetInfo() string {
	// create the output struct
	output := NewParameter()
	// call the plugin operation
	info, err := db.Plugin.GetInfo()
	// if an error is found
	if err != nil {
		// return the error
		return output.ToError(err)
	}
	// set the result value
	output.Set("result", info)
	// return the serialised output back to the RPC client
	return output.ToString()
}
