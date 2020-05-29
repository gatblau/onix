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

type ConfigUseCmd struct {
	cmd     *cobra.Command
	cfgPath string
	cfgName string
}

func NewConfigUseCmd() *ConfigUseCmd {
	c := &ConfigUseCmd{
		cmd: &cobra.Command{
			Use:     "use",
			Short:   "switches to the specified configuration file",
			Example: `dbman config use myapp_dev`,
		}}
	c.cmd.Run = c.Run
	return c
}

func (c *ConfigUseCmd) Run(cmd *cobra.Command, args []string) {
	DM.Use(cfgPath, cfgName)
	fmt.Printf("using configuration from %v\n", DM.Cfg.ConfigFileUsed())
}
