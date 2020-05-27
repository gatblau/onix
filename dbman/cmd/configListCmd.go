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

type ConfigListCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewConfigListCmd() *ConfigListCmd {
	c := &ConfigListCmd{
		cmd: &cobra.Command{
			Use:     "list",
			Short:   "list dbman's configuration values",
			Example: `dbman config list`,
		}}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml or json")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a filename is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *ConfigListCmd) Run(cmd *cobra.Command, args []string) {
	DM.PrintConfig()
}
