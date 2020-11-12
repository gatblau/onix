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
	credentials string
	tls         *bool
}

func NewPushCmd() *PushCmd {
	c := &PushCmd{
		cmd: &cobra.Command{
			Use:   "push [OPTIONS] NAME[:TAG]",
			Short: "uploads an artefact to a remote artefact store",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.tls = c.cmd.Flags().BoolP("tls", "t", true, "use -t=false or --tls=false to disable https for a backend; i.e. use plain http")
	return c
}

func (c *PushCmd) Run(cmd *cobra.Command, args []string) {
	if !*c.tls {
		log.Printf("info: Transport Level Security is disabled\n")
	}
	// check an artefact name has been provided
	if len(args) == 0 {
		log.Fatal("name of the artefact to push is required")
	}
	// get the name of the artefact to push
	nameTag := args[0]
	// validate the name
	artie := core.ParseName(nameTag)
	// create a local registry
	local := registry.NewLocalAPI()
	// attempt upload to remote repository
	local.Push(artie, registry.NewRemoteAPI(*c.tls), c.credentials)
}
