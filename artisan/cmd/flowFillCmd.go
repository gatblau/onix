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
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// list local artefacts
type FlowFillCmd struct {
	cmd           *cobra.Command
	envFilename   string
	buildFilePath string
	stdout        *bool
}

func NewFlowFillCmd() *FlowFillCmd {
	c := &FlowFillCmd{
		cmd: &cobra.Command{
			Use: "fill [flags] [/path/to/flow.yaml] [path/to/pgp/public/key]",
			Short: "fills in a bare flow by adding the required variables, secrets and keys.\n" +
				"Secrets and keys are PGP encrypted by default using the provided public PGP key.",
			Long: `fills in a bare flow by adding the required variables, secrets and keys.\n
Secrets and keys are PGP encrypted by default using the provided public PGP key.`,
		},
	}
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.cmd.Flags().StringVarP(&c.buildFilePath, "build-file", "b", "", "--build-file=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.stdout = c.cmd.Flags().Bool("stdout", false, "--stdout to print the flow to the stdout")
	c.cmd.Run = c.Run
	return c
}

func (c *FlowFillCmd) Run(cmd *cobra.Command, args []string) {
	var flowPath, pubPath string
	if len(args) == 2 {
		flowPath = core.ToAbsPath(args[0])
		pubPath = core.ToAbsPath(args[1])
	} else if len(args) < 2 {
		core.RaiseErr("insufficient arguments: need the paths to the flow the PUBLIC PGP key files")
	} else if len(args) > 2 {
		core.RaiseErr("insufficient arguments: only need the paths to the flow the PUBLIC PGP key files")
	}
	// try to load env from file
	core.LoadEnvFromFile(c.envFilename)
	// loads a bare flow from the path
	g, err := flow.NewFromPath(flowPath, pubPath, c.buildFilePath)
	core.CheckErr(err, "failed to load bare flow")
	// fills in the bare flow
	g.FillIn(registry.NewLocalRegistry())
	// marshals the merged flow to a yaml string
	yaml, err := g.YamlString()
	core.CheckErr(err, "cannot fill in bare flow")
	if *c.stdout {
		// prints the flow to stdout
		fmt.Println(yaml)
	} else {
		dir, file := filepath.Split(flowPath)
		filename := core.FilenameWithoutExtension(file)
		mergedFlowFilename := filepath.Join(dir, fmt.Sprintf("%s-%s.yaml", filename, core.RandomString(5)))
		err = ioutil.WriteFile(mergedFlowFilename, []byte(yaml), os.ModePerm)
		core.CheckErr(err, "cannot save filled in flow")
	}
}
