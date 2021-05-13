package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
)

type PullCmd struct {
	cmd         *cobra.Command
	credentials string
	path        string
}

func NewPullCmd() *PullCmd {
	c := &PullCmd{
		cmd: &cobra.Command{
			Use:   "pull [OPTIONS] NAME[:TAG]",
			Short: "downloads an package from the package registry",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	return c
}

func (c *PullCmd) Run(cmd *cobra.Command, args []string) {
	// check an package name has been provided
	if len(args) == 0 {
		log.Fatal("name of the package to pull is required")
	}
	// get the name of the package to push
	nameTag := args[0]
	// validate the name
	packageName, err := core.ParseName(nameTag)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry()
	// attempt pull from remote registry
	local.Pull(packageName, c.credentials)
}
