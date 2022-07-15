/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/tkn"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FlowMergeCmd merge a flow with env variables
type FlowMergeCmd struct {
	Cmd           *cobra.Command
	envFilename   string
	buildFilePath string
	stdout        *bool
	tkn           *bool
	out           string
	interactive   *bool
	labels        []string
}

func NewFlowMergeCmd() *FlowMergeCmd {
	c := &FlowMergeCmd{
		Cmd: &cobra.Command{
			Use:   "merge [flags] [/path/to/flow_bare.yaml]",
			Short: "fills in a bare flow by adding the required variables, secrets and keys",
			Long:  `fills in a bare flow by adding the required variables, secrets and keys`,
		},
	}
	c.Cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.Cmd.Flags().StringVarP(&c.buildFilePath, "build-file-path", "b", "", "--build-file-path=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.stdout = c.Cmd.Flags().Bool("stdout", false, "prints the output to the console")
	c.tkn = c.Cmd.Flags().Bool("tkn", false, "generates a tekton resources file")
	c.Cmd.Flags().StringVarP(&c.out, "output", "o", "yaml", "--output json or -o json; the output format for the written flow; available formats are:\n"+
		"yaml: output in YAML format\n"+
		"json: output in JSON format\n"+
		"ojason: output as an Onix configuration item format\n")
	c.interactive = c.Cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.Cmd.Flags().StringSliceVarP(&c.labels, "label", "l", []string{}, "add one or more labels to the flow; -l label1=value1 -l label2=value2")
	c.Cmd.Run = c.Run
	return c
}

func (c *FlowMergeCmd) Run(_ *cobra.Command, args []string) {
	var flowPath string
	if len(args) == 1 {
		flowPath = core.ToAbsPath(args[0])
	} else if len(args) < 1 {
		core.RaiseErr("insufficient arguments: need the path to the bare flow file")
	} else if len(args) > 1 {
		core.RaiseErr("too many arguments: only need the path to the bare flow file")
	}
	// add the build file level environment variables
	env := merge.NewEnVarFromSlice(os.Environ())
	// load vars from file
	env2, err := merge.NewEnVarFromFile(c.envFilename)
	core.CheckErr(err, "failed to load environment file '%s'", c.envFilename)
	// merge with existing environment
	env.Merge(env2)
	// loads a bare flow from the path
	f, err := flow.NewWithEnv(flowPath, c.buildFilePath, env, "")
	core.CheckErr(err, "cannot load bare flow")
	// add labels to the flow
	f.AddLabels(c.labels)
	// merges input, surveying for required data if in interactive mode
	err = f.Merge(*c.interactive)
	core.CheckErr(err, "cannot merge bare flow")
	// if tekton format is requested
	if *c.tkn {
		// gets a tekton transpiler
		builder := tkn.NewBuilder(f.Flow)
		// transpile the flow
		buf := builder.BuildBuffer()
		// if stdout required
		if *c.stdout {
			// print to console
			fmt.Println(buf.String())
		} else {
			// write to file
			err = ioutil.WriteFile(tknPath(flowPath), buf.Bytes(), os.ModePerm)
			core.CheckErr(err, "cannot write tekton file")
		}
	} else { // flow format requested
		// if stdout required
		if *c.stdout {
			if c.out == "yaml" {
				// marshals the flow to YAML
				yaml, err := f.YamlString()
				core.CheckErr(err, "cannot marshal bare flow")
				// print to stdout
				fmt.Println(yaml)
			} else if c.out == "json" {
				// marshals the flow to YAML
				json, err := f.JsonString()
				core.CheckErr(err, "cannot marshal bare flow")
				// print to stdout
				fmt.Println(json)
			} else {
				core.RaiseErr("invalid format '%s'", c.out)
			}
		} else {
			// save the flow to file
			if c.out == "yaml" {
				err = f.SaveYAML()
			} else if c.out == "json" {
				err = f.SaveJSON()
			} else if c.out == "ojson" {
				err = f.SaveOnixJSON()
			} else {
				core.RaiseErr("invalid format '%s'", c.out)
			}
			core.CheckErr(err, "cannot save bare flow")
		}
	}
}

func tknPath(path string) string {
	dir, file := filepath.Split(path)
	filename := core.FilenameWithoutExtension(file)
	return filepath.Join(dir, fmt.Sprintf("%s_tkn.yaml", filename[0:len(filename)-len("_bare")]))
}
