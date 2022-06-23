/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import (
    "encoding/json"
    "fmt"
    "github.com/gatblau/onix/artisan/core"
    "github.com/shirou/gopsutil/cpu"
    hostUtil "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "math"
    "net"
    "os/exec"
    "regexp"
    "runtime"
    "strings"
    "time"
)

// const TimeLayout = "02-01-2006 03:04:05.000-0700"

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
    BootTime        time.Time
    MacAddress      []string
    PrimaryMAC      string
    HardwareId      string
}

func NewHostInfo() (*HostInfo, error) {
    i, err := hostUtil.Info()
    if err != nil {
        return nil, err
    }
    // get the host IP address
    hostIp, err := getPrimaryIP()
    if err != nil {
        // if it failed to retrieve IP set to unknown
        hostIp = "unknown"
    }
    primaryMAC, macList, err := macAddr(hostIp)
    if err != nil {
        // if it failed to retrieve media access control addresses set to unknown
        macList = []string{"unknown"}
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
        BootTime:    time.Unix(int64(i.BootTime), 0),
        TotalMemory: memory,
        CPUs:        cpus,
        MacAddress:  macList,
        PrimaryMAC:  primaryMAC,
        HardwareId:  getHwId(),
    }
    // return
    return info, nil
}

func (h *HostInfo) String() string {
    bytes, err := json.Marshal(h)
    if err != nil {
        return err.Error()
    }
    return string(bytes)
}

// retrieve a list of mac addresses for the host
func macAddr(primaryIp string) (string, []string, error) {
    var primaryMAC string
    // get all network interfaces
    ifas, err := net.Interfaces()
    if err != nil {
        return "", nil, err
    }
    // fetch mac addresses from all interfaces
    var as []string
    for _, ifa := range ifas {
        addrs, err2 := ifa.Addrs()
        if err2 != nil {
            return "", nil, fmt.Errorf("cannot retrieve list of unicast interface addresses for network interface: %s", err2)
        }
        // var ip net.IP
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            // process IP address
            if ip.To4() != nil && ip.To4().String() == primaryIp {
                primaryMAC = ifa.HardwareAddr.String()
                break
            }
        }
        a := ifa.HardwareAddr.String()
        if a != "" {
            as = append(as, a)
        }
    }
    return primaryMAC, as, nil
}

// getPrimaryIP gets the IP address of the primary network interface
// In oder to find the primary interface, it triggers a fake udp connection so that the underlying host
// can resolve the route to use
func getPrimaryIP() (string, error) {
    conn, err := net.Dial("udp", "1.2.3.4:80")
    if err != nil {
        return "", fmt.Errorf("cannot dial connection to retrieve primary IP address: %s", err)
    }
    defer conn.Close()
    ipAddress := conn.LocalAddr().(*net.UDPAddr)
    return ipAddress.IP.String(), nil
}

func getHwId() string {
    switch runtime.GOOS {
    case "linux":
        c := exec.Command("sudo", "dmidecode", "-s", "system-uuid")
        out, err := c.Output()
        if err != nil {
            core.ErrorLogger.Printf("cannot retrieve system uuid: %s\n", err)
            return ""
        }
        return string(out[:])
    case "darwin":
        // system_profiler SPHardwareDataType | grep "Hardware UUID"
        c := exec.Command("system_profiler", "SPHardwareDataType")
        out, err := c.Output()
        if err != nil {
            core.ErrorLogger.Printf("cannot find system uuid: %s\n", err)
            return ""
        }
        r, _ := regexp.Compile(".*Hardware UUID: (?P<HW_ID>.*)\\b")
        match := r.FindStringSubmatch(string(out[:]))
        result := make(map[string]string)
        for i, name := range r.SubexpNames() {
            if i != 0 && name != "" {
                result[name] = match[i]
            }
        }
        if _, exists := result["HW_ID"]; !exists {
            core.ErrorLogger.Printf("cannot find system uuid for OS '%s': %s\n", runtime.GOOS, out)
            return ""
        }
        return result["HW_ID"]
    default:
        core.ErrorLogger.Printf("does not support finding system uuid for '%s'\n", runtime.GOOS)
        return ""
    }
    return ""
}
