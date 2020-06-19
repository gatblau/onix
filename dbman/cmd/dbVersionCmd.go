//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	// "fmt"
	// . "github.com/gatblau/onix/dbman/util"
	"github.com/spf13/cobra"
)

type DbVersionCmd struct {
	cmd *cobra.Command
}

func NewDbVersionCmd() *DbVersionCmd {
	c := &DbVersionCmd{
		&cobra.Command{
			Use:   "version",
			Short: "retrieves database version information",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DbVersionCmd) Run(cmd *cobra.Command, args []string) {
	// version, err := DM.RunQuery("")
}
