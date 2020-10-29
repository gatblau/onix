/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artie/build"
	"github.com/gatblau/onix/artie/core"
	"github.com/spf13/cobra"
)

// create a file seal
type BuildCmd struct {
	cmd          *cobra.Command
	branch       string
	gitTag       string
	artefactName string
	builder      *build.Builder
	gitToken     string
	from         string
	fromPath     string
	profile      string
}

func NewBuildCmd() *BuildCmd {
	c := &BuildCmd{
		cmd: &cobra.Command{
			Use:   "build [flags] [source]",
			Short: "build a package",
			Long:  ``,
		},
		builder: build.NewBuilder(),
	}
	c.cmd.Run = c.Run
	// c.cmd.Flags().StringVarP(&c.branch, "branch", "b", "", "the git branch to use")
	// c.cmd.Flags().StringVarP(&c.gitTag, "gitTag", "l", "", "the git tag to use")
	c.cmd.Flags().StringVarP(&c.gitToken, "token", "k", "", "the git access token")
	c.cmd.Flags().StringVarP(&c.artefactName, "artefactName", "t", "", "name and optionally a tag in the 'name:tag' format")
	c.cmd.Flags().StringVarP(&c.fromPath, "path", "f", "", "the path within the git repository where the root of the source to package is")
	c.cmd.Flags().StringVarP(&c.profile, "profile", "p", "", "the build profile to use. if not provided, the default profile defined in the build file is used. if no default profile is found, then the first profile in the build file is used.")
	return c
}

func (b *BuildCmd) Run(cmd *cobra.Command, args []string) {
	b.from = args[0]
	b.builder.Build(b.from, b.fromPath, b.gitToken, core.ParseName(b.artefactName), b.profile)
}
