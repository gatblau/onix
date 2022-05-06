/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"os"
	"runtime"
)

func init() {
	// ensure the registry folder structure is in place
	ensureRegistryDir()
}

// check the local localReg directory exists and if not creates it
func ensureRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(core.RegistryPath())
	// if it does not
	if os.IsNotExist(err) {
		if runtime.GOOS == "linux" && os.Geteuid() == 0 {
			core.WarningLogger.Printf("if the root user creates the local registry then runc commands will fail\n" +
				"as the runtime user will not be able to access its content when it is bind mounted\n" +
				"ensure the local registry path is not owned by the root user\n")
		}
		err = os.Mkdir(core.RegistryPath(), os.ModePerm)
		i18n.Err(err, i18n.ERR_CANT_CREATE_REGISTRY_FOLDER, core.RegistryPath(), core.HomeDir())
	}
	filesPath := core.FilesPath()
	// check the files directory exists
	_, err = os.Stat(filesPath)
	// if it does not
	if os.IsNotExist(err) {
		// create a key pair
		err = os.Mkdir(filesPath, os.ModePerm)
		if err != nil {
			core.RaiseErr(err.Error())
		}
	}
}

func versionLabel() string {
	return fmt.Sprintf("onix-artisan-%s", core.Version)
}
