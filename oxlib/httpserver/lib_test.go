/*
  Onix Config Manager - Http Client
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package httpserver

import (
	"fmt"
	"testing"
)

func TestBasicToken(t *testing.T) {
	token := BasicToken("ab.cd", "1234")
	uname, pwd := ReadBasicToken(token)
	fmt.Println(uname, pwd)
}
