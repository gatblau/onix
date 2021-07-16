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
	"log"
	"os"
	"time"

	"github.com/gatblau/onix/artisan/core"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type RepoManager struct {
	repo       *git.Repository
	repoURI    string
	token      string
	workingDir string
}

// NewRepoManager will initialise RepoManager.
// repoURI location from where repository has to be cloned
// token token required to clone the repository
// workingDir directory where repo will be cloned
// It will return initialised RepoManager or any  error occured
func NewRepoManager(repoURI, token, workingDir string) *RepoManager {
	return &RepoManager{
		repoURI: repoURI, token: token, workingDir: workingDir,
	}
}

// Clone will clone git repo from repo uri to a temp folder. It only accepts a token if authentication is required
// It will return any error if occured
func (g *RepoManager) Clone() error {
	if len(g.repoURI) == 0 {
		return fmt.Errorf("git repo URI is missing in RepoManager ")
	}

	// clone the remote repository
	opts := &git.CloneOptions{
		URL:      g.repoURI,
		Progress: os.Stdout,
	}
	// if authentication token has been provided
	if len(g.token) > 0 {
		core.Debug("git, GitClone token is provided ")
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		opts.Auth = &http.BasicAuth{
			Username: "****", // yes, this can be anything except an empty string
			Password: g.token,
		}
	}
	log.Println("git, GitClone peforming  git.PlainClone ")
	grepo, err := git.PlainClone(g.workingDir, false, opts)
	g.repo = grepo

	return err
}

// Commit will commit the changes.
// It will return any error if occured
func (g *RepoManager) Commit() error {
	var err error
	wrkTree, err := g.repo.Worktree()
	if err != nil {
		log.Printf("git, Commit, error while getting Working tree for git rep %v", err)
		return err
	}
	wrkTree.Add(".")
	log.Println("git, Commit, added current folder into working tree")
	commit, err := wrkTree.Commit("Auto generated- Changes committed by  RepoManager ", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "******",
			Email: "******",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Printf("git, Commit, error while creating commit instance from Working tree %v ", err)
		return err
	}

	log.Println("git, Commit , creating commit ")
	_, err = g.repo.CommitObject(commit)
	if err != nil {
		log.Printf("Commit ,, error while committing to git repo %v ", err)
		return err
	}

	log.Println("git, Commit executed successfully ")
	return err
}

// Push will push the local commited changes back to remote git repo.
// It will return any error if occured
func (g *RepoManager) Push() error {
	var err error
	auth := &http.BasicAuth{
		Username: "*****", // yes, this can be anything except an empty string
		Password: g.token,
	}

	log.Println("git, Push pushing the changes ")
	err = g.repo.Push(&git.PushOptions{Auth: auth})
	if err != nil {
		log.Println("git,Push, error while pushing changes to git repo %v ", err)
		return err
	}
	log.Println("git, pushed local changes to remote successfully")
	return err
}
