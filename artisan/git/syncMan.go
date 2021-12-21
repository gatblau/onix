package git

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/gatblau/onix/artisan/merge"

	"github.com/gatblau/onix/artisan/core"
)

const FILE_EXT_REGEX = "^.*\\.(tem|art)$"

type SyncManager struct {
	repoManager       *RepoManager
	path4Files2BeSync string
	repoPath          string
	fileNamePrefix    string
	workingDir        string
	strictSync        bool
	regEx             *regexp.Regexp
	templateFiles     []string
	recursive         bool
	preserveFiles     bool
}

// NewSyncManagerFromUri will initialise SyncManager by cloning the repo.
// repoUrl location from where repository has to be cloned
// gitToken token required to clone the repository
// fileNamePrefix string to be appended to the file's name, example when tem file is converted to yaml append this to file name.
// path4Files2BeSync relate path from where the source files must be considered for merge and sync
// repoPath folder with in repo where the final files to be copied during sync operation
// It will return initialised SyncManager or any error occurred
func NewSyncManagerFromUri(repoURI, token, fileNamePrefix, path4Files2BeSync, repoPath string,
	strictSync bool, recursive bool, tempPath string, preserveFiles bool, branch string) (*SyncManager, error) {

	var err error
	var workingDir string

	if len(tempPath) == 0 {
		workingDir, err = core.NewTempDir()
	} else {
		err = os.MkdirAll(tempPath, os.ModePerm)
		workingDir = tempPath
	}
	if err != nil {
		return nil, err
	}

	e, err := regexp.Compile(FILE_EXT_REGEX)
	if err != nil {
		return nil, err
	}

	repoManager := NewRepoManager(repoURI, token, workingDir, repoPath, strictSync, branch)
	return &SyncManager{
		repoManager:       repoManager,
		fileNamePrefix:    fileNamePrefix,
		path4Files2BeSync: path4Files2BeSync,
		repoPath:          repoPath,
		workingDir:        workingDir,
		strictSync:        strictSync,
		regEx:             e,
		recursive:         recursive,
		preserveFiles:     preserveFiles,
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
	absRepoPath, err := s.getAbsoluteRepoPath()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	s.copyFolder(absSyncPath, absRepoPath)
	// find all file names with extension .tem or .art
	if err != nil {
		return err
	}
	if len(s.templateFiles) == 0 {
		return fmt.Errorf("no template files (*.tem/*art) found in the path %v\n", absSyncPath)
	}
	// replace environment variable value with the place holder
	envVar := merge.NewEnVarFromSlice(os.Environ())
	merger, err := merge.NewTemplMerger()
	if err != nil {
		return err
	}
	err = merger.LoadTemplates(s.templateFiles)
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

	err = s.artMerge()
	if err != nil {
		return err
	}

	/*delete template files so that they don't get checked in*/
	s.deleteTemplateFiles()
	err = s.repoManager.Commit()
	if err != nil {
		return err
	}
	err = s.repoManager.Push()
	if err != nil {
		return err
	}

	if !s.preserveFiles {
		s.deleteTempFolder()
	}
	return nil
}

// copy the files in a folder recursively
func (s *SyncManager) copyFolder(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if s.recursive {
				if err = s.copyFolder(srcFp, dstFp); err != nil {
					core.ErrorLogger.Printf(err.Error())
				}
			}
		} else {
			if s.regEx.Match([]byte(fd.Name())) {
				s.templateFiles = append(s.templateFiles, dstFp)
			}
			if err = s.copyFile(srcFp, dstFp); err != nil {
				core.ErrorLogger.Printf(err.Error())
			}
		}
	}
	return nil
}

func (s *SyncManager) copyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo
	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		err := srcFd.Close()
		if err != nil {
			fmt.Println("Failed to close source file  %s ", srcFd.Name(), err)
		}
	}()
	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		err := dstFd.Close()
		if err != nil {
			fmt.Println("Failed to close destination file  %s ", dstFd.Name(), err)
		}
	}()
	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// delete all the template files
func (s *SyncManager) deleteTemplateFiles() {
	for _, f := range s.templateFiles {
		var _, err = os.Stat(f)
		if !os.IsNotExist(err) {
			var err = os.Remove(f)
			if err != nil {
				fmt.Println("Failed to delete template file %s ", f, err)
			}
		}
	}
}

// delete the temporary folder in which git clone was performed
func (s *SyncManager) deleteTempFolder() {
	err := os.RemoveAll(s.workingDir)
	if err != nil {
		fmt.Println("Failed to delete temporary folder %s ", s.workingDir, err)
	}
}
