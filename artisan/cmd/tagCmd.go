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
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
)

type TagCmd struct {
	cmd *cobra.Command
}

func NewTagCmd() *TagCmd {
	c := &TagCmd{
		cmd: &cobra.Command{
			Use:     "tag",
			Short:   "add a tag to an existing artefact",
			Long:    `create a tag TARGET_ARTEFACT that refers to SOURCE_ARTEFACT`,
			Example: `art tag SOURCE_ARTEFACT[:TAG] TARGET_ARTEFACT[:TAG]`,
		}}
	c.cmd.Run = c.Run
	return c
}

func (c *TagCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		core.RaiseErr("source and target artefact tags are required")
	}
	l := registry.NewLocalRegistry()
	l.Tag(core.ParseName(args[0]), core.ParseName(args[1]))
}
