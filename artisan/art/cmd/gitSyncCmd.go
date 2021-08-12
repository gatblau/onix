package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/git"
	"github.com/spf13/cobra"
)

// GitSyncCmd sync git resources.
type GitSyncCmd struct {
	cmd               *cobra.Command
	repoPath          string
	repoURI           string
	token             string
	recursive         *bool
	path4Files2BeSync string
	yamlFilePrefix    string
	strictSync        bool
}

// NewGitSyncCmd create a new GitSyncCmd.
func NewGitSyncCmd() *GitSyncCmd {
	c := &GitSyncCmd{
		cmd: &cobra.Command{
			Use: "sync [flags] [path/to/template/files]\n" +
				"  the path to the *.tem or *.art files is optional, if no path is specified, the current path [.] is used",
			Short: "synchronise a remote git repository with the content of a local folder containing files and/or merged Artisan templates",
			Long: "\nprogrammatically update the content of a git repository by \n\n" +
				"a) merging a set of Artisan templates with environment variables\n" +
				"b) updating a local git repository using the merged files; and \n" +
				"c) committing and pushing the local git changes back to its remote origin",
			Example: `art sync -p path/within/git-repo -t the-git-authentication-token -u https://git-host/git-repo.git [. or path/to/.tem or .art files]`,
		},
	}

	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.repoPath, "repo-path", "p", "", "the path within the git repository to be synchronised")
	c.cmd.Flags().StringVarP(&c.repoURI, "uri", "u", "", "the URI of the git repository to synchronise")
	c.cmd.Flags().StringVarP(&c.token, "token", "t", "", "the token to authenticate with the git repository")
	c.cmd.Flags().StringVarP(&c.yamlFilePrefix, "yaml-file-prefix", "x", "", "The prefix to be added to yaml file name after merge")
	c.cmd.Flags().BoolVarP(&c.strictSync, "strict", "s", false, "whether strict synchronisation to be followed, by delete existing repo path and create new folder with same name")
	//c.recursive = c.cmd.Flags().BoolP("recursive", "r", false, "whether to perform recursive sync. true or false (currently not implemented) ")
	// this is causing problem, the recursive value is coming as args in Run function.
	return c
}

// Run to execute git sync
func (g *GitSyncCmd) Run(cmd *cobra.Command, args []string) {
	// art sync --repo-uri git-url --repo-path  -u user:password --recursive  .
	switch len(args) {
	case 0:
		g.path4Files2BeSync = "."
	case 1:
		g.path4Files2BeSync = args[0]
	default:
		core.RaiseErr("too many arguments")
	}
	sm, err := git.NewSyncManagerFromUri(g.repoURI, g.token, g.yamlFilePrefix, g.path4Files2BeSync, g.repoPath, g.strictSync)
	core.CheckErr(err, "Failed to initialise SyncManagerFromUri ")
	err = sm.MergeAndSync()
	core.CheckErr(err, "Failed to perform sync operation")
}
