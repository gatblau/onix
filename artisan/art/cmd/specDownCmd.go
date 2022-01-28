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
)

// SpecDownCmd downloads the contents of a spec from a remote source
type SpecDownCmd struct {
	cmd       *cobra.Command
	creds     string
	localPath string
}

func NewSpecDownCmd() *SpecDownCmd {
	c := &SpecDownCmd{
		cmd: &cobra.Command{
			Use: "down [FLAGS] URI",
			Short: "downloads a specification (tarball files) from a remote URI to a file system folder but does not " +
				"actually perform the import",
			Long: `Usage: art import spec-download [FLAGS] URI

Use this command to download a specification (package export tarball files) from a remote URI to a file system folder 
but does not actually perform the import.
This is useful when files have to be inspected or transferred before the can be imported.

Example:
   # download a specification from a remote location
   art spec down s3s://my-s3-service.com/my-app/v1.0 -c S3_ID:S3_SECRET -p ./my-local-folder
 
   # perform some file validation / scan of downloaded assets
   scan [DEVICE]

   # import the specification from the local folder 
   art spec import ./my-local-folder
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to retrieve the specification from an endpoint")
	c.cmd.Flags().StringVarP(&c.localPath, "path", "p", "", "if specified, download the spec tarball files to the path")
	return c
}

func (c *SpecDownCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if args != nil && len(args) < 1 {
		core.RaiseErr("the URI of the specification is required")
	}
	// import the tar archive(s)
	err := export.DownloadSpec(args[0], c.creds, c.localPath)
	core.CheckErr(err, "cannot download spec")
}
