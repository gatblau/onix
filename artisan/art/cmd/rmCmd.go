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
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
)

// list local packages
type RmCmd struct {
	cmd *cobra.Command
}

func NewRmCmd() *RmCmd {
	c := &RmCmd{
		cmd: &cobra.Command{
			Use:   "rm PACHAGE [PACKAGE...]",
			Short: "removes one or more packages from the local package registry",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *RmCmd) Run(cmd *cobra.Command, args []string) {
	// check one or more package names have been provided
	if len(args) == 0 {
		log.Fatal("missing name(s) of the package(s) to remove")
	}
	//  create a local registry
	local := registry.NewLocalRegistry()
	// get the name(s) of the package(s) to remove
	local.Remove(c.toArtURIs(args))
}

func (c *RmCmd) toArtURIs(args []string) []*core.PackageName {
	var uris = make([]*core.PackageName, 0)
	for _, arg := range args {
		uri, err := core.ParseName(arg)
		i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
		uris = append(uris, uri)
	}
	return uris
}
