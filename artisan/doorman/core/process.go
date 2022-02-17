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
	"github.com/gatblau/onix/artisan/doorman/types"
)

// ProcessAsync starts processing a pipeline asynchronously
func ProcessAsync(ev types.NewSpecEvent) {
	go process(ev)
}

func process(ev types.NewSpecEvent) {
	fmt.Printf("start processing of event from %s %s\n", ev.URI, ev.Bucket)
}
