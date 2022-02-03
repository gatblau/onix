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
	"github.com/gatblau/onix/artisan/export"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
)

// SpecSignCmd re-signs an existing package
type SpecSignCmd struct {
	cmd     *cobra.Command
	pkPath  string
	pubPath string
}

func NewSpecSignCmd() *SpecSignCmd {
	c := &SpecSignCmd{
		cmd: &cobra.Command{
			Use:   "sign [OPTIONS] SPEC-FILE-PATH",
			Short: "re-signs an existing package",
			Long: `re-signs an existing package
Usage: art spec sign [OPTIONS] SPEC-FILE-PATH

Use this command to apply a new digital signature to packages in a specification.
If the path to the spec.yaml file is not specified, the current folder is assumed.

Example:
   # re-signs all packages in the spec.yaml in the current folder using the passed-in private pgp key
   art spec sign -k new-private-key.pgp .
 
   # note: to create a pgp key pair see 'art pgp gen -h'
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.pkPath, "pk", "k", "", "-k=/path/to/private/key or --key=/path/to/private/key to use a private PGP key to re-create the package digital signature; if not specified, the key deployed in the local registry is used")
	c.cmd.Flags().StringVarP(&c.pubPath, "pub", "p", "", "-p=/path/to/public/key or --pub=/path/to/public/key to use a public PGP key to verify the original package digital signature; if not specified, no signature verification of the original package is performed")
	return c
}

func (c *SpecSignCmd) Run(cmd *cobra.Command, args []string) {
	// if no path to the spec.yaml has been provided
	if args == nil || len(args) == 0 {
		// assume current path
		args = []string{"."}
	}
	core.CheckErr(signSpec(args[0], c.pkPath, c.pubPath), "cannot sign spec artefacts")
}

func signSpec(specPath, pk, pub string) error {
	spec, err := export.NewSpec(specPath, "")
	if err != nil {
		return err
	}
	local := registry.NewLocalRegistry()
	for _, pac := range spec.Packages {
		err = local.Sign(pac, pk, pub)
		if err != nil {
			return err
		}
	}
	return nil
}
