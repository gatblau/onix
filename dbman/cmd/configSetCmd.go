//   Onix Config Db - Dbman
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

type ConfigSetCmd struct {
	cmd *cobra.Command
	key string
}

func NewConfigSetCmd() *ConfigSetCmd {
	c := &ConfigSetCmd{
		cmd: &cobra.Command{
			Use:     "set [key] [value]",
			Short:   "set the specified configuration value",
			Example: `dbman config set SchemaURI https://raw.githubusercontent.com/gatblau/ox-db/master`,
		}}
	c.cmd.Run = c.Run
	return c
}

func (c *ConfigSetCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Printf("!!! I found an incorrect number of arguments and cannot process the command.\n" +
			"??? You need to pass [key] and [value] as arguments to the set command.\n")
		return
	}
	key := args[0]
	value := args[1]
	DM.SetConfig(key, value)
	DM.SaveConfig()
}
