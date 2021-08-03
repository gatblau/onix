package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

type VersionCmd struct {
	cmd *cobra.Command
}

func NewVersionCmd() *VersionCmd {
	c := &VersionCmd{
		cmd: &cobra.Command{
			Use:   "version",
			Short: "vmdiskfs utility version",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *VersionCmd) Run(cmd *cobra.Command, args []string) {
	log.Println("Version: 0.0.1")
}
