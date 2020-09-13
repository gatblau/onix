/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	hostUtil "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"math"
	"strings"
)

type HostInfo struct {
	// unique identifier for the host
	HostID string
	// a host label (domain name) to uniquely identify it in various forms of electronic communication
	HostName string
	// the host Operating System
	OS string
	// OS parameters
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	// is the host a virtual or physical machine?
	Virtual bool
	// Memory
	TotalMemory float64
	// CPU
	CPUs int
}

func NewHostInfo() (*HostInfo, error) {
	info := new(HostInfo)
	i, err := hostUtil.Info()
	if err != nil {
		return nil, err
	}
	info.HostName = i.Hostname
	info.OS = i.OS
	info.HostID = i.HostID
	info.Virtual = strings.ToLower(i.VirtualizationRole) == "guest"
	m, err := mem.VirtualMemory()
	if err == nil {
		info.TotalMemory = math.Round(float64(m.Total) * 9.31 * math.Pow(10, -10))
	}
	c, err := cpu.Info()
	if err == nil {
		info.CPUs = len(c)
	}
	return info, nil
}

func (h *HostInfo) String() string {
	bytes, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
