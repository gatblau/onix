/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const AppName = "artisan"

// HomeDir gets the user home directory
func HomeDir() string {
	// if ARTISAN_HOME is defined use it
	if artHome := os.Getenv("ARTISAN_HOME"); len(artHome) > 0 {
		return artHome
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func WorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}

// RegistryPath gets the root path of the local registry
func RegistryPath(path string) string {
	if len(path) > 0 {
		path, _ = filepath.Abs(path)
		return filepath.Join(path, fmt.Sprintf(".%s", AppName))
	}
	return filepath.Join(HomeDir(), fmt.Sprintf(".%s", AppName))
}

func KeysPath(path string) string {
	return filepath.Join(RegistryPath(path), "keys")
}

func FilesPath(path string) string {
	return filepath.Join(RegistryPath(path), "files")
}

func LangPath(path string) string {
	return filepath.Join(RegistryPath(path), "lang")
}

// TmpPath temporary path for file operations
func TmpPath(path string) string {
	return filepath.Join(RegistryPath(path), "tmp")
}

func TmpExists(path string) {
	tmp := TmpPath(path)
	// ensure tmp folder exists for temp file operations
	_, err := os.Stat(tmp)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(tmp, os.ModePerm)
	}
}

func LangExists(path string) {
	lang := LangPath(path)
	// ensure lang folder exists for temp file operations
	_, err := os.Stat(lang)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(lang, os.ModePerm)
	}
}

// RunPath temporary path for running package functions
func RunPath(path string) string {
	return filepath.Join(RegistryPath(path), "tmp", "run")
}

func RunPathExists(path string) {
	runPath := RunPath(path)
	// ensure tmp folder exists for  running package functions
	_, err := os.Stat(runPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(runPath, os.ModePerm)
	}
}

// EnsureRegistryPath check the local registry directory exists and if not creates it
func EnsureRegistryPath(path string) error {
	// check the home directory exists
	_, err := os.Stat(RegistryPath(path))
	// if it does not
	if os.IsNotExist(err) {
		if runtime.GOOS == "linux" && os.Geteuid() == 0 {
			WarningLogger.Printf("if the root user creates the local registry then runc commands will fail\n" +
				"as the runtime user will not be able to access its content when it is bind mounted\n" +
				"ensure the local registry path is not owned by the root user\n")
		}
		err = os.MkdirAll(RegistryPath(path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create registry folder: %s\n", err)
		}
	}
	filesPath := FilesPath(path)
	// check the files' directory exists
	_, err = os.Stat(filesPath)
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(filesPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create local registry files folder: %s\n", err)
		}
	}
	return nil
}
