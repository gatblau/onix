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
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/i18n"
	"os"
	"path"
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
	keysPath := path.Join(core.RegistryPath(), "keys")
	// check the keys directory exists
	_, err = os.Stat(keysPath)
	// if it does not
	if os.IsNotExist(err) {
		// create a key pair
		err = os.Mkdir(keysPath, os.ModePerm)
		if err != nil {
			core.RaiseErr(err.Error())
		}
		host, _ := os.Hostname()
		crypto.GeneratePGPKeys(keysPath, "root", fmt.Sprintf("root-%s", host), "", "", versionLabel(), 2048)
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
	return fmt.Sprintf("onix-artisan-%s", Version)
}
