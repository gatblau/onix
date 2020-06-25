//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/dbman/plugins"
	. "github.com/gatblau/onix/dbman/util"
	"github.com/spf13/cobra"
)

type DbVersionCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewDbVersionCmd() *DbVersionCmd {
	c := &DbVersionCmd{
		cmd: &cobra.Command{
			Use:   "version",
			Short: "retrieves database version information",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml, json, csv")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a filename is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *DbVersionCmd) Run(cmd *cobra.Command, args []string) {
	resultStr := DM.DbPlugin().GetVersion()
	r := plugins.NewParameterFromJSON(resultStr)
	if r.HasError() {
		fmt.Printf("!!! I cannot retrieve database version: '%s'", r.Error())
		return
	}
	// if a filename was specified then save the content to file
	if len(c.filename) > 0 {
		r.Save(c.format, c.filename)
	} else {
		// otherwise print to the console
		fmt.Println(r.Sprint(c.format))
	}
}
