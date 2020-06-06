//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	. "github.com/gatblau/onix/dbman/util"
	"github.com/spf13/cobra"
)

type DbInitCmd struct {
	cmd *cobra.Command
}

func NewDbInitCmd() *DbInitCmd {
	c := &DbInitCmd{
		&cobra.Command{
			Use:   "init",
			Short: "initialise the database",
			Long: `execute admin level scripts to create the database, database user, etc in advanced of the creation 
of the schema and database objects`,
			Example: `dbman db init`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DbInitCmd) Run(cmd *cobra.Command, args []string) {
	err := DM.InitialiseDb()
	if err != nil {
		fmt.Printf("!!! I cannot initialise the database: %v", err)
	} else {
		fmt.Printf("? I have completed the database initialisation")
	}
}
