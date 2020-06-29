//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	. "github.com/gatblau/onix/dbman/core"
	"github.com/spf13/cobra"
)

type DbCreateCmd struct {
	cmd *cobra.Command
}

func NewDbCreateCmd() *DbCreateCmd {
	c := &DbCreateCmd{
		cmd: &cobra.Command{
			Use:   "create",
			Short: "create a new database based on the current Application Version",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DbCreateCmd) Run(cmd *cobra.Command, args []string) {
	output, err, elapsed := DM.Create()
	fmt.Print(output.String())
	if err != nil {
		return
	}
	fmt.Printf("? I have created the database in %v\n", elapsed)
}
