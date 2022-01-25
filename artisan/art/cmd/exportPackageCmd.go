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

// ExportPackageCmd save one or more packages from the local registry to a tar archive to allow copying without using registries
type ExportPackageCmd struct {
	cmd         *cobra.Command
	srcCreds    string
	targetCreds string
	output      string
}

func NewSavePackageCmd() *ExportPackageCmd {
	c := &ExportPackageCmd{
		cmd: &cobra.Command{
			Use:   "package [FLAGS] PACKAGE [PACKAGE...]",
			Short: "export one or more packages to a tar archive",
			Long: `Usage: art export package [FLAGS] PACKAGE [PACKAGE...]

Exports one or more packages to a tar archive, streamed to STDOUT by default or to a URI that can be for the file system or an S3 endpoint

Examples:
   # exports my-package-1 and my-package-2 to a tar archive by redirecting STDOUT to file (using the redirection operator '>')
   art export package my-package-1 my-package-2 > archive.tar 
   
   # exports my-package-1 and my-package-2 to a tar archive by specifying relative file path via URI (using the -o flag)
   art export package my-package-1 my-package-2 -o ./test/archive.tar 

   # exports my-package-1 and my-package-2 from remote artisan registry to an authenticated and TLS enabled s3 bucket
   art export package my-package-1 my-package-2 -u REG_USER:REG_PWD -o s3s://endpoint/bucket/archive.tar -c S3_ID:S3_SECRET
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.output, "output", "o", "", "the URI where the tar archive will be saved, instead of STDOUT; URI can be file system (absolute or relative path) or s3 bucket (s3:// or s3s:// using TLS)")
	c.cmd.Flags().StringVarP(&c.srcCreds, "user", "u", "", "the credentials used to pull packages from an authenticated artisan registry, if the packages are not already in the local registry")
	c.cmd.Flags().StringVarP(&c.targetCreds, "creds", "c", "", "the credentials to write packages to a destination, if such destination implements authentication (e.g. s3)")
	return c
}

func (c *ExportPackageCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if len(args) < 1 {
		log.Fatal("at least the name of one package to save is required")
	}
	// validate the package names
	names, err := core.ValidateNames(args)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry()
	// export packages into tar bytes
	err = local.Save(names, c.srcCreds, c.output, c.targetCreds)
	core.CheckErr(err, "cannot export package(s)")
}
