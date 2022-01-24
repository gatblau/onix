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
	"github.com/spf13/cobra"
	"log"
)

// SaveImageCmd save one or more container images to a tar archive
type SaveImageCmd struct {
	cmd         *cobra.Command
	targetCreds string
	output      string
	packageName string
}

func NewSaveImageCmd() *SaveImageCmd {
	c := &SaveImageCmd{
		cmd: &cobra.Command{
			Use:   "image [FLAGS] IMAGE",
			Short: "save a container image as an artisan package",
			Long: `Usage: art save image [FLAGS] IMAGE 

Save a container image as an artisan package
If a target URI is specified, the package is save to the target, otherwise it remains in the local artisan registry
Note the container images must be in the local container registry (already pulled)

Examples:
   # create a tar archive of my-contaner-image and put it in an artisan package
   art save image my-contaner-image -t my-package-name
   
   # create a tar archive of my-contaner-image and put it in an artisan package 
   # then exports the package as a tar archive to the specified path
   art save image my-contaner-image -t my-package-name -o ./test/archive.tar 

   # create a tar archive of my-contaner-image and put it in an artisan package 
   # then exports the package as a tar archive to the specified s3 location 
   art save image my-contaner-image -t my-package-name -o s3s://endpoint/bucket/archive.tar -c S3_ID:S3_SECRET
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.output, "output", "o", "", "the URI where the tar archive will be saved, instead of STDOUT; URI can be file system (absolute or relative path) or s3 bucket (s3:// or s3s:// using TLS)")
	c.cmd.Flags().StringVarP(&c.targetCreds, "creds", "c", "", "the credentials to write packages to a destination, if such destination implements authentication (e.g. s3)")
	c.cmd.Flags().StringVarP(&c.packageName, "package-name", "t", "", "the name of the package to create")
	return c
}

func (c *SaveImageCmd) Run(cmd *cobra.Command, args []string) {
	// check an image name has been provided
	if len(args) < 1 {
		log.Fatal("at least the name of one images to save is required")
	}
	core.CheckErr(export.SaveImage(args[0], c.packageName, c.output, c.targetCreds), "cannot save image")
}
