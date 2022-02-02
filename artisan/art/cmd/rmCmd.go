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
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
)

// RmCmd remove local packages
type RmCmd struct {
	cmd *cobra.Command
	all *bool
}

func NewRmCmd() *RmCmd {
	c := &RmCmd{
		cmd: &cobra.Command{
			Use:   "rm PACKAGE [PACKAGE...]",
			Short: "removes one or more packages from the local package registry",
			Long:  `removes one or more packages from the local package registry`,
		},
	}
	c.all = c.cmd.Flags().BoolP("all", "a", false, "remove all packages")
	c.cmd.Run = c.Run
	return c
}

func (c *RmCmd) Run(cmd *cobra.Command, args []string) {
	// check one or more package names have been provided if remove all is not specified
	if len(args) == 0 && !*c.all {
		core.RaiseErr("missing name(s) of the package(s) to remove")
	}
	// cannot provide all flag and package name
	if len(args) > 0 && *c.all {
		core.RaiseErr("a package name %s should not be provided with the --all|-a flag", args[0])
	}
	//  create a local registry
	local := registry.NewLocalRegistry()
	if *c.all {
		core.CheckErr(local.Remove(local.AllPackages()), "cannot remove packages")
	} else {
		core.CheckErr(local.Remove(args), "cannot remove package")
	}
}
