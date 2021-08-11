package cmd

//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
import (
	"fmt"
	"github.com/gatblau/onix/dbman/core"
	"os"
	"testing"
)

func TestCheck(t *testing.T) {
	dm, err := core.NewDbMan()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	core.DM = dm
	results := core.DM.CheckConfigSet()
	for check, result := range results {
		fmt.Printf("[%v] => %v\n", check, result)
	}
}
