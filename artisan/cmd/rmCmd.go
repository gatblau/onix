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
	"log"
)

// list local artefacts
type RmCmd struct {
	cmd *cobra.Command
}

func NewRmCmd() *RmCmd {
	c := &RmCmd{
		cmd: &cobra.Command{
			Use:   "rm ARTEFACT [ARTEFACT...]",
			Short: "removes one or more artefacts from the local artefact store",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *RmCmd) Run(cmd *cobra.Command, args []string) {
	// check one or more artefact names have been provided
	if len(args) == 0 {
		log.Fatal("missing name(s) of the artefact(s) to remove")
	}
	//  create a local registry
	local := registry.NewLocalRegistry()
	// get the name(s) of the artefact(s) to remove
	local.Remove(c.toArtURIs(args))
}

func (c *RmCmd) toArtURIs(args []string) []*core.PackageName {
	var uris = make([]*core.PackageName, 0)
	for _, arg := range args {
		uri := core.ParseName(arg)
		uris = append(uris, uri)
	}
	return uris
}
