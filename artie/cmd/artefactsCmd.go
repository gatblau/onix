/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artie/registry"
	"github.com/spf13/cobra"
)

// list local artefacts
type ArtefactsCmd struct {
	cmd   *cobra.Command
	quiet *bool
}

func NewArtefactsCmd() *ArtefactsCmd {
	c := &ArtefactsCmd{
		cmd: &cobra.Command{
			Use:   "artefacts",
			Short: "list artefacts",
			Long:  ``,
		},
	}
	c.quiet = c.cmd.Flags().BoolP("quiet", "q", false, "only show numeric IDs")
	c.cmd.Run = c.Run
	return c
}

func (b *ArtefactsCmd) Run(cmd *cobra.Command, args []string) {
	local := registry.NewLocalAPI()
	if *b.quiet {
		local.ListQ()
	} else {
		local.List()
	}
}
