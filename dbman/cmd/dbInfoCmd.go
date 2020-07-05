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

type DbInfoCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewDbInfoCmd() *DbInfoCmd {
	c := &DbInfoCmd{
		cmd: &cobra.Command{
			Use:   "server",
			Short: "gets database server information",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml, json, csv")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a filename is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *DbInfoCmd) Run(cmd *cobra.Command, args []string) {
	info, err := DM.GetDbInfo()
	if err != nil {
		fmt.Printf("!!! I cannot get database server information: %v\n", err)
		return
	}
	Print(info, c.format, c.filename)
}
