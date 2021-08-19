package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	hostUtil "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"math"
	"net"
	"strings"
	"time"
)

const TimeLayout = "02-01-2006 03:04:05.000-0700"

// HostInfo abstracts host information
type HostInfo struct {
	MachineId       string
	HostName        string
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	Virtual         bool
	TotalMemory     float64
	CPUs            int
	HostIP          string
	BootTime        string
}

func NewHostInfo() (*HostInfo, error) {
	i, err := hostUtil.Info()
	if err != nil {
		return nil, err
	}
	// get the host IP address
	hostIp, err := externalIP()
	if err != nil {
		// if it failed to retrieve IP set to unknown
		hostIp = "unknown"
	}
	var (
		memory float64
		cpus   int
	)
	m, err := mem.VirtualMemory()
	if err == nil {
		memory = math.Round(float64(m.Total) * 9.31 * math.Pow(10, -10))
	} else {
		memory = -1
	}
	c, err := cpu.Info()
	if err == nil {
		cpus = len(c)
	} else {
		cpus = -1
	}
	info := &HostInfo{
		MachineId:   strings.ReplaceAll(i.HostID, "-", ""),
		HostIP:      hostIp,
		HostName:    i.Hostname,
		OS:          i.OS,
		Virtual:     strings.ToLower(i.VirtualizationRole) == "guest",
		BootTime:    time.Unix(int64(i.BootTime), 0).Format(TimeLayout),
		TotalMemory: memory,
		CPUs:        cpus,
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

// externalIP return host external IP
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loop back interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?\n")
}
