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
	c.cmd.Flags().StringVarP(&c.cfgPath, "path", "p", "", "set the path where the configuration files are written")
	c.cmd.Flags().StringVarP(&c.cfgName, "name", "n", "", "set the name of the configuration to use; e.g. different configurations can be kept for different environments")
	return c
}

func (c *ConfigUseCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fmt.Print("oops! I do not take any arguments, use flags instead\n")
	}
	fmt.Printf("I am currently using the configuration from %v\n", DM.Cfg.ConfigFileUsed())
	if len(c.cfgPath) > 0 || len(c.cfgName) > 0 {
		DM.Use(c.cfgPath, c.cfgName)
		fmt.Printf("I have changed it to %v\n", DM.Cfg.ConfigFileUsed())
	}
}
