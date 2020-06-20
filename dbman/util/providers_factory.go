//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugins"
	"plugin"
)

// load a database provider plugin
func LoadDbProviderPlugin(cfg *Config) (DatabaseProvider, error) {
	// the database provider to load
	var provider DatabaseProvider
	// get the plugin name to use
	pluginName := fmt.Sprintf("%.so", cfg.Get(DbProvider))
	// load the plugin
	plug, err := plugin.Open(pluginName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("!!! I cannot load the db provider file '%s': %s", pluginName, err))
	}
	instance, err := plug.Lookup("DbProvider")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("!!! I cannot lookup the db provider instance'%s': %s", pluginName, err))
	}
	// check the provider instance complies with the provider interface
	provider, ok := instance.(DatabaseProvider)
	if !ok {
		return nil, errors.New(fmt.Sprintf("!!! I cannot use the db provider in '%s' as it does not comply with the DatabaseProvider interface", pluginName))
	}
	// pass the DbMan's configuration to the provider
	provider.Setup(cfg)
	// return the database provider
	return provider, nil
}
