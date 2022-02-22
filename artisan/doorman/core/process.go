/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
)

// ProcessAsync starts processing a pipeline asynchronously
func ProcessAsync(uri string) {
	go process(uri)
}

func process(uri string) {
	fmt.Printf("start processing of event from %s\n", uri)
}
