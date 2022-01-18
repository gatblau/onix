/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	os "os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// a read / write lock for repository metadata operations
type lock struct {
}

// acquires a read/write lock for the specified repository
func (l *lock) acquire(repository string) (int, error) {
	lockName := l.name(repository)
	_, err := os.Stat(lockName)
	// if the lock does not exist we are good to go
	if os.IsNotExist(err) {
		// create it
		_, err := os.Create(lockName)
		// if the creation failed
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	// if we got here is because the lock already exist therefore
	// a new lock cannot be obtained until the existing one is released
	return -1, nil
}

// releases the lock for the specified repository
func (l *lock) release(repository string) (int, error) {
	lockName := l.name(repository)
	// if the lock exists
	if l.fileExists(lockName) {
		// delete the lock file
		err := os.Remove(lockName)
		// if it failed to delete the lock file
		if err != nil {
			// return the error
			return 0, err
		}
	}
	return 1, nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func (l *lock) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// returns the name of the lock for a repository
func (l *lock) name(repository string) string {
	filename := strings.ReplaceAll(core.StringCheckSum(repository), "/", "")
	return l.fqn(fmt.Sprintf("%s.lock", filename))
}

// wait until a lock can be released and then release it
func (l *lock) tryRelease(repository string, lockLifespanInSeconds time.Duration) error {
	lockName := l.name(repository)
	info, err := os.Stat(lockName)
	// if the lock file does not exist returns unlocked
	if os.IsNotExist(err) {
		return nil
	}
	for {
		// if the lock can be released (i.e. the current time is greater than the last mod time of the lock file plus lockLifespanInSceonds)
		if time.Now().After(info.ModTime().Add(lockLifespanInSeconds * time.Second)) {
			released, msg := l.release(repository)
			if released > 0 {
				return nil
			}
			return fmt.Errorf("cannot release lock: %s", msg)
		}
	}
}

// the path where lock files are written
func (l *lock) path() string {
	return filepath.Join(core.HomeDir(), fmt.Sprintf(".%s", core.AppName), "locks")
}

// ensures the locks path is there
func (l *lock) ensurePath() {
	_, err := os.Stat(l.path())
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(l.path(), os.ModePerm)
			u, _ := user.Current()
			core.CheckErr(err, fmt.Sprintf("cannot create locks path, user='%s', uid='%s', path='%s'", u.Username, u.Uid, l.path()))
		}
	}
}

func (l *lock) fqn(filename string) string {
	return filepath.Join(l.path(), filename)
}
