/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

// list local artefacts
type FlowFillCmd struct {
	cmd           *cobra.Command
	envFilename   string
	buildFilePath string
}

func NewFlowFillCmd() *FlowFillCmd {
	c := &FlowFillCmd{
		cmd: &cobra.Command{
			Use:   "fill [flags] [/path/to/flow.yaml]",
			Short: "fills in a bare flow by adding the required variables, secrets and keys",
			Long:  `fills in a bare flow by adding the required variables, secrets and keys`,
		},
	}
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.cmd.Flags().StringVarP(&c.buildFilePath, "build-file", "b", "", "--build-file=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.cmd.Run = c.Run
	return c
}

func (c *FlowFillCmd) Run(cmd *cobra.Command, args []string) {
	var flowPath string
	switch len(args) {
	case 0:
		flowPath = ""
	case 1:
		flowPath = args[0]
		if !path.IsAbs(flowPath) {
			abs, err := filepath.Abs(flowPath)
			core.CheckErr(err, "cannot convert '%s' to absolute path", flowPath)
			flowPath = abs
		}
	default:
		core.RaiseErr("too many arguments")
	}
	// try to load env from file
	core.LoadEnvFromFile(c.envFilename)
	// loads a bare flow from the path
	g, err := flow.NewFromPath(flowPath, c.buildFilePath)
	core.CheckErr(err, "failed to load bare flow")
	// fills in the bare flow
	g.FillIn()
	// marshals the merged flow to a yaml string
	yaml, err := g.YamlString()
	core.CheckErr(err, "cannot fill in bare flow")
	// prints the flow to stdout
	fmt.Println(yaml)
}
