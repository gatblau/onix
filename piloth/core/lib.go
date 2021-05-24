package core

/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"log"
	"os"
	"os/user"
	"path"
)

// HomeDir pilot's home directory
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// IsRegistered is the host registered?
func IsRegistered() bool {
	_, err := os.Stat(regpath())
	return os.IsExist(err)
}

// SetRegistered set the host as registered
func SetRegistered() error {
	regFile, err := os.Create(regpath())
	if err != nil {
		return err
	}
	regFile.Close()
	return nil
}

func regpath() string {
	return path.Join(HomeDir(), ".pilot_reg")
}
