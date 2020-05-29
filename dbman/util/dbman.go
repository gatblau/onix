//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"github.com/gatblau/oxc"
)

var DM *DbMan

type DbMan struct {
	Cfg  *AppCfg
	info *ScriptSource
}

func NewDbMan(cfgFilePath string, cfgFileName string) (*DbMan, error) {
	cfg := NewAppCfg(cfgFilePath, cfgFileName)
	scriptClient, err := oxc.NewClient(NewClientConf(cfg))
	if err != nil {
		return nil, err
	}
	rInfo, err := NewScriptSource(cfg, scriptClient)
	if err != nil {
		return nil, err
	}
	return &DbMan{
		Cfg:  cfg,
		info: rInfo,
	}, nil
}

func (dm *DbMan) GetReleasePlan() (*Plan, error) {
	return dm.info.fetchPlan()
}

func (dm *DbMan) GetReleaseInfo(appVersion string) (*Release, error) {
	return dm.info.fetchRelease(appVersion)
}

func (dm *DbMan) SaveConfig() {
	dm.Cfg.save()
}

func (dm *DbMan) SetConfig(key string, value string) {
	dm.Cfg.set(key, value)
}

func (dm *DbMan) GetConfig(key string) {
	dm.Cfg.get(key)
}

func (dm *DbMan) PrintConfig() {
	dm.Cfg.print()
}

func (dm *DbMan) Use(filepath string, filename string) {
	dm.Cfg.load(filepath, filename)
}
