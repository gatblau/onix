/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"math/rand"
	"strings"
	"time"
)

// RandomPwd creates a new password of the specified length with or without symbols
func RandomPwd(length int, addSymbols bool) string {
	rand.Seed(time.Now().UnixNano())
	var chars []rune
	// if addSymbols is true then
	if addSymbols {
		// adds special characters to password character set
		chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_=!@£#$%&+")
	} else {
		// excludes special characters from password character set
		chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	}
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
