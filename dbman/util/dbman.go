//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import "github.com/gatblau/oxc"

var DM *DbMan

type DbMan struct {
	cfg  *Config
	info *ScriptSource
}

func NewDbMan(cfgFilePath string) (*DbMan, error) {
	cfg, err := NewConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	scriptClient, err := oxc.NewClient(NewClientConf(cfg))
	if err != nil {
		return nil, err
	}
	rInfo, err := NewScriptSource(cfg, scriptClient)
	if err != nil {
		return nil, err
	}
	return &DbMan{
		cfg:  cfg,
		info: rInfo,
	}, nil
}

func (dm *DbMan) GetReleasePlan() (*ReleasePlan, error) {
	return dm.info.fetchPlan()
}
