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
)

// list local packages
type PGPImportCmd struct {
	cmd       *cobra.Command
	group     string // the repository group for which the key should be used - if empty then the root is used
	name      string // the repository name for which the key should be used
	isPrivate *bool  //  whether the key is public or private
}

func NewPGPImportCmd() *PGPImportCmd {
	c := &PGPImportCmd{
		cmd: &cobra.Command{
			Use:   "import [flags] path/to/key/file",
			Short: "import PGP/RSA keys into the local registry",
			Long:  `a private PGP/RSA key is used to digitally sign an package upon build, a public RSA key is used to verify the digital signature when the package is opened`,
		},
	}
	c.isPrivate = c.cmd.Flags().BoolP("private", "k", false, "A flag indicating if the key to import is the private or public key.")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "The repository group to which the key should be applied. \nIf not specified, the key is applied to the registry root and it is used for all repositories.")
	c.cmd.Flags().StringVarP(&c.name, "name", "n", "", "The repository name to which the key should be applied. \nIf not specified, the key is applied to the repository group and it is used for all repositories under the group.")
	c.cmd.Run = c.Run
	return c
}

func (b *PGPImportCmd) Run(cmd *cobra.Command, args []string) {
	// check if a path to the key has been provided
	if len(args) == 0 {
		core.RaiseErr("the path to the key file to be imported must be provided when calling this command")
	}
	if len(args) > 1 {
		core.RaiseErr("more than one argument have been provided, only the path to the key file is required")
	}
	l := registry.NewLocalRegistry()
	l.ImportKey(args[0], *b.isPrivate, b.group, b.name)
}
