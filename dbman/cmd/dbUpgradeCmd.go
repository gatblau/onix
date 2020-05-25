//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type DbUpgradeCmd struct {
	cmd *cobra.Command
}

func NewDbUpgradeCmd() *DbUpgradeCmd {
	c := &DbUpgradeCmd{
		&cobra.Command{
			Use:   "upgrade [version]",
			Short: "upgrades the current schema to a specific version",
			Long:  `if version is not specified, then rolling upgrades to the latest version are executed`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DbUpgradeCmd) Run(cmd *cobra.Command, args []string) {
	fmt.Println("restore called")
}
