package peerdisco

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/pilotctl/types"
	"time"
)

type Links []Link

func (l Links) Equals(links Links) bool {
	panic("not implemented")
}

// Link information about the discovered pilot peer
type Link struct {
	UUID     string    `json:"uuid"`
	Address  string    `json:"address,omitempty"`
	Name     string    `json:"name"`
	BootTime time.Time `json:"boot_time"`
}

func NewLink(hostInfo types.HostInfo) Link {
	return Link{
		UUID:     hostInfo.HostUUID,
		Name:     hostInfo.HostName,
		BootTime: hostInfo.BootTime,
	}
}
