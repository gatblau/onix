//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func NewDatabase(cfg *Config) (*Database, error) {
	provider, client, err := getDbProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Database{
		provider: provider,
		client:   client,
	}, nil
}

type Database struct {
	provider DatabaseProvider
	client   *plugin.Client
}

// safely terminate the rpc client
func (db *Database) Close() {
	db.client.Kill()
}

func (db *Database) Provider() DatabaseProvider {
	return db.provider
}

// load a database provider plugin
func getDbProvider(cfg *Config) (DatabaseProvider, *plugin.Client, error) {
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

// return the executable path
func execPath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
