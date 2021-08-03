package cmd

func InitialiseRootCmd() *RootCmd {
	rootCmd := NewRootCmd()
	parseCmd := NewParseCmd()
	downloadCmd := NewDownloadCmd()
	uploadCmd := NewUploadCmd()
	convertCmd := NewConvertCmd()
	downloadOptCmd := NewDownloadByConfigCmd()
	uploadOptCmd := NewUploadByConfigCmd()
	convertOptCmd := NewConvertByConfigCmd()
	downloadConfigCmd := NewDownloadConfigCmd()
	sendWebhookCmd := NewSendWebHookCmd()
	version := NewVersionCmd()
	rootCmd.Cmd.AddCommand(
		parseCmd.cmd,
		downloadCmd.cmd,
		uploadCmd.cmd,
		convertCmd.cmd,
		downloadOptCmd.cmd,
		uploadOptCmd.cmd,
		convertOptCmd.cmd,
		downloadConfigCmd.cmd,
		sendWebhookCmd.cmd,
		version.cmd,
	)
	return rootCmd
}
