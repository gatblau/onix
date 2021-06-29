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
	cmd                  *cobra.Command
	workingDir           string
	repoPath             string
	repoURI              string
	env                  *core.Envar
	token                string
	recursive            *bool
	path4Files2BeSync    string
	absRepoPath          string
	absPath4Files2BeSync string
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
	var cmdName = "GitSyncCmd"
	tmpDir, err := core.NewTempDir()
	if err != nil {
		l.Fatal("Run, error occurred during core.NewTempDir operation ", err)
	} else {
		g.workingDir = tmpDir
		l.Debug("GitSyncCmd is [path4Files2BeSync = " + g.path4Files2BeSync + "] [repoPath = " + g.repoPath + "] [repoURI = " + g.repoURI + "] [workingDir = " + g.workingDir + "]")

		// perform git cloning
		gitRepo, err := core.GitClone(g.repoURI, g.token, g.workingDir, cmdName)
		if err != nil {
			_ = os.RemoveAll(g.workingDir)
			l.Fatal("Run, error occurred during core.GitClone operation ", err)
		} else {

			l.Info("Run, core.GitClone completed ...")
			//set the absolute paths
			g.setAbsolutePaths()
			//get all the tem files from the folder provided by calling application
			temFilesWithPath := g.getAllTemFilesFromPath()
			if len(temFilesWithPath) > 0 && len(g.absPath4Files2BeSync) > 0 && len(g.absRepoPath) > 0 {
				// perform Sync operation by doing merging of environment variable values and then
				// move the yaml files generated after merging, to repo path
				err = core.Sync(temFilesWithPath, g.absPath4Files2BeSync, g.absRepoPath)
				if err != nil {
					_ = os.RemoveAll(g.workingDir)
					l.Fatal("Run, error occurred while executing core.Sync function using files %s", temFilesWithPath)
					l.Fatal("Run, error occurred while executing core.Sync function from folder "+g.absPath4Files2BeSync+" to folder "+g.absRepoPath, err)
				} else {
					l.Info("Run, core.Sync function executed successfully ")
					// commit and push the yaml files copied at repo path
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

/*
  setAbsolutePaths, will find and set the absolute path for folder which has tem files and
  the repo path where final yaml files generated after merge has to be copied
*/
func (g *GitSyncCmd) setAbsolutePaths() {
	absPath4TemFiles, err := core.AbsPath(g.path4Files2BeSync)
	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		l.Fatal("Run, error occurred when getting absolute path of relative path "+g.path4Files2BeSync, err)
	} else {
		g.absPath4Files2BeSync = absPath4TemFiles
		var repoAbsPath = g.workingDir
		if g.repoPath != "." {
			repoAbsPath = filepath.Join(repoAbsPath, g.repoPath)
		}
		g.absRepoPath = repoAbsPath
	}
}

/*
getAllTemFilesFromPath, this fuction will find and return slice containing all the tem file with
absolute path from the folder specified by the calling application
*/
func (g *GitSyncCmd) getAllTemFilesFromPath() []string {
	var temFilesWithPath []string
	l.Debug("gitSyncCmd, getAllTemFilesFromPath, using absolute path to get tem file %s", g.absPath4Files2BeSync)
	if len(g.absPath4Files2BeSync) > 0 {
		files, err := ioutil.ReadDir(g.absPath4Files2BeSync)
		if err != nil {
			l.Fatal("getAllTemFilesFromPath, error occurred when readDir with abs path "+g.absPath4Files2BeSync, err)
			_ = os.RemoveAll(g.workingDir)
			return temFilesWithPath
		}
		l.Info("gitSyncCmd, getAllTemFilesFromPath, Read files from path ")
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".tem" {
				temFilesWithPath = append(temFilesWithPath, filepath.Join(g.absPath4Files2BeSync, file.Name()))
			}
		}
		l.Info("gitSyncCmd, getAllTemFilesFromPath, total number of tem files found is %s", len(temFilesWithPath))
	}

	return temFilesWithPath
}
