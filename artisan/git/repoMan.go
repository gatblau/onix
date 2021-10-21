package git

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type RepoManager struct {
	repo       *git.Repository
	repoURI    string
	token      string
	workingDir string
	repoPath   string
	strictSync bool
	branch     string
}

// NewRepoManager will initialise RepoManager.
// repoURI location from where repository has to be cloned
// token token required to clone the repository
// workingDir directory where repo will be cloned
// It will return initialised RepoManager or any  error occurred
func NewRepoManager(repoURI, token, workingDir, repoPath string, strictSync bool, branch string) *RepoManager {
	return &RepoManager{
		repoURI: repoURI, token: token, workingDir: workingDir, repoPath: repoPath, strictSync: strictSync, branch: branch,
	}
}

// Clone will clone git repo from repo uri to a temp folder. It only accepts a token if authentication is required
// It will return any error if occurred
func (g *RepoManager) Clone() error {
	if len(g.repoURI) == 0 {
		return fmt.Errorf("git repo URI is missing in RepoManager ")
	}
	var opts *git.CloneOptions
	if len(g.branch) > 0 {
		opts = &git.CloneOptions{
			URL:           g.repoURI,
			Progress:      os.Stdout,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.branch)),
		}
	} else {
		opts = &git.CloneOptions{
			URL:      g.repoURI,
			Progress: os.Stdout,
		}
	}

	// if authentication token has been provided
	if len(g.token) > 0 {
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		opts.Auth = &http.BasicAuth{
			Username: "****", // yes, this can be anything except an empty string
			Password: g.token,
		}
	}
	grepo, err := git.PlainClone(g.workingDir, false, opts)
	if err != nil {
		return err
	}
	g.repo = grepo

	wrkTree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	if g.strictSync {
		wrkTree.Remove(g.repoPath)
	}

	return err
}

// Commit will commit the changes.
// It will return any error if occured
func (g *RepoManager) Commit() error {
	var err error
	wrkTree, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	wrkTree.Add(".")
	commit, err := wrkTree.Commit("automated commit by artisan", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "******",
			Email: "******",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	_, err = g.repo.CommitObject(commit)
	if err != nil {
		return err
	}
	fmt.Println(" commit performed  .... ")
	return err
}

// Push will push the local committed changes back to remote git repo.
// It will return any error if occurred
func (g *RepoManager) Push() error {
	var err error
	auth := &http.BasicAuth{
		Username: "*****", // yes, this can be anything except an empty string
		Password: g.token,
	}
	err = g.repo.Push(&git.PushOptions{Auth: auth})
	if err != nil {
		return err
	}
	fmt.Println(" push performed  .... ")
	return err
}
