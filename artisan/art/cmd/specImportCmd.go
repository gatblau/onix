/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	. "github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
)

// SpecImportCmd Import the contents from a tarball to create an artisan package in the local registry
type SpecImportCmd struct {
	cmd             *cobra.Command
	creds           string
	ignoreSignature *bool
	pubPath         string
	filter          string
}

func NewSpecImportCmd() *SpecImportCmd {
	c := &SpecImportCmd{
		cmd: &cobra.Command{
			Use:   "import [FLAGS] URI",
			Short: "imports an application release specification (e.g. one or more tarball files) into the local registry",
			Long: `Usage: art spec import [FLAGS] URI

Import one or more tarball files into the local registry using a specification (spec.yaml file).

Examples:
   # import a specification from a file system folder
   art spec import ./test
   
   # import a specification from an S3 bucket folder
   art spec import s3s://my-s3-service.com/my-app/v1.0 -c S3_ID:S3_SECRET
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to retrieve the specification from an endpoint")
	c.cmd.Flags().StringVarP(&c.filter, "filter", "f", "", "a regular expression used to select spec artefacts to be imported; any artefacts not matched by the filter are skipped (e.g. -f \"^quay.*$\")")
	return c
}

func (c *SpecImportCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if args != nil && len(args) < 1 {
		core.RaiseErr("the URI of the specification is required")
	}
	// if not ignoring signature and public key path is not provided, then uses the local registry default public key
	if len(c.pubPath) == 0 && !*c.ignoreSignature {
		// works out the FQN of the public root key
		_, pub := crypto.KeyNames(core.KeysPath(""), "root", "pgp")
		c.pubPath = pub
		fmt.Printf("verifying signatures with local root public key: %s\n", pub)
	}
	// import the tar archive(s)
	_, err := ImportSpec(ImportOptions{
		TargetUri:   args[0],
		TargetCreds: c.creds,
		Filter:      c.filter,
	})
	core.CheckErr(err, "cannot import spec")
}
