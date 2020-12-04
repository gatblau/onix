/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"github.com/gatblau/onix/artie/tkn"
	"github.com/spf13/cobra"
)

// list local artefacts
type PipeSetupCmd struct {
	cmd *cobra.Command
}

func NewPipeSetupCmd() *PipeSetupCmd {
	c := &PipeSetupCmd{
		cmd: &cobra.Command{
			Use:   "setup",
			Short: "setups the pipelines required to automatically build applications using artie",
			Long:  ``,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (b *PipeSetupCmd) Run(cmd *cobra.Command, args []string) {
	pipeline := tkn.NewPipeline("", "")
	merged := new(bytes.Buffer)
	err := pipeline.Merge(merged)
	core.CheckErr(err, "cannot merge pipeline template")
	fmt.Print(merged)
}
