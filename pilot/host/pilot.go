/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package host

import (
	"github.com/gatblau/onix/pilot/core"
)

// pilot host mode
type Pilot struct {
	cfg  *core.Config
	info *HostInfo
}

func NewPilot() (*Pilot, error) {
	// read configuration
	cfg := &core.Config{}
	err := cfg.Load()
	if err != nil {
		return nil, err
	}
	info, err := NewHostInfo()
	if err != nil {
		return nil, err
	}
	p := &Pilot{
		cfg:  cfg,
		info: info,
	}
	// return a new pilot
	return p, nil
}

func (p *Pilot) Start() {
	// check registration
	// start ping loop
	// start submit loop
}
