/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package crypto

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func TestLoader(t *testing.T) {
	l := newKeyLoader()
	for _, k := range l.m {
		fmt.Printf("%s\n", k)
	}
}

func TestLoadPrivate(t *testing.T) {
	name, _ := core.ParseName("localhost:8082/test/testpk/aabb:latest")
	l := newKeyLoader()
	k, err := l.Key(name, true)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	fmt.Printf("Key Id: %d\n", k.entity.PrimaryKey.KeyId)
}
