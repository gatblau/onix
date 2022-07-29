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
	utilCmd := InitialiseUtilCommand()
	specCmd := InitialiseSpecCommand()
	buildCmd := NewBuildCmd()
	lsCmd := NewListCmd()
	pushCmd := NewPushCmd()
	rmCmd := NewRmCmd()
	tagCmd := NewTagCmd()
	runCmd := NewRunCmd()
	runCCmd := NewRunCCmd()
	mergeCmd := NewMergeCmd()
	pullCmd := NewPullCmd()
	openCmd := NewOpenCmd()
	flowCmd := InitialiseFlowCommand()
	manifCmd := NewManifestCmd()
	exeCmd := NewExeCmd()
	exeCCmd := NewExeCCmd()
	envCmd := InitialiseEnvCommand()
	pruneCmd := NewPruneCmd()
	rootCmd.Cmd.AddCommand(
		utilCmd.Cmd,
		specCmd.Cmd,
		buildCmd.Cmd,
		lsCmd.Cmd,
		pushCmd.Cmd,
		rmCmd.Cmd,
		tagCmd.Cmd,
		runCmd.Cmd,
		runCCmd.Cmd,
		mergeCmd.Cmd,
		pullCmd.Cmd,
		openCmd.cmd,
		flowCmd.Cmd,
		manifCmd.Cmd,
		exeCmd.cmd,
		exeCCmd.Cmd,
		envCmd.Cmd,
		pruneCmd.Cmd,
	)
	return rootCmd
}

func InitialiseUtilCommand() *UtilCmd {
	utilCmd := NewUtilCmd()
	utilPwdCmd := NewUtilPwdCmd()
	utilNameCmd := NewUtilNameCmd()
	utilExtractCmd := NewUtilExtractCmd()
	utilB64Cmd := NewUtilBase64Cmd()
	utilStampCmd := NewUtilStampCmd()
	curlCmd := NewCurlCmd()
	waitCmd := NewWaitCmd()
	tknCmd := InitialiseTknCommand()
	langCmd := InitialiseLangCommand()
	gitSyncCmd := NewGitSyncCmd()
	serveCmd := NewServeCmd()
	appCmd := NewAppCmd()
	utilCmd.Cmd.AddCommand(
		utilPwdCmd.Cmd,
		utilExtractCmd.Cmd,
		utilNameCmd.Cmd,
		utilB64Cmd.Cmd,
		utilStampCmd.Cmd,
		waitCmd.Cmd,
		curlCmd.Cmd,
		tknCmd.Cmd,
		langCmd.Cmd,
		gitSyncCmd.Cmd,
		serveCmd.Cmd,
		appCmd.Cmd,
	)
	return utilCmd
}

func InitialiseSpecCommand() *SpecCmd {
	specCmd := NewSpecCmd()
	specExportCmd := NewSpecExportCmd()
	specImportCmd := NewSpecImportCmd()
	specDownCmd := NewSpecDownCmd()
	specUpCmd := NewSpecUpCmd()
	specPushCmd := NewSpecPushCmd()
	specInfoCmd := NewSpecInfoCmd()
	specPullCmd := NewSpecPullCmd()
	specCmd.Cmd.AddCommand(specExportCmd.Cmd)
	specCmd.Cmd.AddCommand(specImportCmd.Cmd)
	specCmd.Cmd.AddCommand(specDownCmd.Cmd)
	specCmd.Cmd.AddCommand(specUpCmd.Cmd)
	specCmd.Cmd.AddCommand(specPushCmd.Cmd)
	specCmd.Cmd.AddCommand(specInfoCmd.Cmd)
	specCmd.Cmd.AddCommand(specPullCmd.Cmd)
	return specCmd
}

func InitialiseEnvCommand() *EnvCmd {
	envCmd := NewEnvCmd()
	envPackageCmd := NewEnvPackageCmd()
	envFlowCmd := NewEnvFlowCmd()
	envCmd.Cmd.AddCommand(envFlowCmd.Cmd, envPackageCmd.Cmd)
	return envCmd
}

func InitialiseLangCommand() *LangCmd {
	langCmd := NewLangCmd()
	langFetchCmd := NewLangFetchCmd()
	langUpdateCmd := NewLangUpdateCmd()
	langCmd.Cmd.AddCommand(langFetchCmd.Cmd, langUpdateCmd.Cmd)
	return langCmd
}

func InitialiseFlowCommand() *FlowCmd {
	flowCmd := NewFlowCmd()
	flowMergeCmd := NewFlowMergeCmd()
	flowRunCmd := NewFlowRunCmd()
	flowCmd.Cmd.AddCommand(flowMergeCmd.Cmd, flowRunCmd.Cmd)
	return flowCmd
}

func InitialiseTknCommand() *TknCmd {
	tknCmd := NewTknCmd()
	tknGenCmd := NewTknGenCmd()
	tknCmd.Cmd.AddCommand(tknGenCmd.Cmd)
	return tknCmd
}
