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
type PipeCiCmd struct {
	cmd     *cobra.Command
	profile string
}

func NewPipeCiCmd() *PipeCiCmd {
	c := &PipeCiCmd{
		cmd: &cobra.Command{
			Use:   "ci [flags] [build-file-path]",
			Short: "deploy a CI pipeline for a specific application type to integrate application code using artie",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.profile, "profile", "p", "", "the build profile to use. if not provided, the default profile defined in the build file is used. if no default profile is found, then the first profile in the build file is used.")
	c.cmd.Run = c.Run
	return c
}

func (b *PipeCiCmd) Run(cmd *cobra.Command, args []string) {
	var buildFile = "."
	if len(args) == 1 {
		buildFile = args[0]
	} else if len(args) > 1 {
		core.RaiseErr("only one argument is required")
	}
	pipeline := tkn.NewPipeline(buildFile, b.profile)
	merged := new(bytes.Buffer)
	err := pipeline.Merge(merged)
	core.CheckErr(err, "cannot merge pipeline template")
	fmt.Print(merged)
}
