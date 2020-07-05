package main

import (
	"fmt"
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
	v, err := dbProvider.GetVersion()
	// check for error
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(v)
}

func TestPgSQLProvider_GetDbInfo(t *testing.T) {
	// read the config
	cfg := core.NewConfig("", "")
	// transform cfg into map
	conf, _ := NewConf(cfg.All())
	// creates the provider
	dbProvider := &PgSQLProvider{
		cfg: conf,
	}
	// test get version
	i, err := dbProvider.GetInfo()
	// check for error
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(i)
}
