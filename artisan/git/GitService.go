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
	fileType          string
	workingDir        string
<<<<<<< HEAD
	strictSync        bool
=======
>>>>>>> 556ebf1fa3a975112b38b30ef8771bd794ddf143
}

// NewSyncManagerFromUri will initialise SyncManager by cloning the repo.
// repoUrl location from where repository has to be cloned
// gitToken token required to clone the repository
// fileType type of source files to be considered for merge and sync
// fileNamePrefix string to be append to the file's name, example when tem file is converted to yaml append this to file name.
// path4Files2BeSync relate path from where the source files must be consider for merge and sync
// repoPath folder with in repo where the final files to be copied during sync operation
// It will return initialised SyncManager or any error occured
func NewSyncManagerFromUri(repoURI, token, fileType, fileNamePrefix, path4Files2BeSync, repoPath string,
	strictSync bool) (*SyncManager, error) {
	workingDir, err := core.NewTempDir()
	if err != nil {
		log.Printf("GitService,NewSyncManagerFromUri, error while invoking core.NewTempDir operation ", err)
		return nil, err
	}

	repoManager := NewRepoManager(repoURI, token, workingDir, repoPath, strictSync)
	core.Debug(" GitService,PerformSync, RepoManager structs created ")
	return &SyncManager{
		repoManager: repoManager, fileType: fileType, fileNamePrefix: fileNamePrefix,
		path4Files2BeSync: path4Files2BeSync, repoPath: repoPath, workingDir: workingDir,
		strictSync: strictSync,
	}, nil
}

// getAbsoluteFilePathToSync will convert given path from where files to be merged can be found,to absolute path
// It will return absolute path or any error if occured
func (s *SyncManager) getAbsoluteFilePathToSync() (string, error) {
	absPathOfFiles2BeSync, err := core.AbsPath(s.path4Files2BeSync)
	if err != nil {
		log.Printf("GitService,getAbsoluteFilePathToSync, error while getting absolute path for source file path '%v' Error:- %v ", s.path4Files2BeSync, err)
		return "", err
	}
	core.Debug("source folder path detail %v ", absPathOfFiles2BeSync)
	return absPathOfFiles2BeSync, nil
}

// getAbsoluteRepoPath will create absolute path by appending repoPath to new temp folder
// It will return absolute path where merged files to copied before sync operation or any error if occured
func (s *SyncManager) getAbsoluteRepoPath() (string, error) {
	core.Debug(" GitService,getAbsoluteRepoPath, creating abs path absPathOfRepoFolder using working dir %v , repo path %v ", s.workingDir, s.repoPath)
	absPathOfRepoFolder := filepath.Join(s.workingDir, s.repoPath)
	core.Debug("target folder path detail %v ", absPathOfRepoFolder)

	return absPathOfRepoFolder, nil
}

// artMerge will find all valid files from source absolute path and merge them with values from environment variable
// It will return any error if occured
func (s *SyncManager) artMerge() error {

	absSyncPath, err := s.getAbsoluteFilePathToSync()
	if err != nil {
		return err
	}
	files, err := core.GetFiles(absSyncPath, s.fileType)
	if err != nil {
		log.Printf("GitService,artMergeFiles, error while getting '%v' files from path '%v' Error:- '%v'", s.fileType, absSyncPath, err)
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("GitService,artMergeFiles, No file found for file type %v in the path %v ", s.fileType, absSyncPath)
	}

	log.Printf("GitService, artMergeFiles file size is %v ", fmt.Sprintf("%v", len(files)))
	// replace environment variable value with the place holder
	envVar := core.NewEnVarFromSlice(os.Environ())
	core.MergeFiles(files, envVar)
	return nil
}

// MergeAndSync will clone the repo at target path and then perform merging of tem files and finally push the changes back to remote git repo
// It will return any error if occured
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
	core.Debug(" GitService,PerformSync, core.CommitAndPush function executed successfully ")
	return nil
}

// mergeAndCopy will get tem files merged and move yaml files to target repo path from where repo is synchronised
// It will return any error if occured
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
	// move each yaml file generated after merge to absolute repo path, so that it can be commited and
	// push to remote git
	log.Printf("git, mergeAndCopy, moving yaml files generated after merge from path %v to local git repo path %v ", absPath4Files2BeSync, absRepoPath)
	files, err := ioutil.ReadDir(absPath4Files2BeSync)
	if err != nil {
		log.Printf("git, mergeAndCopy, error while reading files from path %v , Error :- %v ", absPath4Files2BeSync, err)
		return err
	} else {
		err = os.MkdirAll(absRepoPath, os.ModePerm)
		if err != nil {
			log.Printf("git, mergeAndCopy, error while creating folder %v , Error :- %v ", absRepoPath, err)
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
					log.Printf("git, mergeAndCopy , error while moving file [ %v ] from path [ %v ] to path [ %v ] \n Error :- %s ", file.Name(), absPath4Files2BeSync, absRepoPath, err)
					return err
				}
			}
		}
	}
	log.Println("git, merge and copy completed ")
	return nil
}
