/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
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
type PullCmd struct {
	cmd         *cobra.Command
	credentials string
	noTLS       *bool
	path        string
}

func NewPullCmd() *PullCmd {
	c := &PullCmd{
		cmd: &cobra.Command{
			Use:   "pull [OPTIONS] NAME[:TAG]",
			Short: "downloads an artefact from the artefact registry",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.noTLS = c.cmd.Flags().BoolP("no-tls", "t", false, "use -t or --no-tls to connect to a artisan registry over plain HTTP")
	return c
}

func (c *PullCmd) Run(cmd *cobra.Command, args []string) {
	if !*c.noTLS {
		log.Printf("info: Transport Level Security is disabled\n")
	}
	// check an artefact name has been provided
	if len(args) == 0 {
		log.Fatal("name of the artefact to pull is required")
	}
	// get the name of the artefact to push
	nameTag := args[0]
	// validate the name
	artie := core.ParseName(nameTag)
	// create a local registry
	local := registry.NewLocalRegistry()
	// attempt pull from remote registry
	local.Pull(artie, c.credentials, *c.noTLS)
}
