/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func TestLock(t *testing.T) {
	name, _ := core.ParseName("gatblau/boot")
	repoName := name.Repository()
	lock := new(lock)
	lock.ensurePath()
	locked, err := lock.acquire(repoName)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	if locked < 0 {
		t.FailNow()
	}
	unlocked, err := lock.release(repoName)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	if unlocked < 1 {
		t.FailNow()
	}
}
