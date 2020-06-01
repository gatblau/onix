//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	. "github.com/gatblau/onix/dbman/util"
	"github.com/spf13/cobra"
)

type CheckCmd struct {
	cmd *cobra.Command
}

func NewCheckCmd() *CheckCmd {
	c := &CheckCmd{
		cmd: &cobra.Command{
			Use:   "check",
			Short: "performs a health check of dbman's current configuration",
			Long:  ``,
		}}
	c.cmd.Run = c.Run
	return c
}

func (c *CheckCmd) Run(cmd *cobra.Command, args []string) {
	DM.Check()
}
