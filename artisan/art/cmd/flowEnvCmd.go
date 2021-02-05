/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// list local artefacts
type FlowEnvCmd struct {
	cmd           *cobra.Command
	buildFilePath string
	stdout        *bool
	out           string
	flowPath      string
}

func NewFlowEnvCmd() *FlowEnvCmd {
	c := &FlowEnvCmd{
		cmd: &cobra.Command{
			Use:   "env [flags] [/path/to/flow_bare.yaml]",
			Short: "return the variables required by a given flow and can include a build.yaml",
			Long:  `return the variables required by a given flow and can include a build.yaml`,
		},
	}
	c.cmd.Flags().StringVarP(&c.buildFilePath, "build-file-path", "b", "", "--build-file-path=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.cmd.Flags().StringVarP(&c.out, "output", "o", "env", "--output yaml or -o yaml; the output format (e.g. env, json, yaml)")
	c.stdout = c.cmd.Flags().Bool("stdout", false, "prints the output to the console")
	c.cmd.Run = c.Run
	return c
}

func (c *FlowEnvCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		c.flowPath = core.ToAbsPath(args[0])
	} else if len(args) < 1 {
		core.RaiseErr("insufficient arguments: need the path to the bare flow file")
	} else if len(args) > 1 {
		core.RaiseErr("too many arguments: only need the path to the bare flow file")
	}
	// loads a bare flow from the path
	f, err := flow.LoadFlow(c.flowPath)
	core.CheckErr(err, "cannot load bare flow")

	// loads the build.yaml
	var b *data.BuildFile
	// if there is a build file, load it
	if len(c.buildFilePath) > 0 {
		b, err = data.LoadBuildFile(path.Join(c.buildFilePath, "build.yaml"))
	}
	// discover the input required by the flow / build file
	input := f.GetCombinedInput(b)
	var output []byte
	switch strings.ToLower(c.out) {
	// if the requested format is env
	case "env":
		output = toEnv(input)
	case "yaml":
		output, err = yaml.Marshal(input)
		core.CheckErr(err, "cannot marshal input")
	case "json":
		output, err = json.MarshalIndent(input, "", " ")
		core.CheckErr(err, "cannot marshal input")
	}
	if *c.stdout {
		// print to console
		fmt.Println(string(output))
	} else {
		// save to disk
		dir := filepath.Dir(c.flowPath)
		var filename string
		switch strings.ToLower(c.out) {
		case "yaml":
			fallthrough
		case "yml":
			filename = "env.yaml"
		case "json":
			filename = "env.json"
		default:
			filename = ".env"
		}
		err := ioutil.WriteFile(path.Join(dir, filename), output, os.ModePerm)
		core.CheckErr(err, "cannot write '%s' file", filename)
	}
}

func toEnv(i *data.Input) []byte {
	buf := &bytes.Buffer{}
	for _, v := range i.Var {
		buf.WriteString(fmt.Sprintf("# %s\n", v.Description))
		if len(v.Default) > 0 {
			buf.WriteString(fmt.Sprintf("%s=%s\n", v.Name, v.Default))
		} else {
			buf.WriteString(fmt.Sprintf("%s=\n", v.Name))
		}
	}
	return buf.Bytes()
}
