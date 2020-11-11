package registry

import (
	"github.com/gatblau/onix/artie/core"
	"log"
)

type Factory struct {
	conf *core.ServerConfig
}

func NewBackendFactory() *Factory {
	return &Factory{
		conf: &core.ServerConfig{},
	}
}

func (f *Factory) Get() Remote {
	// get the configured factory
	switch f.conf.Backend() {
	case core.FileSystem:
		return new(RemoteFs)
	case core.Nexus3:
		return NewNexus3Registry(
			f.conf.BackendDomain(), // the nexus scheme://domain:port
		)
	}
	log.Fatal("backend not recognised")
	return nil
}
