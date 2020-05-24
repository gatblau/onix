//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

func InitialiseRootCmd() *RootCmd {
	rootCmd := NewRootCmd()
	releaseCmd := InitialiseReleaseCmd(rootCmd)
	rootCmd.Command.AddCommand(releaseCmd.cmd)
	return rootCmd
}

func InitialiseReleaseCmd(rootCmd *RootCmd) *ReleaseCmd {
	releaseCmd := NewReleaseCmd(rootCmd.cfg)
	releaseInfoCmd := NewReleaseInfoCmd()
	releasePlanCmd := NewReleasePlanCmd()
	releaseCmd.cmd.AddCommand(releaseInfoCmd.cmd, releasePlanCmd.cmd)
	return releaseCmd
}
