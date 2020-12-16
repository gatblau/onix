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
type OpenCmd struct {
	cmd         *cobra.Command
	credentials string
	tls         *bool
	verify      *bool
	path        string
	pubPath     string
}

func NewOpenCmd() *OpenCmd {
	c := &OpenCmd{
		cmd: &cobra.Command{
			Use:   "open NAME[:TAG] [path]",
			Short: "opens an artefact in the specified path",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.tls = c.cmd.Flags().BoolP("tls", "t", true, "-t=false or --tls=false to disable https for a backend; i.e. use plain http")
	c.verify = c.cmd.Flags().BoolP("verify", "v", true, "-v=false or --verify=false to signature verification")
	c.cmd.Flags().StringVarP(&c.pubPath, "pub", "p", "", "--pub=/path/to/public/key or -p=/path/to/public/key")
	return c
}

func (c *OpenCmd) Run(cmd *cobra.Command, args []string) {
	if !*c.tls {
		log.Printf("info: Transport Level Security is disabled\n")
	}
	// check an artefact name has been provided
	if len(args) < 1 {
		log.Fatal("name of the artefact to open is required")
	}
	// get the name of the artefact to push
	nameTag := args[0]
	path := ""
	if len(args) == 2 {
		path = args[1]
	}
	// validate the name
	artie := core.ParseName(nameTag)
	// create a local registry
	local := registry.NewLocalRegistry()
	// attempt to open from local registry
	local.Open(artie, c.credentials, *c.tls, path, c.pubPath, *c.verify)
}
