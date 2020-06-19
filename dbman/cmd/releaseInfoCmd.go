//   Onix Config DatabaseProvider - Dbman
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

// decorator for the release info cobra command
type ReleaseInfoCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewReleaseInfoCmd() *ReleaseInfoCmd {
	c := &ReleaseInfoCmd{
		cmd: &cobra.Command{
			Use:   "info [app version]",
			Short: "shows the release manifest that matches an application version",
			Long: `the release manifest contains a list of scripts to be applied to the database to make it compatible 
with the application`,
		},
	}
	c.cmd.Run = c.run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml or json")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a filename is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *ReleaseInfoCmd) run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("!!! I need a release tag to be provided, see help: dbman release info --help\n")
		return
	}
	// get the app version from the first argument
	appVer := args[0]
	// get the release manifest that matches the app version
	release, err := DM.GetReleaseInfo(appVer)
	if err != nil {
		fmt.Sprintf("!!! I cannot fetch release information: %v", err)
		return
	}
	Print(release, c.format, c.filename)
}
