//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// implemented by database specific implementations
type DatabaseProvider interface {
	// setup the provider
	Setup(config interface{})
	// execute the specified db scripts
	RunCommand(cmd interface{}) interface{}
	// execute a query
	RunQuery(query interface{}, params ...interface{}) interface{}
	// get db version
	GetVersion() interface{}
	// set the version
	SetVersion(args interface{}) error
}

// Database Provider RPC client
type DatabaseProviderRPC struct {
	client *rpc.Client
}

func (db *DatabaseProviderRPC) Setup(config interface{}) {
	err := db.client.Call("Plugin.Setup", config, new(interface{}))
	if err != nil {
		panic(err)
	}
}

func (db *DatabaseProviderRPC) RunCommand(cmd interface{}) interface{} {
	var result interface{}
	err := db.client.Call("Plugin.RunCommand", cmd, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) RunQuery(query interface{}, params ...interface{}) interface{} {
	var result interface{}
	args := make([]interface{}, 2)
	args[0] = query
	args[1] = params
	err := db.client.Call("Plugin.RunQuery", args, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) GetVersion() interface{} {
	var result interface{}
	err := db.client.Call("Plugin.GetVersion", new(interface{}), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) SetVersion(args interface{}) error {
	var result map[string]interface{}
	err := db.client.Call("Plugin.SetVersion", args, &result)
	if err != nil {
		panic(err)
	}
	return err
}

func (db *DatabaseProviderRPC) cmdResult(result interface{}) (string, error) {
	r := result.(map[string]interface{})
	strObj := r["string"]
	errObj := r["error"]
	return strObj.(string), errObj.(error)
}

func (db *DatabaseProviderRPC) queryResult(result interface{}) (Table, error) {
	r := result.(map[string]interface{})
	tableObj := r["table"]
	errObj := r["error"]
	return tableObj.(Table), errObj.(error)
}

func (db *DatabaseProviderRPC) versionResult(result interface{}) (string, string, error) {
	r := result.(map[string]interface{})
	appVerObj := r["appVersion"]
	dbVerObj := r["dbVersion"]
	errObj := r["error"]
	return appVerObj.(string), dbVerObj.(string), errObj.(error)
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type DatabaseProviderRPCServer struct {
	// This is the real implementation
	Impl DatabaseProvider
}

func (s *DatabaseProviderRPCServer) RunCommand(args interface{}, resp *interface{}) error {
	*resp = s.Impl.RunCommand(args)
	return nil
}

func (s *DatabaseProviderRPCServer) RunQuery(args interface{}, resp *interface{}) error {
	*resp = s.Impl.RunQuery(args)
	return nil
}

func (s *DatabaseProviderRPCServer) GetVersion(args interface{}, resp *interface{}) error {
	*resp = s.Impl.GetVersion()
	return nil
}

func (s *DatabaseProviderRPCServer) SetVersion(args interface{}, resp *error) error {
	*resp = s.Impl.SetVersion(args)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a DatabaseProviderRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return DatabaseProviderRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type DatabaseProviderPlugin struct {
	// Impl Injection
	Impl DatabaseProvider
}

func (p *DatabaseProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DatabaseProviderRPCServer{Impl: p.Impl}, nil
}

func (p *DatabaseProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DatabaseProviderRPC{client: c}, nil
}
