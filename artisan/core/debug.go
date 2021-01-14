package core

import (
	"fmt"
	"os"
)

// writes a debug message to the console
func Debug(msg string, a ...interface{}) {
	if inDebugMode() {
		fmt.Printf("DEBUG => %s\n", fmt.Sprintf(msg, a...))
	}
}

// check for a ARTISAN_DEBUG variable set
func inDebugMode() bool {
	return len(os.Getenv("ARTISAN_DEBUG")) > 0
}
