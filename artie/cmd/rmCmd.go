/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artie/core"
	"github.com/gatblau/onix/artie/registry"
	"github.com/spf13/cobra"
	"log"
)

// list local artefacts
type RmCmd struct {
	cmd   *cobra.Command
	local *registry.LocalAPI
}

func NewRmCmd() *RmCmd {
	c := &RmCmd{
		cmd: &cobra.Command{
			Use:   "rm ARTEFACT [ARTEFACT...]",
			Short: "removes one or more artefacts from the local artefact store",
			Long:  ``,
		},
		local: registry.NewLocalAPI(),
	}
	c.cmd.Run = c.Run
	return c
}

func (c *RmCmd) Run(cmd *cobra.Command, args []string) {
	// check one or more artefact names have been provided
	if len(args) == 0 {
		log.Fatal("missing name(s) of the artefact(s) to remove")
	}
	// get the name(s) of the artefact(s) to remove
	c.local.Remove(c.toArtURIs(args))
}

func (c *RmCmd) toArtURIs(args []string) []*core.ArtieName {
	var uris = make([]*core.ArtieName, 0)
	for _, arg := range args {
		uri := core.ParseName(arg)
		uris = append(uris, uri)
	}
	return uris
}
