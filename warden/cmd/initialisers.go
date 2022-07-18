/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

func InitialiseRootCmd() *RootCmd {
	rootCmd := NewRootCmd()
	launchCmd := InitialiseLaunchCommand()
	rootCmd.Cmd.AddCommand(
		launchCmd.cmd,
	)
	return rootCmd
}

func InitialiseLaunchCommand() *LaunchCmd {
	launchCmd := NewLaunchCmd()
	basicCmd := NewBasicCmd()
	tapCmd := NewTapCmd()
	launchCmd.cmd.AddCommand(basicCmd.cmd, tapCmd.cmd)
	return launchCmd
}
