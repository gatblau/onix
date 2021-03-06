/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

// merges environment variables into one or more files
type MergeCmd struct {
	cmd         *cobra.Command
	envFilename string
}

func NewMergeCmd() *MergeCmd {
	c := &MergeCmd{
		cmd: &cobra.Command{
			Use:   "merge [flags] [files]",
			Short: "merges environment variables in the specified files",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
	return c
}

func (c *MergeCmd) Run(cmd *cobra.Command, args []string) {
	core.MergeFilesOld(args, c.envFilename)
}
