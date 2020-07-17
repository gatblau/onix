//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/dbman/core"
	"github.com/spf13/cobra"
)

type DbQueriesCmd struct {
	cmd     *cobra.Command
	format  string
	verbose bool
}

func NewDbQueriesCmd() *DbQueriesCmd {
	c := &DbQueriesCmd{
		cmd: &cobra.Command{
			Use:   "queries",
			Short: "list the available queries",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "yaml", "the format of the output - yaml, json, csv")
	c.cmd.Flags().BoolVarP(&c.verbose, "verbose", "v", false, `if true, the output will contain additional query information.`)
	return c
}

func (c *DbQueriesCmd) Run(cmd *cobra.Command, args []string) {
	// get the release manifest for the current application version
	_, manifest, err := core.DM.GetReleaseInfo(core.DM.Cfg.GetString(core.AppVersion))
	if err != nil {
		fmt.Printf("!!! I cannot fetch release information: %v\n", err)
		return
	}
	fmt.Println(manifest.GetQueriesInfo(c.format, c.verbose))
}
