//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// decorator for the release plan cobra command
type ReleasePlanCmd struct {
	cmd *cobra.Command
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
	return c
}

func (c *ReleasePlanCmd) run(cmd *cobra.Command, args []string) {
	fmt.Println("release plan called")
}
