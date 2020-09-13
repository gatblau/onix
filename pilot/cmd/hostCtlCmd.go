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
			Short: "launch pilot in host mode",
			Long:  `hosts register virtual or physical hosts with the configuration database and manage running processes in the host`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *HostCtlCmd) Run(cmd *cobra.Command, args []string) {
	host, err := NewHost()
	if err != nil {
		fmt.Printf("cannot launch pilot host: %v", err)
		os.Exit(-1)
	}
	host.Start()
}
