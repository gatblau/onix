/*
  Onix Config Manager - Art
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/pak/core"
	"github.com/spf13/cobra"
)

// create a file seal
type BuildCmd struct {
	cmd         *cobra.Command
	branch      string
	gitTag      string
	packageName string
	builder     *core.Builder
	gitToken    string
	from        string
	fromPath    string
}

func NewBuildCmd() *BuildCmd {
	c := &BuildCmd{
		cmd: &cobra.Command{
			Use:   "build [git_repo_url]",
			Short: "build a package",
			Long:  ``,
		},
		builder: core.NewBuilder(),
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.branch, "branch", "b", "", "the git branch to use")
	c.cmd.Flags().StringVarP(&c.gitTag, "gitTag", "l", "", "the git tag to use")
	c.cmd.Flags().StringVarP(&c.gitToken, "token", "k", "", "the git access token")
	c.cmd.Flags().StringVarP(&c.packageName, "packageName", "t", "", "name and optionally a tag in the 'name:tag' format")
	c.cmd.Flags().StringVarP(&c.gitToken, "path", "p", "", "the path within the git repository where the root of the source to package is")
	return c
}

func (b *BuildCmd) Run(cmd *cobra.Command, args []string) {
	b.from = args[0]
	b.builder.Build(b.from, b.gitToken, b.packageName)
}
