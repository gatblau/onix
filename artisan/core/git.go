package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	l "github.com/sirupsen/logrus"
)

/*
GitClone will clone git repo from input repo url to a temp folder "gitsynccmd"
currents user's home. It only accepts a token if authentication is required
if the token is not provided (empty string) then no authentication is used
*/
func GitClone(repoUrl string, gitToken string, sourceDir string, cmdName string) (*git.Repository, error) {
	// clone the remote repository
	opts := &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	}
	// if authentication token has been provided
	if len(gitToken) > 0 {
		l.Debug("git, GitClone token is provided ")
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		opts.Auth = &http.BasicAuth{
			Username: cmdName, // yes, this can be anything except an empty string
			Password: gitToken,
		}
	}
	l.Info("git, GitClone peforming  git.PlainClone ")
	return git.PlainClone(sourceDir, false, opts)
}

/*
Sync will replace environment variables in the tem files with the respective value
and add asset folder into git working tree
*/
func Sync(temFilesWithPath []string, absPath4SrcFiles string, absPath4Repo string, filePrefix string) error {
	var err error
	l.Infof("git, Sync tem file size is  ", len(temFilesWithPath))
	// replace environment variable value with the place holder
	envVar := NewEnVarFromSlice(os.Environ())
	MergeFiles(temFilesWithPath, envVar)
	l.Info("git, Sync tem file merge completed")
	// move each yaml file generated after merge to absolute repo path, so that it can be commited and
	// push to remote git
	l.Infof("git, Sync, moving yaml files generated after merge from path %s to local git repo path %s ", absPath4SrcFiles, absPath4Repo)
	files, err := ioutil.ReadDir(absPath4SrcFiles)
	if err != nil {
		l.Fatalf("git, Sync, error while reading files from path %s , Error :- %s ", absPath4SrcFiles, err)
		return err
	} else {
		os.MkdirAll(absPath4Repo, os.ModePerm)
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".yaml" {
				// move yaml files from tem files folder to repo path
				oldLocation := filepath.Join(absPath4SrcFiles, file.Name())
				newFileName := file.Name()
				if len(filePrefix) > 0 {
					newFileName = filePrefix + "-" + newFileName
				}
				newLocation := filepath.Join(absPath4Repo, newFileName)
				err := os.Rename(oldLocation, newLocation)
				if err != nil {
					l.Fatalf("git, Sync , error while moving file [ %s ] from path [ %s ] to path [ %s ] \n Error :- %s ", file.Name(), absPath4SrcFiles, absPath4Repo, err)
					return err
				}
			}
		}
	}
	l.Info("git, Sync completed ")
	return err
}

// CommitAndPush will commit and push the changes back to remote git repo
func CommitAndPush(repo *git.Repository, token string, cmdName string) error {
	// commit changes
	var err error
	wrkTree, err := repo.Worktree()
	if err != nil {
		l.Fatalf("git, CommitAndPush, error while getting Working tree for git rep", err)
		return err
	}
	wrkTree.Add(".")
	l.Info("git, CommitAndPush added current folder into working tree")
	commit, err := wrkTree.Commit("Changes committed by  GitSyncCmd for the ?????", &git.CommitOptions{
		Author: &object.Signature{
			Name:  cmdName,
			Email: "******",
			When:  time.Now(),
		},
	})
	if err != nil {
		l.Fatalf("git, CommitAndPush, error while creating commit instance from Working tree ", err)
		return err
	}

	l.Info("git, CommitAndPush creating commit ")
	_, err = repo.CommitObject(commit)
	if err != nil {
		log.Fatalf("CommitAndPush, error while committing to git repo ", err)
		log.Fatal(err)
		return err
	}

	auth := &http.BasicAuth{
		Username: cmdName, // yes, this can be anything except an empty string
		Password: token,
	}

	l.Info("git, CommitAndPush pushing the changes ")
	err = repo.Push(&git.PushOptions{Auth: auth})
	if err != nil {
		l.Fatal("CommitAndPush, error while pushing changes to git repo ", err)
		return err
	}
	l.Info("git, CommitAndPush pushed the changes ")
	return err
}
