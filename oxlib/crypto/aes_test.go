/*
  Onix Config Manager - crypto utils
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package crypto

import (
	"fmt"
	"testing"
)

func TestEncryptAES256(t *testing.T) {
	pt := "This is a secret"
	c := EncryptAES(pt)
	fmt.Println(pt)
	fmt.Println(c)
	pt2 := DecryptAES(c)
	fmt.Println(pt2)
}
