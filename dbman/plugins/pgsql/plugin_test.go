package main

import (
	. "github.com/gatblau/onix/dbman/plugins"
	"github.com/hashicorp/go-hclog"
	"os"
	"testing"
)

func TestPgSQLProvider_GetVersion(t *testing.T) {
	// creates a logger
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	// read the config
	cfg := NewAppCfg("", "")
	// transform cfg into map
	conf, _ := NewConf(cfg.All())
	// creates the provider
	dbProvider := &PgSQLProvider{
		log: logger,
		cfg: conf,
	}
	// test get version
	jsonResult := dbProvider.GetVersion()
	// reads the result
	r := NewParameterFromJSON(jsonResult)
	// check for error
	if r.HasError() {
		t.Error(r.Error())
		t.Fail()
	}
}
