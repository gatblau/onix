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
	"github.com/gatblau/onix/artisan/tkn"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// list local artefacts
type FlowMergeCmd struct {
	cmd           *cobra.Command
	envFilename   string
	buildFilePath string
	stdout        *bool
	tkn           *bool
}

func NewFlowMergeCmd() *FlowMergeCmd {
	c := &FlowMergeCmd{
		cmd: &cobra.Command{
			Use:   "merge [flags] [/path/to/flow_bare.yaml]",
			Short: "fills in a bare flow by adding the required variables, secrets and keys",
			Long:  `fills in a bare flow by adding the required variables, secrets and keys`,
		},
	}
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.cmd.Flags().StringVarP(&c.buildFilePath, "build-file-path", "b", ".", "--build-file-path=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.stdout = c.cmd.Flags().Bool("stdout", false, "prints the output to the console")
	c.tkn = c.cmd.Flags().Bool("tkn", false, "generates a tekton resources file")
	c.cmd.Run = c.Run
	return c
}

func (c *FlowMergeCmd) Run(cmd *cobra.Command, args []string) {
	var flowPath string
	if len(args) == 1 {
		flowPath = core.ToAbsPath(args[0])
	} else if len(args) < 1 {
		core.RaiseErr("insufficient arguments: need the path to the bare flow file")
	} else if len(args) > 1 {
		core.RaiseErr("too many arguments: only need the path to the bare flow file")
	}
	// loads a bare flow from the path
	flow, err := flow.NewFromPath(flowPath, c.buildFilePath)
	core.CheckErr(err, "cannot load bare flow")
	// survey for required inputs
	err = flow.Merge()
	core.CheckErr(err, "cannot merge bare flow")
	// if tekton format is requested
	if *c.tkn {
		// gets a tekton transpiler
		builder := tkn.NewBuilder(flow.Flow)
		// transpile the flow
		buf := builder.Create()
		// if stdout required
		if *c.stdout {
			// print to console
			fmt.Println(buf)
		} else {
			// write to file
			err = ioutil.WriteFile(tknPath(flowPath), buf.Bytes(), os.ModePerm)
			core.CheckErr(err, "cannot write tekton file")
		}
	} else { // flow format requested
		// if stdout required
		if *c.stdout {
			// marshals the flow to YAML
			yaml, err := flow.YamlString()
			core.CheckErr(err, "cannot marshal bare flow")
			// print to stdout
			fmt.Println(yaml)
		} else {
			// save the flow to file
			err = flow.Save()
			core.CheckErr(err, "cannot save bare flow")
		}
	}
}

func tknPath(path string) string {
	dir, file := filepath.Split(path)
	filename := core.FilenameWithoutExtension(file)
	return filepath.Join(dir, fmt.Sprintf("%s_tkn.yaml", filename[0:len(filename)-len("_bare")]))
}
