package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gatblau/onix/artisan/core"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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
	recursive         bool
	path4Files2BeSync string
}

// NewGitSyncCmd create a new GitSyncCmd.
func NewGitSyncCmd() *GitSyncCmd {
	c := &GitSyncCmd{
		cmd: &cobra.Command{
			Use:   "sync",
			Short: "sync local git repo with remote",
			Long:  ``,
		},
	}

	c.cmd.Flags().StringVarP(&c.repoPath, "repo-path", "p", "", "the git repo to sync")
	c.cmd.Flags().StringVarP(&c.repoURI, "uri", "u", "", "the git repo uri to use")
	c.cmd.Flags().StringVarP(&c.token, "token", "t", "", "the token to login to git repo")
	c.cmd.Flags().BoolVarP(&c.recursive, "recursive", "r", "", "whether to perform recursive sync. true or false (currently not implemented) ")
	c.cmd.Run = c.Run
	return c
}

//Run to execute git sync
func (g *GitSyncCmd) Run(cmd *cobra.Command, args []string) {
	//art sync --repo-uri git-url --repo-path  -u user:password --recursive  .
	println("Ani is great ..!")
	switch len(args) {
	case 0:
		g.path4Files2BeSync = "."
	case 1:
		g.path4Files2BeSync = args[0]
	default:
		core.RaiseErr("too many arguments")
	}
	//Query:- we need two set of input
	// 1> input to connect git repo like repoURL, assetFolder, git token
	// 2> environment variable with values for doing art merge
	// 3> Should we validate inputs parameters provided to below fuctions

	// init(gitURL string, gitToken string)
	// cloneRepo(repoUrl string, gitToken string)
	// syncFiles(repo *git.Repository)
	// commitAndPush(wrkTree *git.Worktree, repo *git.Repository)
	core.newTempDir()
}

/*
init create a temp directory under current user's home directory.
This temp directory will be used to perform git operations like clone,pull,push
and finally this will be deleted if required.
*/
func (g *GitSyncCmd) init(gitURL string, gitToken string) {

	filepath.Join(core.TempPath(), core.RandomString(10))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	} else {
		// create a temp directory with name "gitsynccmd" under current user's home directory
		wrkDir, err := ioutil.TempDir(homeDir, "gitsynccmd")
		if err != nil {
			log.Fatal(err)
		} else {
			// store this reference back in gitSyncCmd structs
			g.workingDir = wrkDir
		}
	}
}

/*
cloneRepo will clone git repo from input repo url to a temp folder "gitsynccmd"
currents user's home
*/
func (g *GitSyncCmd) cloneRepo(repoUrl string, gitToken string) *git.Repository {
	g.repoURI = repoUrl
	opts := &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	}
	opts.Auth = &http.BasicAuth{
		Username: "SyncCmd", // yes, this can be anything except an empty string
		Password: gitToken,
	}
	repo, err := git.PlainClone(g.workingDir, false, opts)
	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		log.Fatal(err)
	}
	return repo
}

/*
syncFiles will replace environment variables in the files with the respective value
and add asset folder into git working tree
*/
func (g *GitSyncCmd) syncFiles(repo *git.Repository) *git.Worktree {
	wrkTree, err := repo.Worktree()
	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		log.Fatal(err)
	} else {
		// merge all .yaml.tem files
		//move .yaml files generated to g.workingDir /g.repoName /g.assetFolder

		//add asset folder to git working tree
		_, err = wrkTree.Add(filepath.Join(g.workingDir, g.repoName, g.assetFolder))
		if err != nil {
			_ = os.RemoveAll(g.workingDir)
			log.Fatal(err)
		}
	}

	return wrkTree
}

/*
 commitAndPush will commint and push the changes back to remote git repo
*/
func (g *GitSyncCmd) commitAndPush(wrkTree *git.Worktree, repo *git.Repository) {
	//commit changes
	commit, err := wrkTree.Commit("Changes committed for the ?????", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "SyncCmd",
			Email: "******",
			When:  time.Now(),
		},
	})

	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		log.Fatal(err)
	}
	obj, err := repo.CommitObject(commit)
	if err != nil {
		_ = os.RemoveAll(g.workingDir)
		log.Fatal(err)
	}
	err = repo.Push(&git.PushOptions{})
	fmt.Println(obj)
}
