/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	// names with numbers
	for i := 0; i < 20; i++ {
		fmt.Println(RandomName(99))
	}
	// names with no numbers
	for i := 0; i < 20; i++ {
		fmt.Println(RandomName(0))
	}
}
