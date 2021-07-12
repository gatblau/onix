package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"path/filepath"

	"github.com/gatblau/onix/artisan/core"
	l "github.com/sirupsen/logrus"
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
}

// NewGitSyncCmd create a new GitSyncCmd.
func NewGitSyncCmd() *GitSyncCmd {
	c := &GitSyncCmd{
		cmd: &cobra.Command{
			Use: "sync [flags] [path/to/template/files]\n" +
				"  the path to the tem files is optional, if no path is specified, the current path [.] is used",
			Short: "synchronise a remote git repository with the content of a local folder containing files and/or merged Artisan templates",
			Long: "\nprogrammatically update the content of a git repository by \n\n" +
				"a) merging a set of Artisan templates with environment variables\n" +
				"b) updating a local git repository using the merged files; and \n" +
				"c) committing and pushing the local git changes back to its remote origin",
			Example: `art sync -p path/within/git-repo -t the-git-authentication-token -u https://git-host/git-repo.git [. or path/to/.tem files]`,
		},
	}

	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.repoPath, "repo-path", "p", "", "the path within the git repository to be synchronised")
	c.cmd.Flags().StringVarP(&c.repoURI, "uri", "u", "", "the URI of the git repository to synchronise")
	c.cmd.Flags().StringVarP(&c.token, "token", "t", "", "the token to authenticate with the git repository")
	c.cmd.Flags().StringVarP(&c.yamlFilePrefix, "yaml-file-prefix", "x", "", "The prefix to be added to yaml file name after merge")
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

	l.Debug(" Executing gitSyncCmd ")

	var cmdName = "GitSyncCmd"
	var fileType = ".tem"
	workingDir, err := core.NewTempDir()
	core.CheckErr(err, "Run, error occurred during core.NewTempDir operation ")
	l.Debug("Run, temp file created successfully ")

	gitRepo, err := core.GitClone(g.repoURI, g.token, workingDir, cmdName)
	core.CheckErr(err, "Run, error occurred during core.GitClone operation")
	l.Debug("Run, core.GitClone operation completed successfully ")

	absPathOfTemFiles, err := core.AbsPath(g.path4Files2BeSync)
	core.CheckErr(err, "Run, core.AbsPath, error occurred while getting absolute path for path '%s'", g.path4Files2BeSync)
	l.Debugf("source folder path detail %s ", absPathOfTemFiles)

	absPathOfRepoFolder := filepath.Join(workingDir, g.repoPath)
	l.Debugf("target folder path detail %s ", absPathOfRepoFolder)

	filesWithPath, err := core.GetFiles(absPathOfTemFiles, fileType)
	core.CheckErr(err, "Run, core.GetFiles error while getting '%s' files from path '%s' Error:- '%s'", fileType, absPathOfTemFiles)
	l.Debug("Run, core.GetFiles operation executed successfully ")

	err = core.Sync(filesWithPath, absPathOfTemFiles, absPathOfRepoFolder, g.yamlFilePrefix)
	core.CheckErr(err, "Run, error occurred while executing core.Sync function  from folder '%s' to folder '%s' , Error:- '%s'", absPathOfTemFiles, absPathOfRepoFolder)
	l.Debug("Run, core.Sync operation executed successfully using,Source path,Target path, Target file name prefix %s ", g.yamlFilePrefix)

	err = core.CommitAndPush(gitRepo, g.token, cmdName)
	core.CheckErr(err, "Run, error occurred while executing core.CommitAndPush function ")
	l.Debug("Run, core.CommitAndPush function executed successfully ")
}
