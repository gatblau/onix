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
type PushCmd struct {
	cmd         *cobra.Command
	credentials string
	noTLS       *bool
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
	c.noTLS = c.cmd.Flags().BoolP("no-tls", "t", false, "use -t or --no-tls to connect to a artisan registry over plain HTTP")
	return c
}

func (c *PushCmd) Run(cmd *cobra.Command, args []string) {
	if *c.noTLS {
		log.Printf("info: Transport Level Security is disabled\n")
	}
	// check an artefact name has been provided
	if len(args) == 0 {
		log.Fatal("name of the artefact to push is required")
	}
	// get the name of the artefact to push
	nameTag := args[0]
	// validate the name
	packageName, err := core.ParseName(nameTag)
	core.CheckErr(err, "invalid package name")
	// create a local registry
	local := registry.NewLocalRegistry()
	// attempt upload to remote repository
	local.Push(packageName, c.credentials, *c.noTLS)
}
