package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/library/vmdiskfs/parser"
	"github.com/spf13/cobra"
	"os"
)

type ParseCmd struct {
	cmd *cobra.Command
}

func NewParseCmd() *ParseCmd {
	c := &ParseCmd{
		cmd: &cobra.Command{
			Use:   "parse",
			Short: "parses push notification received",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *ParseCmd) Run(cmd *cobra.Command, args []string) {
	provider := os.Getenv("AWS_PROVIDER")
	parse := parser.ParsedInformationWriter(provider)
	if parse == true {
		fmt.Println("Successfully Parsed")
	} else {
		fmt.Println("Failed, due to empty data")
	}
}
