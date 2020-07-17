//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

import (
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
	"strings"
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
	// declare the database provider instance
	var provider DatabasePlugin

	// what database provider to use?
	dbProvider := cfg.GetString(DbProvider)
	// if the provider starts with _ then it is a native provider
	if strings.HasPrefix(dbProvider, "_") {
		switch strings.ToLower(dbProvider) {
		// the PostgreSQL native db provider
		case "_pgsql":
			// create a plugin instance
			provider = &PgSQLProvider{}
		default:
			// there is not any native provider implemented for the required name
			return nil, nil, errors.New(fmt.Sprintf("!!! I do not support a native database provider called '%s'", dbProvider))
		}
		// retrieves the configuration
		conf, _ := NewConf(cfg.ToString())
		// passes the configuration to the provider
		provider.Setup(conf)
		// creates an instance of the plugin wrapper
		d := &DatabasePluginDecorator{Plugin: provider}
		// returns the wrapper
		return d, nil, nil
	}

	// if the provider name does not start with _ then it is a plugin
	// create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Error,
	})

	// start by launching the plugin process
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

	// connect to the db plugin via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("I cannot load db provider %s (%s)\n", dbProvider, err))
	}

	// request the plugin
	raw, err := rpcClient.Dispense(dbProvider)
	if err != nil {
		return nil, nil, err
	}

	// We should have a DatabaseProvider now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	db := raw.(DatabaseProvider)

	// return the provider
	return db, client, nil
}
