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

// decorator for the release info cobra command
type ReleaseInfoCmd struct {
	cmd *cobra.Command
}

func NewReleaseInfoCmd() *ReleaseInfoCmd {
	c := &ReleaseInfoCmd{
		cmd: &cobra.Command{
			Use:   "info [version]",
			Short: "shows specific release information",
			Long:  ``,
		}}
	c.cmd.Run = c.run
	return c
}

func (c *ReleaseInfoCmd) run(cmd *cobra.Command, args []string) {
	fmt.Println("release info called")
}
