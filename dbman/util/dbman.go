//   Onix Config Db - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"fmt"
	"github.com/gatblau/oxc"
)

var DM *DbMan

type DbMan struct {
	// configuration
	Cfg *AppCfg
	// scrips manager
	script *ScriptManager
	// database manager
	db Db
}

func NewDbMan() (*DbMan, error) {
	// create an instance of the current configuration set
	cfg := NewAppCfg("", "")
	// create an instance of the script http client
	scriptClient, err := oxc.NewClient(NewOxClientConf(cfg))
	if err != nil {
		return nil, err
	}
	// create an instance of the script manager
	rInfo, err := NewScriptManager(cfg, scriptClient)
	if err != nil {
		return nil, err
	}
	pgdb := NewPgDb(cfg)
	return &DbMan{
		Cfg:    cfg,
		script: rInfo,
		db:     pgdb,
	}, nil
}

func (dm *DbMan) GetReleasePlan() (*Plan, error) {
	return dm.script.fetchPlan()
}

func (dm *DbMan) GetReleaseInfo(appVersion string) (*Release, error) {
	return dm.script.fetchRelease(appVersion)
}

func (dm *DbMan) SaveConfig() {
	dm.Cfg.save()
}

func (dm *DbMan) SetConfig(key string, value string) {
	dm.Cfg.set(key, value)
}

func (dm *DbMan) GetConfig(key string) {
	dm.Cfg.Get(key)
}

func (dm *DbMan) PrintConfig() {
	dm.Cfg.print()
}

func (dm *DbMan) Use(filepath string, name string) {
	dm.Cfg.load(filepath, name)
}

func (dm *DbMan) GetCurrentConfigFile() string {
	return dm.Cfg.ConfigFileUsed()
}

func (dm *DbMan) GetCurrentDir() string {
	return dm.Cfg.root.path()
}

func (dm *DbMan) Check() {
	fmt.Printf("checking: can I connect to schema.uri? : '%v'\n", dm.Cfg.Get(SchemaURI))
	_, err := dm.script.fetchPlan()
	if err != nil {
		fmt.Printf("oops! check failed: %v\n", err)
	} else {
		fmt.Printf("yeah! check suceeded!\n")
	}
	fmt.Printf("checking: can I connect to db.provider? : '%v'\n", dm.Cfg.Get(DbProvider))
	_, err = dm.db.CanConnect()
	if err != nil {
		fmt.Printf("oops! check failed: %v\n", err)
	} else {
		fmt.Printf("yeah! check suceeded!\n")
	}
}
