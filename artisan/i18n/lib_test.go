/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package i18n

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	str := String("", LBL_LS_HEADER)
	fmt.Print(str)
	Printf(INFO_PUSHED, "aaa/bbb")
}
