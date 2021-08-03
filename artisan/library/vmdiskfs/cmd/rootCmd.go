package cmd

import "github.com/spf13/cobra"

type RootCmd struct {
	Cmd *cobra.Command
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		Cmd: &cobra.Command{
			Use:   "vmdiskfs",
			Short: "Artisan: the Onix DevOps CLI",
			Long: `
++++++++++++++++++++++++++++++++++++++++++++++++
| Tool for parsing notifications, downloading, |
| converting and publishing to object storage  |
| Tool is part of artisan vmdiskfs runtime     |
++++++++++++++++++++++++++++++++++++++++++++++++
`,
		},
	}
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootCmd) initConfig() {

}
