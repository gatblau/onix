/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"strings"
)

// encode strings to be used in tekton pipelines names
func encode(value string) string {
	length := 30
	value = strings.ToLower(value)
	value = strings.Replace(value, " ", "-", -1)
	if len(value) > length {
		value = value[0:length]
	}
	return value
}
