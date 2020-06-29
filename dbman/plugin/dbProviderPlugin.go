//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// the implementation of the DatabaseProvider plugin
type DatabaseProviderPlugin struct {
	// Impl Injection
	Impl DatabaseProvider
}

func (p *DatabaseProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DatabaseProviderRPCServer{Impl: p.Impl}, nil
}

func (p *DatabaseProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DatabaseProviderRPC{
		Client: c,
	}, nil
}
