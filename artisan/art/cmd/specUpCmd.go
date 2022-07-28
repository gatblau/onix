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
	. "github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
)

// SpecUpCmd uploads the contents of a spec from a remote source
type SpecUpCmd struct {
	Cmd    *cobra.Command
	creds  string
	output string
}

func NewSpecUpCmd() *SpecUpCmd {
	c := &SpecUpCmd{
		Cmd: &cobra.Command{
			Use:   "up [FLAGS] SPEC-FOLDER",
			Short: "uploads a specification (tarball files) to a remote URI from a file system folder",
			Long: `Usage: art spec up [FLAGS] URI

Use this command to upload a specification (package export tarball files) to a remote URI form a file system folder.

Example:
   # upload the specification tarball files in my-local-folder to a remote s3s:// location 
   art spec up -o s3s://my-s3-service.com/my-app/v1.0 -c S3_ID:S3_SECRET ./my-local-folder
`,
		},
	}
	c.Cmd.Run = c.Run
	c.Cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to authenticate with the upload endpoint")
	c.Cmd.Flags().StringVarP(&c.output, "output", "o", "", "the URI where the tar archive(s) will be uploaded; URI can be s3 bucket (s3:// or s3s:// using TLS)")
	c.Cmd.MarkFlagRequired("output")
	return c
}

func (c *SpecUpCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if args != nil && len(args) < 1 {
		core.RaiseErr("the path to the local spec files is required")
	}
	err := UploadSpec(UpDownOptions{
		TargetUri:   c.output,
		TargetCreds: c.creds,
		LocalPath:   args[0],
	})
	core.CheckErr(err, "cannot upload spec")
}
