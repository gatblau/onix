/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

// create a file seal
type BuildCmd struct {
	cmd          *cobra.Command
	branch       string
	gitTag       string
	artefactName string
	gitToken     string
	from         string
	fromPath     string
	profile      string
	copySource   *bool
	interactive  *bool
}

func NewBuildCmd() *BuildCmd {
	c := &BuildCmd{
		cmd: &cobra.Command{
			Use:   "build [flags] [source]",
			Short: "build a package",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	// c.cmd.Flags().StringVarP(&c.branch, "branch", "b", "", "the git branch to use")
	// c.cmd.Flags().StringVarP(&c.gitTag, "gitTag", "l", "", "the git tag to use")
	c.cmd.Flags().StringVarP(&c.gitToken, "token", "k", "", "the git access token")
	c.cmd.Flags().StringVarP(&c.artefactName, "artefactName", "t", "", "name and optionally a tag in the 'name:tag' format")
	c.cmd.Flags().StringVarP(&c.fromPath, "path", "f", "", "the path within the git repository where the root of the source to package is")
	c.cmd.Flags().StringVarP(&c.profile, "profile", "p", "", "the build profile to use. if not provided, the default profile defined in the build file is used. if no default profile is found, then the first profile in the build file is used.")
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.copySource = c.cmd.Flags().BoolP("copy", "c", false, "indicates if a copy should be made of the project files before building the artefact. it is only applicable if the source is in the file system. ")
	return c
}

func (b *BuildCmd) Run(cmd *cobra.Command, args []string) {
	// validate build path
	switch len(args) {
	case 0:
		b.from = "."
	case 1:
		b.from = args[0]
	default:
		core.RaiseErr("too many arguments")
	}
	builder := build.NewBuilder()
	builder.Build(b.from, b.fromPath, b.gitToken, core.ParseName(b.artefactName), b.profile, *b.copySource, *b.interactive)
}
