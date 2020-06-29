//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	. "github.com/gatblau/onix/dbman/core"
	"github.com/spf13/cobra"
	"os"
)

type ConfigDeleteCmd struct {
	cmd *cobra.Command
}

func NewConfigDeleteCmd() *ConfigDeleteCmd {
	c := &ConfigDeleteCmd{
		cmd: &cobra.Command{
			Use:     "delete [set name]",
			Short:   "delete the configuration set specified by its name",
			Example: `dbman config delete dev`,
		}}
	c.cmd.Run = c.Run
	return c
}

func (c *ConfigDeleteCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("!!! I need to know the configuration set name")
		return
	}
	err := os.Remove(fmt.Sprintf("%v/.dbman_%v.toml", DM.GetConfigSetDir(), args[0]))

	if err != nil {
		fmt.Printf("!!! I could not remove configuration %v: %v\n", args[0], err)
		return
	}
}
