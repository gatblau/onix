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
	"strings"
)

// RmCmd remove local packages
type RmCmd struct {
	cmd      *cobra.Command
	all      bool
	registry bool
	creds    string
}

func NewRmCmd() *RmCmd {
	c := &RmCmd{
		cmd: &cobra.Command{
			Use:   "rm PACKAGE [PACKAGE...]",
			Short: "removes one or more packages from the local package registry or a remote registry",
			Long:  `removes one or more packages from the local package registry or a remote registry`,
			Example: `
# delete all packages in local registry
art rm -a

# delete specific packages from the local registry
art rm package-1 package-2

# delete all packages in the remote registry at localhost:8082
# note the use of the wildcard "localhost:8082/*" - any other valid regular expressions are acceptable
art rm -ru <user>:<pwd> "localhost:8082/*"
`,
		},
	}
	c.cmd.Flags().BoolVarP(&c.all, "all", "a", false, "remove all packages")
	c.cmd.Flags().BoolVarP(&c.registry, "registry", "r", false, "indicates to use a remote artisan registry")
	c.cmd.Flags().StringVarP(&c.creds, "user", "u", "", "the credentials used to retrieve the information from the remote registry")
	c.cmd.Run = c.Run
	return c
}

func (c *RmCmd) Run(cmd *cobra.Command, args []string) {
	// check one or more package names have been provided if remove all is not specified
	if len(args) == 0 && !c.all {
		core.RaiseErr("missing name(s) of the package(s) to remove")
	}
	// cannot provide all flag and package name
	if len(args) > 0 && c.all {
		core.RaiseErr("a package name %s should not be provided with the --all|-a flag", args[0])
	}
	// if no remote specified then it is a local operation
	if !c.registry {
		//  create a local registry
		local := registry.NewLocalRegistry("")
		if c.all {
			// prune dangling packages first
			core.CheckErr(local.Prune(), "cannot prune packages")
			// remove all packages
			core.CheckErr(local.Remove(local.AllPackages()), "cannot remove packages")
		} else {
			core.CheckErr(local.Remove(args), "cannot remove package")
		}
	} else {
		uname, pwd := core.RegUserPwd(c.creds)
		remote, err := registry.NewRemoteRegistry(domain(args[0]), uname, pwd, "")
		core.CheckErr(err, "invalid remote")
		// otherwise, it is a remote operation
		if c.all {
			core.RaiseErr("--all flag is not valid for remote registries, use a filter expression instead")
		}
		core.CheckErr(remote.Remove(args[0]), "cannot remove packages from remote registry")
	}
}

func domain(name string) string {
	if strings.Contains(name, "//") {
		core.RaiseErr("invalid host name")
	}
	return strings.Split(name, "/")[0]
}
