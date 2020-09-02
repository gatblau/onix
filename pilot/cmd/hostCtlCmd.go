/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	. "gatblau.org/onix/pilot/core"
	"github.com/spf13/cobra"
	"os"
)

type HostCtlCmd struct {
	cmd *cobra.Command
}

func NewHostCtlCmd() *HostCtlCmd {
	c := &HostCtlCmd{
		cmd: &cobra.Command{
			Use:   "host",
			Short: "launches pilot in Host Control mode",
			Long:  `pilot registers the host, listens for configuration manager events and syncs configuration data for applications running on the host`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *HostCtlCmd) Run(cmd *cobra.Command, args []string) {
	pilot, err := NewPilot()
	if err != nil {
		fmt.Printf("cannot launch pilot: %v", err)
		os.Exit(-1)
	}
	P = pilot
	P.Host()
}
