/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
    "encoding/json"
    "fmt"
    "github.com/gatblau/onix/artisan/core"
    ctl "github.com/gatblau/onix/pilotctl/types"
    core2 "github.com/gatblau/onix/piloth/core"
    "github.com/spf13/cobra"
    "strings"
)

// ConfigCmd retrieves host/device configuration
type ConfigCmd struct {
    cmd *cobra.Command
}

func NewConfigCmd() *ConfigCmd {
    c := &ConfigCmd{
        cmd: &cobra.Command{
            Use:   "config [host-uuid|mac-addr|hw-id|machine-id|ALL]",
            Short: "retrieves host/device configuration",
            Long:  `retrieves host/device configuration`,
        },
    }
    c.cmd.Run = c.Run
    return c
}

func (c *ConfigCmd) Run(cmd *cobra.Command, args []string) {
    // collects device/host information
    hostInfo, err := ctl.NewHostInfo()
    if err != nil {
        core.RaiseErr("cannot collect host information")
    }
    // load host uuid from activation key if present
    hostInfo.HostUUID = hostUUID()
    switch len(args) {
    case 0:
        // prints the all host information
        i, _ := json.MarshalIndent(hostInfo, "", "  ")
        // prints the host information
        fmt.Printf("%s\n", i)
    case 1:
        if strings.ToUpper(args[0]) == "ALL" {
            i, _ := json.MarshalIndent(hostInfo, "", "  ")
            // prints the host information
            fmt.Printf("%s\n", i)
        } else if args[0] == "host-uuid" {
            fmt.Printf("%s\n", hostInfo.HostUUID)
        } else if args[0] == "mac-addr" {
            // prints the host UUID
            fmt.Printf("%s\n", hostInfo.PrimaryMAC)
        } else if args[0] == "hw-id" {
            // prints the host hardware system uuid
            fmt.Printf("%s\n", hostInfo.HardwareId)
        } else if args[0] == "machine-id" {
            // prints the host hardware system uuid
            fmt.Printf("%s\n", hostInfo.MachineId)
        } else {
            // shows usage message
            fmt.Printf(`unknown argument '%s', valid arguments are: 
- host-uuid:  the host unique identifier against a control plane; it is generated by the discovery process and is part of the activation key
- mac-addr:   MAC address of the host primary interface; the host unique identifier against the discovery service; can optionally be replaced with the device hardware id
- hw-id:      device hardware id; the alternative device unique identifier against the discovery service
- machine-id: operating system generated from a random source during system installation or first boot 
- ALL:        all host/device attributes
`, args[0])
        }
    default:
        // shows usage message
        fmt.Printf(`invalid arguments '%s'`, args[0])
    }
}

func hostUUID() string {
    core2.TRA, core2.CE = core2.NewTracer(false)
    var (
        ak  *core2.AKInfo
        err error
    )
    // prints the host UUID
    ak, err = core2.LoadActivationKey()
    if err != nil {
        core.InfoLogger.Printf("Host UUID is unknown: %s\n")
        return "unknown"
    }
    return ak.HostUUID
}
