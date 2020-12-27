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
)

// list local artefacts
type PGPExportCmd struct {
	cmd       *cobra.Command
	group     string // the repository group for which the key should be used - if empty then the root is used
	name      string // the repository name for which the key should be used
	isPrivate *bool  //  whether the key is public or private
}

func NewPGPExportCmd() *PGPExportCmd {
	c := &PGPExportCmd{
		cmd: &cobra.Command{
			Use:   "export [flags] path/to/exported/file",
			Short: "export a (private or public) PGP/RSA key stored in the local registry",
			Long:  `a private PGP/RSA key is used to digitally sign an artefact upon build, a public RSA key is used to verify the digital signature when the artefact is opened`,
		},
	}
	c.isPrivate = c.cmd.Flags().BoolP("private", "k", false, "flag that indicates if the key to export is the private or public key.")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "the repository group location of the exported key.")
	c.cmd.Flags().StringVarP(&c.name, "name", "n", "", "the repository name location of the exported key.")
	c.cmd.Run = c.Run
	return c
}

func (b *PGPExportCmd) Run(cmd *cobra.Command, args []string) {
	// check if a path to the key has been provided
	if len(args) == 0 {
		core.RaiseErr("the path to the key file to be exported must be provided when calling this command")
	}
	if len(args) > 1 {
		core.RaiseErr("more than one argument have been provided, only the path to the key file is required")
	}
	l := registry.NewLocalRegistry()
	l.ExportKey(args[0], *b.isPrivate, b.group, b.name)
}
