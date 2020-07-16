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

type DbDeployCmd struct {
	cmd *cobra.Command
}

func NewDbDeployCmd() *DbDeployCmd {
	c := &DbDeployCmd{
		cmd: &cobra.Command{
			Use:   "deploy",
			Short: "deploy schemas and objects in an existing database based on the current Application Version",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DbDeployCmd) Run(cmd *cobra.Command, args []string) {
	output, err, elapsed := DM.Deploy()
	fmt.Print(output.String())
	if err != nil {
		fmt.Printf("!!! I cannot deploy the database\n")
		fmt.Printf("%v\n", err)
		fmt.Printf("? the execution time was %v\n", elapsed)
		return
	}
	fmt.Printf("? I have deployed the database in %v\n", elapsed)
}
