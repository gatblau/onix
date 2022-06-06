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

// PushCmd push a package to a remote registry
type PushCmd struct {
	cmd         *cobra.Command
	credentials string
}

func NewPushCmd() *PushCmd {
	c := &PushCmd{
		cmd: &cobra.Command{
			Use:   "push [FLAGS] NAME[:TAG]",
			Short: "uploads an package to a remote package store",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	return c
}

func (c *PushCmd) Run(cmd *cobra.Command, args []string) {
	// check an package name has been provided
	if len(args) == 0 {
		log.Fatal("name of the package to push is required")
	}
	// get the name of the package to push
	nameTag := args[0]
	// validate the name
	packageName, err := core.ParseName(nameTag)
	i18n.Err("", err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry("")
	// attempt upload to remote repository
	core.CheckErr(local.Push(packageName, c.credentials), i18n.Sprintf("", i18n.ERR_CANT_PUSH_PACKAGE))
}
