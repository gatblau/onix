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

// SpecPullCmd pulls all artefacts defined in a spec
type SpecPullCmd struct {
	cmd   *cobra.Command
	creds string
	user  string
}

func NewSpecPullCmd() *SpecPullCmd {
	c := &SpecPullCmd{
		cmd: &cobra.Command{
			Use:   "pull [FLAGS] URI",
			Short: "pull all artefacts in a specification to the localhost ",
			Long: `Usage: art spec pull [FLAGS] URI

Use this command to pull all artefacts such as packages and images required to be exported.

Example:
   # in the simplest case, pull all images in a locally stored spec file 
   # assume spec contains no packages hence no need for using -u flag
   art spec pull .

   # pull artefacts defined in a specification located at an S3 bucket
   art spec pull s3s://my-s3-service.com/my-app/v1.0 -c S3_ID:S3_SECRET -u reg_user:reg_pwd

   # pull artefacts defined in a specification located in the local file system
   art spec pull ./my-app/v1.0 -u reg_user:reg_pwd
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to retrieve the specification from an endpoint")
	c.cmd.Flags().StringVarP(&c.creds, "user", "u", "", "the credentials used to retrieve packages from a remote artisan registry; for container images you should already be logged in (e.g. docker login)")
	return c
}

func (c *SpecPullCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if args != nil && len(args) < 1 {
		core.RaiseErr("the URI of the specification is required")
	}
	// import the tar archive(s)
	err := export.PullSpec(args[0], c.creds, c.user)
	core.CheckErr(err, "cannot pull spec artefacts")
}
