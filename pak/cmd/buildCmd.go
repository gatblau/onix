/*
  Onix Config Manager - Pak
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/pak/core"
	"github.com/spf13/cobra"
	"os"
)

// create a file seal
type BuildCmd struct {
	cmd      *cobra.Command
	branch   string
	tag      string
	packer   *core.Builder
	gitToken string
}

func NewBuildCmd() *BuildCmd {
	c := &BuildCmd{
		cmd: &cobra.Command{
			Use:   "build [git_repo_url]",
			Short: "build a package",
			Long:  ``,
		},
		packer: core.NewBuilder(),
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.branch, "branch", "b", "", "the git branch to use")
	c.cmd.Flags().StringVarP(&c.tag, "tag", "t", "", "the git tag to use")
	c.cmd.Flags().StringVarP(&c.gitToken, "token", "k", "", "the git access token")
	return c
}

func (c *BuildCmd) Run(cmd *cobra.Command, args []string) {
	// get the url of the remote git repository containing the source code
	gitRepoUrl := args[0]
	// execute the build
	c.packer.Build(gitRepoUrl, c.gitToken)
}

// return the working path
func path() string {
	basePath, _ := os.Getwd()
	return fmt.Sprintf("%s/.pak", basePath)
}
