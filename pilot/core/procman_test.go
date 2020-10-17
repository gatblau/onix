/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"testing"
	"time"
)

var (
	homeDir    string
	appCmdArgs []string
)

const (
	appDir = "%s/go/src/github.com/gatblau/onix/probare"
	appCmd = "probare"
)

func init() {
	// determines the home directory
	hd, err := home()
	if err != nil {
		panic(err)
	}
	homeDir = hd
	// check the app to launch is there
	fqn := fmt.Sprintf("%s/%s", fmt.Sprintf(appDir, homeDir), appCmd)
	if _, err := os.Stat(fqn); os.IsNotExist(err) {
		panic(errors.New(fmt.Sprintf("%s application is needed in path %s", appCmd, fqn)))
	}
}

func TestTerminateApp(t *testing.T) {
	ps := NewProcessManager()
	// start the app process
	err := ps.start(fmt.Sprintf(appDir, homeDir), appCmd, appCmdArgs)
	check(t, err)
	// give the app time to fully start up
	time.Sleep(1 * time.Second)
	// request termination within 3 seconds
	err = ps.stop(3 * time.Second)
	check(t, err)
}

func TestRestartApp(t *testing.T) {
	ps := NewProcessManager()
	// start the app process
	err := ps.start(fmt.Sprintf(appDir, homeDir), appCmd, appCmdArgs)
	check(t, err)
	// give the app time to fully start up
	time.Sleep(1 * time.Second)
	// request restart with termination within 3 seconds
	err = ps.restart(3 * time.Second)
	check(t, err)
	// request termination within 1 second
	err = ps.stop(1 * time.Second)
}

// generic check for error
func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// get the user home
func home() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
