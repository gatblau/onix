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
type PushCmd struct {
	cmd         *cobra.Command
	local       *registry.FileRegistry
	remote      registry.Remote
	credentials string
}

func NewPushCmd() *PushCmd {
	c := &PushCmd{
		cmd: &cobra.Command{
			Use:   "push [OPTIONS] NAME[:TAG]",
			Short: "uploads an artefact to a remote artefact store",
			Long:  ``,
		},
		local:  registry.NewFileRegistry(),
		remote: registry.NewNexus3Registry(),
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	return c
}

func (b *PushCmd) Run(cmd *cobra.Command, args []string) {
	// check an artefact name has been provided
	if len(args) == 0 {
		log.Fatal("name of the artefact to push is required")
	}
	// get the name of the artefact to push
	nameTag := args[0]
	// validate the name
	named := core.ParseName(nameTag)
	// attempt upload to remote repository
	b.local.Push(named, b.remote, b.credentials)
}
