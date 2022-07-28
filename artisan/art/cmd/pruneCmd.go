/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
)

// PruneCmd remove all dangling packages
type PruneCmd struct {
	Cmd *cobra.Command
}

func NewPruneCmd() *PruneCmd {
	c := &PruneCmd{
		Cmd: &cobra.Command{
			Use:   "prune",
			Short: "remove all dangling packages",
			Long:  `remove all dangling packages`,
		},
	}
	c.Cmd.Run = c.Run
	return c
}

func (b *PruneCmd) Run(cmd *cobra.Command, args []string) {
	local := registry.NewLocalRegistry("")
	core.CheckErr(local.Prune(), "")
}
