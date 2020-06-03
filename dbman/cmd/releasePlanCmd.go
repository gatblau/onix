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

// decorator for the release plan cobra command
type ReleasePlanCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewReleasePlanCmd() *ReleasePlanCmd {
	c := &ReleasePlanCmd{
		cmd: &cobra.Command{
			Use:   "plan",
			Short: "displays the release plan",
			Long:  `A release plan is the list of all releases available`,
		},
	}
	c.cmd.Run = c.run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml or json")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a file name is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *ReleasePlanCmd) run(cmd *cobra.Command, args []string) {
	// fetch the release plan
	plan, err := DM.GetReleasePlan()
	if err != nil {
		fmt.Printf("!!! cannot get release plan: %v", err)
		return
	}
	Print(plan, c.format, c.filename)
}
