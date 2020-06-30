//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

import (
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"log"
	"os"
	"os/exec"
)

func NewDatabase(cfg *Config) (*DatabaseProviderManager, error) {
	provider, client, err := getDbProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &DatabaseProviderManager{
		provider: provider,
		client:   client,
	}, nil
}

// manages the lifecycle of a DatabaseProvider
type DatabaseProviderManager struct {
	provider DatabaseProvider
	client   *plugin.Client
}

// safely terminate the rpc client
func (db *DatabaseProviderManager) Close() {
	db.client.Kill()
}

func (db *DatabaseProviderManager) Provider() DatabaseProvider {
	return db.provider
}

// load a database provider plugin
func getDbProvider(cfg *Config) (DatabaseProvider, *plugin.Client, error) {
	// what database provider to use?
	dbProvider := cfg.GetString(DbProvider)

	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Error,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "dbman-db-provider",
			MagicCookieValue: fmt.Sprintf("dbman-db-%s", dbProvider),
		},
		Plugins: map[string]plugin.Plugin{
			dbProvider: &DatabaseProviderPlugin{},
		},
		Cmd:    exec.Command(fmt.Sprintf("./dbman-db-%s", dbProvider)),
		Logger: logger,
	})
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(dbProvider)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a DatabaseProvider now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	db := raw.(DatabaseProvider)

	// return
	return db, client, nil
}
