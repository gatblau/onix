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
	Setup(config map[string]interface{})
	// execute the specified db scripts
	RunCommand(cmd map[string]interface{}) map[string]interface{}
	// execute a query
	RunQuery(query map[string]interface{}, params ...interface{}) map[string]interface{}
	// get db version
	GetVersion() map[string]interface{}
	// set the version
	SetVersion(args map[string]interface{}) error
}

// Database Provider RPC client
type DatabaseProviderRPC struct {
	client *rpc.Client
}

func (db *DatabaseProviderRPC) Setup(config map[string]interface{}) {
	err := db.client.Call("Plugin.Setup", config, nil)
	if err != nil {
		panic(err)
	}
}

func (db *DatabaseProviderRPC) RunCommand(cmd map[string]interface{}) (string, error) {
	var result map[string]interface{}
	err := db.client.Call("Plugin.RunCommand", cmd, &result)
	if err != nil {
		panic(err)
	}
	return db.cmdResult(result)
}

func (db *DatabaseProviderRPC) RunQuery(query map[string]interface{}) (Table, error) {
	var result map[string]interface{}
	err := db.client.Call("Plugin.RunQuery", query, &result)
	if err != nil {
		panic(err)
	}
	return db.queryResult(result)
}

func (db *DatabaseProviderRPC) GetVersion() (string, string, error) {
	var result map[string]interface{}
	err := db.client.Call("Plugin.GetVersion", nil, &result)
	if err != nil {
		panic(err)
	}
	return db.versionResult(result)
}

func (db *DatabaseProviderRPC) SetVersion(args map[string]interface{}) error {
	var result map[string]interface{}
	err := db.client.Call("Plugin.SetVersion", args, &result)
	if err != nil {
		panic(err)
	}
	return err
}

func (db *DatabaseProviderRPC) cmdResult(result map[string]interface{}) (string, error) {
	strObj := result["string"]
	errObj := result["error"]
	return strObj.(string), errObj.(error)
}

func (db *DatabaseProviderRPC) queryResult(result map[string]interface{}) (Table, error) {
	tableObj := result["table"]
	errObj := result["error"]
	return tableObj.(Table), errObj.(error)
}

func (db *DatabaseProviderRPC) versionResult(result map[string]interface{}) (string, string, error) {
	appVerObj := result["appVersion"]
	dbVerObj := result["dbVersion"]
	errObj := result["error"]
	return appVerObj.(string), dbVerObj.(string), errObj.(error)
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type DatabaseProviderRPCServer struct {
	// This is the real implementation
	Impl DatabaseProvider
}

func (s *DatabaseProviderRPCServer) RunCommand(args map[string]interface{}, resp *map[string]interface{}) error {
	*resp = s.Impl.RunCommand(args)
	return nil
}

func (s *DatabaseProviderRPCServer) RunQuery(args map[string]interface{}, resp *map[string]interface{}) error {
	*resp = s.Impl.RunQuery(args)
	return nil
}

func (s *DatabaseProviderRPCServer) GetVersion(args interface{}, resp *map[string]interface{}) error {
	*resp = s.Impl.GetVersion()
	return nil
}

func (s *DatabaseProviderRPCServer) SetVersion(args map[string]interface{}, resp *error) error {
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
