/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"fmt"
	"os/user"
	"testing"
	"time"
)

func TestLaunchAndStopApp(t *testing.T) {
	ps := new(procMan)
	homeDir, err := home()
	check(t, err)
	err = ps.start(fmt.Sprintf("%s/go/src/github.com/gatblau/onix/probare", homeDir), "probare", []string{})
	check(t, err)
	time.Sleep(1 * time.Second)
	err = ps.stop(5 * time.Second)
	check(t, err)
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func home() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
