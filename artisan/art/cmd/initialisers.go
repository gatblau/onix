/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

func InitialiseRootCmd() *RootCmd {
	rootCmd := NewRootCmd()
	buildCmd := NewBuildCmd()
	artefactsCmd := NewListCmd()
	pushCmd := NewPushCmd()
	rmCmd := NewRmCmd()
	tagCmd := NewTagCmd()
	runCmd := NewRunCmd()
	runCCmd := NewRunCCmd()
	mergeCmd := NewMergeCmd()
	pullCmd := NewPullCmd()
	openCmd := NewOpenCmd()
	pgpCmd := InitialisePGPCommand()
	flowCmd := InitialiseFlowCommand()
	tknCmd := InitialiseTknCommand()
	manifCmd := NewManifestCmd()
	exeCmd := NewExeCmd()
	exeCCmd := NewExeCCmd()
	waitCmd := NewWaitCmd()
	langCmd := InitialiseLangCommand()
	rootCmd.Cmd.AddCommand(
		buildCmd.cmd,
		artefactsCmd.cmd,
		pushCmd.cmd,
		rmCmd.cmd,
		tagCmd.cmd,
		runCmd.cmd,
		runCCmd.cmd,
		mergeCmd.cmd,
		pullCmd.cmd,
		openCmd.cmd,
		pgpCmd.cmd,
		flowCmd.cmd,
		manifCmd.cmd,
		exeCmd.cmd,
		exeCCmd.cmd,
		tknCmd.cmd,
		waitCmd.cmd,
		langCmd.cmd,
	)
	return rootCmd
}

func InitialiseLangCommand() *LangCmd {
	langCmd := NewLangCmd()
	langFetchCmd := NewLangFetchCmd()
	langUpdateCmd := NewLangUpdateCmd()
	langCmd.cmd.AddCommand(langFetchCmd.cmd, langUpdateCmd.cmd)
	return langCmd
}

func InitialiseFlowCommand() *FlowCmd {
	flowCmd := NewFlowCmd()
	flowMergeCmd := NewFlowMergeCmd()
	flowEnvCmd := NewFlowEnvCmd()
	flowRunCmd := NewFlowRunCmd()
	flowCmd.cmd.AddCommand(flowMergeCmd.cmd, flowEnvCmd.cmd, flowRunCmd.cmd)
	return flowCmd
}

func InitialiseTknCommand() *TknCmd {
	tknCmd := NewTknCmd()
	tknGenCmd := NewTknGenCmd()
	tknCmd.cmd.AddCommand(tknGenCmd.cmd)
	return tknCmd
}

func InitialisePGPCommand() *PGPCmd {
	pgpCmd := NewPGPCmd()
	pgpGenCmd := NewPGPGenCmd()
	pgpImportCmd := NewPGPImportCmd()
	pgpExportCmd := NewPGPExportCmd()
	pgpEncryptCmd := NewPGPEncryptCmd()
	pgpDecryptCmd := NewPGPDecryptCmd()
	pgpCmd.cmd.AddCommand(pgpGenCmd.cmd, pgpImportCmd.cmd, pgpExportCmd.cmd, pgpEncryptCmd.cmd, pgpDecryptCmd.cmd)
	return pgpCmd
}
