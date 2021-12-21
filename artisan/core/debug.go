package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
)

// Debug writes a debug message to the console
func Debug(msg string, a ...interface{}) {
	if InDebugMode() {
		DebugLogger.Printf("%s\n", fmt.Sprintf(msg, a...))
	}
}

// check for a ARTISAN_DEBUG variable set
func InDebugMode() bool {
	return len(os.Getenv("ARTISAN_DEBUG")) > 0
}
