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

type SideCarCmd struct {
	cmd *cobra.Command
}

func NewSideCarCmd() *SideCarCmd {
	c := &SideCarCmd{
		cmd: &cobra.Command{
			Use:   "sidecar",
			Short: "launch pilot in sidecar mode",
			Long:  `sidecars synchronise configuration files for applications following changes in the configuration database and protect local configurations from being updated by unauthorised parties`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *SideCarCmd) Run(cmd *cobra.Command, args []string) {
	sidecar, err := NewSidecar()
	if err != nil {
		fmt.Printf("cannot launch pilot sidecar: %v", err)
		os.Exit(-1)
	}
	sidecar.Start()
}
