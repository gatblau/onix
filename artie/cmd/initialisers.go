/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

func InitialiseRootCmd() *RootCmd {
	rootCmd := NewRootCmd()
	buildCmd := NewBuildCmd()
	artefactsCmd := NewArtefactsCmd()
	pushCmd := NewPushCmd()
	rmCmd := NewRmCmd()
	tagCmd := NewTagCmd()
	serveCmd := NewServeCmd()
	versionCmd := NewVersionCmd()
	rootCmd.Command.AddCommand(buildCmd.cmd, artefactsCmd.cmd, pushCmd.cmd, rmCmd.cmd, tagCmd.cmd, serveCmd.cmd, versionCmd.cmd)
	return rootCmd
}
