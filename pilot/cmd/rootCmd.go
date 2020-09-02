/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type RootCmd struct {
	*cobra.Command
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		&cobra.Command{
			Use:   "pilot",
			Short: "configuration manager host agent",
			Long:  `An Onix agent to synchronise configuration changes with the processes running on a host or pod and vice-versa`,
		},
	}
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootCmd) initConfig() {
}
