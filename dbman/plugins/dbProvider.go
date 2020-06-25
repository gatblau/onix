//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

import (
	// "fmt"
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// implemented by database specific implementations
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

// Database Provider RPC client
type DatabaseProviderRPC struct {
	client *rpc.Client
}

func (db *DatabaseProviderRPC) Setup(config string) string {
	var result string
	err := db.client.Call("Plugin.Setup", config, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) GetVersion() string {
	var result string
	err := db.client.Call("Plugin.GetVersion", "", &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) RunCommand(cmd string) string {
	var result string
	err := db.client.Call("Plugin.RunCommand", cmd, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) SetVersion(args string) string {
	var result string
	err := db.client.Call("Plugin.SetVersion", args, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DatabaseProviderRPC) RunQuery(query string) string {
	var result string
	err := db.client.Call("Plugin.RunQuery", query, &result)
	if err != nil {
		panic(err)
	}
	return result
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type DatabaseProviderRPCServer struct {
	// This is the real implementation
	Impl DatabaseProvider
}

func (s *DatabaseProviderRPCServer) Setup(args string, resp *string) error {
	*resp = s.Impl.Setup(args)
	return nil
}

func (s *DatabaseProviderRPCServer) GetVersion(args string, resp *string) error {
	*resp = s.Impl.GetVersion()
	return nil
}

func (s *DatabaseProviderRPCServer) RunCommand(args string, resp *string) error {
	*resp = s.Impl.RunCommand(args)
	return nil
}

func (s *DatabaseProviderRPCServer) SetVersion(args string, resp *string) error {
	*resp = s.Impl.SetVersion(args)
	return nil
}

func (s *DatabaseProviderRPCServer) RunQuery(args string, resp *string) error {
	*resp = s.Impl.RunQuery(args)
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
