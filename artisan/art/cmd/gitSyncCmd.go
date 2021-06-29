package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gatblau/onix/artisan/core"
	l "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GitSyncCmd sync git resources.
type GitSyncCmd struct {
	cmd               *cobra.Command
	workingDir        string
	repoPath          string
	repoURI           string
	env               *core.Envar
	token             string
	recursive         *bool
	path4Files2BeSync string
}

// NewGitSyncCmd create a new GitSyncCmd.
func NewGitSyncCmd() *GitSyncCmd {
	c := &GitSyncCmd{
		cmd: &cobra.Command{
			Use:   "sync",
			Short: "sync local git repo with remote",
			Long:  `art sync -p repo/vm-manifest -t mytoken -u https://k-p-ani@bitbucket.org/k-p-ani/vm-auto-deployment.git   . or path/to/.tem files`,
		},
	}

	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.repoPath, "repo-path", "p", "", "the git repo to sync")
	c.cmd.Flags().StringVarP(&c.repoURI, "uri", "u", "", "the git repo uri to use")
	c.cmd.Flags().StringVarP(&c.token, "token", "t", "", "the token to login to git repo")
	//c.recursive = c.cmd.Flags().BoolP("recursive", "r", false, "whether to perform recursive sync. true or false (currently not implemented) ")
	// this is causing problem, the recursive value is coming as args in Run function.
	return c
}

//Run to execute git sync
func (g *GitSyncCmd) Run(cmd *cobra.Command, args []string) {
	//art sync --repo-uri git-url --repo-path  -u user:password --recursive  .
	switch len(args) {
	case 0:
		g.path4Files2BeSync = "."
	case 1:
		g.path4Files2BeSync = args[0]
	default:
		core.RaiseErr("too many arguments")
	}

	g.workingDir = core.NewTempDir()
	l.Debug("GitSyncCmd is [path4Files2BeSync = " + g.path4Files2BeSync + "] [repoPath = " + g.repoPath + "] [repoURI = " + g.repoURI + "] [workingDir = " + g.workingDir + "]")
	var cmdName = "GitSyncCmd"
	gitRepo, err := core.GitClone(g.repoURI, g.token, g.workingDir, cmdName)
	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		l.Fatal("Run, error occurred during core.GitClone operation ", err)
	} else {
		l.Info("Run, core.GitClone completed ...")
		absPath4TemFiles, err := core.AbsPath(g.path4Files2BeSync)
		if err != nil {
			_ = os.RemoveAll(g.workingDir)
			l.Fatal("Run, error occurred when getting absolute path of relative path "+g.path4Files2BeSync, err)
		} else {
			l.Info("Run, absolute path generated for path ", g.path4Files2BeSync)
			temFilesWithPath, err := getAllTemFilesFromPath(absPath4TemFiles)
			if err != nil {
				_ = os.RemoveAll(g.workingDir)
				l.Fatal("Run, error occurred when getting all tem files from path "+absPath4TemFiles, err)
			} else {
				l.Info("successfully call getAllTemFilesFromPath function for fetching tem file from path ", absPath4TemFiles)
				l.Debug("Run, getAllTemFilesFromPath function returned total tem files ", len(temFilesWithPath))
				var repoAbsPath = g.workingDir
				if g.repoPath != "." {
					repoAbsPath = filepath.Join(repoAbsPath, g.repoPath)
				}
				err = core.Sync(temFilesWithPath, absPath4TemFiles, repoAbsPath)
				if err != nil {
					_ = os.RemoveAll(g.workingDir)
					l.Fatal("Run, error occurred while executing core.Sync function using files %s", temFilesWithPath)
					l.Fatal("Run, error occurred while executing core.Sync function from folder "+absPath4TemFiles+" to folder "+repoAbsPath, err)
				} else {
					l.Info("Run, core.Sync function executed successfully ")
					err = core.CommitAndPush(gitRepo, g.token, cmdName)
					if err != nil {
						_ = os.RemoveAll(g.workingDir)
						l.Fatal("Run, error occurred while executing core.CommitAndPush function from folder ", err)
					}
					l.Info("Run, core.CommitAndPush function executed successfully ")
				}
			}
		}
	}
}

func getAllTemFilesFromPath(absPath4TemFiles string) ([]string, error) {
	var err error
	var temFilesWithPath []string
	files, err := ioutil.ReadDir(absPath4TemFiles)
	if err != nil {
		l.Fatal("getAllTemFilesFromPath, error occurred when readDir with abs path "+absPath4TemFiles, err)
		return temFilesWithPath, err
	}
	l.Info("gitSyncCmd, getAllTemFilesFromPath, Read files from path ")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tem" {
			temFilesWithPath = append(temFilesWithPath, filepath.Join(absPath4TemFiles, file.Name()))
		}
	}
	l.Info("gitSyncCmd, getAllTemFilesFromPath, total number of tem files found is %s", len(temFilesWithPath))
	return temFilesWithPath, err
}
