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
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
)

// list local packages
type OpenCmd struct {
	cmd             *cobra.Command
	credentials     string
	noTLS           *bool
	ignoreSignature *bool
	path            string
	pubPath         string
}

func NewOpenCmd() *OpenCmd {
	c := &OpenCmd{
		cmd: &cobra.Command{
			Use:   "open NAME[:TAG] [path]",
			Short: "opens an package in the specified path",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.noTLS = c.cmd.Flags().BoolP("no-tls", "t", false, "use -t or --no-tls to connect to a artisan registry over plain HTTP")
	c.ignoreSignature = c.cmd.Flags().BoolP("ignore-sig", "s", false, "-s or --ignore-sig to ignore signature verification")
	c.cmd.Flags().StringVarP(&c.pubPath, "pub", "p", "", "-p=/path/to/public/key or --pub=/path/to/public/key to load a public PGP key to verify the package digital signature")
	return c
}

func (c *OpenCmd) Run(cmd *cobra.Command, args []string) {
	if *c.noTLS {
		log.Printf("info: Transport Level Security is disabled\n")
	}
	// check an package name has been provided
	if len(args) < 1 {
		log.Fatal("name of the package to open is required")
	}
	// get the name of the package to push
	nameTag := args[0]
	path := ""
	if len(args) == 2 {
		path = args[1]
	}
	// validate the name
	artie, err := core.ParseName(nameTag)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry()
	// attempt to open from local registry
	local.Open(artie, c.credentials, *c.noTLS, path, c.pubPath, *c.ignoreSignature)
}
