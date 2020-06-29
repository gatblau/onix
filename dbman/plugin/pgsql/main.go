//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	. "github.com/gatblau/onix/dbman/plugin"
	"github.com/hashicorp/go-plugin"
)

// entry point for the PGSQL plugin
func main() {
	// launch the plugin as an rpc server
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "dbman-db-provider",
			MagicCookieValue: "dbman-db-pgsql",
		},
		Plugins: map[string]plugin.Plugin{
			"pgsql": &DatabaseProviderPlugin{
				Impl: new(PgSQLProvider),
			},
		},
	})
}
