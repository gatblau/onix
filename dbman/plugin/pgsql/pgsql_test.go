package main

import (
	"github.com/gatblau/onix/dbman/core"
	. "github.com/gatblau/onix/dbman/plugin"
	"testing"
)

func TestPgSQLProvider_GetVersion(t *testing.T) {
	// read the config
	cfg := core.NewConfig("", "")
	// transform cfg into map
	conf, _ := NewConf(cfg.All())
	// creates the provider
	dbProvider := &PgSQLProvider{
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
