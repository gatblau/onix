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

type TagCmd struct {
	Cmd *cobra.Command
}

func NewTagCmd() *TagCmd {
	c := &TagCmd{
		Cmd: &cobra.Command{
			Use:     "tag",
			Short:   "add a tag to an existing package",
			Long:    `create a tag TARGET_PACKAGE that refers to SOURCE_PACKAGE`,
			Example: `art tag SOURCE_PACKAGE[:TAG] TARGET_PACKAGE[:TAG]`,
		}}
	c.Cmd.Run = c.Run
	return c
}

func (c *TagCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		core.RaiseErr("source and target package tags are required")
	}
	l := registry.NewLocalRegistry("")
	core.CheckErr(l.Tag(args[0], args[1]), "cannot tag package")
}
