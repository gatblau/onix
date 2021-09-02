package build

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/artisan/merge"
	"testing"
)

func TestExe(t *testing.T) {
	out, err := Exe("printenv", ".", merge.NewEnVarFromSlice([]string{}), false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
}
