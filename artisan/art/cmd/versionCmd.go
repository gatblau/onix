/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

// list local artefacts
type VersionCmd struct {
	cmd *cobra.Command
}

func NewVersionCmd() *VersionCmd {
	c := &VersionCmd{
		cmd: &cobra.Command{
			Use:   "version",
			Short: "show the current version",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (b *VersionCmd) Run(cmd *cobra.Command, args []string) {
	fmt.Printf("build: %s\n", core.Version)
}
