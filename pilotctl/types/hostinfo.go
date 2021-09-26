package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	hostUtil "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"math"
	"net"
	"os"
	"strings"
	"time"
)

const TimeLayout = "02-01-2006 03:04:05.000-0700"

// HostInfo abstracts host information
type HostInfo struct {
	HostUUID        string
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
	MacAddress      []string
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
	macAddr, err := macAddr()
	if err != nil {
		// if it failed to retrieve media access control addresses set to unknown
		macAddr = []string{"unknown"}
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
		MacAddress:  macAddr,
	}
	// initialises host uuid
	info.InitHostUUID()
	return info, nil
}

// InitHostUUID check if there is a hostUUID file and if not generates one
// loads the UUID value in host info
// and works out a unique reference for the host that is not the machine id, as machine id is not unique
// e.g. cloning a VM will have the same machine-id, hence using a combination of hostname and machine id for uniqueness
// not using system UUID as it requires administrative privileges
func (h *HostInfo) InitHostUUID() (created bool, hostUUID string, err error) {
	// check if .hostUUID exists
	_, err = os.Stat(".hostUUID")
	// if not, creates one
	if err != nil {
		hostUUID = newUUID()
		err = os.WriteFile(".hostUUID", []byte(hostUUID), os.ModePerm)
		// if file could not be created
		if err != nil {
			// it should not continue
			return false, "", fmt.Errorf("cannot create .hostUUID file: %s", err)
		}
		// UUID was created
		created = true
	} else {
		// if the file exists
		bytes, err := ioutil.ReadFile(".hostUUID")
		if err != nil {
			// it should not continue
			return false, "", fmt.Errorf("cannot create .hostUUID file: %s", err)
		}
		hostUUID = fmt.Sprintf("%s", bytes[:])
		created = false
	}
	return created, hostUUID, nil
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

// create a Universally Unique Identifier without hyphens
func newUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// retrieve a list of mac addresses for the host
func macAddr() ([]string, error) {
	// get all network interfaces
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// fetch mac addresses from all interfaces
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
