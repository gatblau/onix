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
	"github.com/gatblau/onix/artisan/merge"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gatblau/onix/artisan/core"
)

type SyncManager struct {
	repoManager       *RepoManager
	path4Files2BeSync string
	repoPath          string
	fileNamePrefix    string
	workingDir        string
	strictSync        bool
}

// NewSyncManagerFromUri will initialise SyncManager by cloning the repo.
// repoUrl location from where repository has to be cloned
// gitToken token required to clone the repository
// fileNamePrefix string to be appended to the file's name, example when tem file is converted to yaml append this to file name.
// path4Files2BeSync relate path from where the source files must be considered for merge and sync
// repoPath folder with in repo where the final files to be copied during sync operation
// It will return initialised SyncManager or any error occurred
func NewSyncManagerFromUri(repoURI, token, fileNamePrefix, path4Files2BeSync, repoPath string, strictSync bool) (*SyncManager, error) {
	workingDir, err := core.NewTempDir()
	if err != nil {
		return nil, err
	}
	repoManager := NewRepoManager(repoURI, token, workingDir, repoPath, strictSync)
	return &SyncManager{
		repoManager:       repoManager,
		fileNamePrefix:    fileNamePrefix,
		path4Files2BeSync: path4Files2BeSync,
		repoPath:          repoPath,
		workingDir:        workingDir,
		strictSync:        strictSync,
	}, nil
}

// getAbsoluteFilePathToSync will convert given path from where files to be merged can be found,to absolute path
// It will return absolute path or any error if occurred
func (s *SyncManager) getAbsoluteFilePathToSync() (string, error) {
	absPathOfFiles2BeSync, err := core.AbsPath(s.path4Files2BeSync)
	if err != nil {
		return "", err
	}
	return absPathOfFiles2BeSync, nil
}

// getAbsoluteRepoPath will create absolute path by appending repoPath to new temp folder
// It will return absolute path where merged files to copied before sync operation or any error if occurred
func (s *SyncManager) getAbsoluteRepoPath() (string, error) {
	absPathOfRepoFolder := filepath.Join(s.workingDir, s.repoPath)
	return absPathOfRepoFolder, nil
}

// artMerge will find all valid files from source absolute path and merge them with values from environment variable
// It will return any error if occurred
func (s *SyncManager) artMerge() error {
	absSyncPath, err := s.getAbsoluteFilePathToSync()
	if err != nil {
		return err
	}
	// find all file names with extension .tem or .art
	files, err := core.FindFiles(absSyncPath, "^.*\\.(tem|art)$")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no template files (*.tem/*art) found in the path %v\n", absSyncPath)
	}
	// replace environment variable value with the place holder
	envVar := merge.NewEnVarFromSlice(os.Environ())
	merger, err := merge.NewTemplMerger()
	if err != nil {
		return err
	}
	err = merger.LoadTemplates(files)
	if err != nil {
		return err
	}
	err = merger.Merge(envVar)
	if err != nil {
		return err
	}
	return merger.Save()
}

// MergeAndSync will clone the repo at target path and then perform merging of tem files and finally push the changes back to remote git repo
// It will return any error if occurred
func (s *SyncManager) MergeAndSync() error {
	err := s.repoManager.Clone()
	if err != nil {
		return err
	}
	err = s.mergeAndCopy()
	if err != nil {
		return err
	}
	err = s.repoManager.Commit()
	if err != nil {
		return err
	}
	err = s.repoManager.Push()
	if err != nil {
		return err
	}
	return nil
}

// mergeAndCopy will get tem files merged and move yaml files to target repo path from where repo is synchronised
// It will return any error if occurred
func (s *SyncManager) mergeAndCopy() error {
	absRepoPath, err := s.getAbsoluteRepoPath()
	if err != nil {
		return err
	}
	absPath4Files2BeSync, err := s.getAbsoluteFilePathToSync()
	if err != nil {
		return err
	}
	err = s.artMerge()
	if err != nil {
		return err
	}
	// move each yaml file generated after merge to absolute repo path, so that it can be committed and
	// push to remote git
	files, err := ioutil.ReadDir(absPath4Files2BeSync)
	if err != nil {
		return err
	} else {
		err = os.MkdirAll(absRepoPath, os.ModePerm)
		if err != nil {
			return err
		}
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".yaml" {
				// move yaml files from tem files folder to repo path
				oldLocation := filepath.Join(absPath4Files2BeSync, file.Name())
				newFileName := file.Name()
				if len(s.fileNamePrefix) > 0 {
					newFileName = s.fileNamePrefix + "_" + newFileName
				}
				newLocation := filepath.Join(absRepoPath, newFileName)
				err := os.Rename(oldLocation, newLocation)
				if err != nil {
					return err
				}
			}
		}
	}
	log.Println("git, merge and copy completed ")
	return nil
}
