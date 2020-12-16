/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"github.com/gatblau/onix/artisan/cmd"
)

func main() {
	rootCmd := cmd.InitialiseRootCmd()

	// Execute adds all child commands to the root command and sets flags appropriately.
	// This is called by main.main(). It only needs to happen once to the rootCmd.
	rootCmd.Cmd.Execute()
}
