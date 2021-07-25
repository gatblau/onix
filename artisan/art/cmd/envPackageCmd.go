package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// list local packages
type EnvPackageCmd struct {
	cmd           *cobra.Command
	buildFilePath string
	stdout        *bool
	out           string
}

func NewEnvPackageCmd() *EnvPackageCmd {
	c := &EnvPackageCmd{
		cmd: &cobra.Command{
			Use: "package [flags] [package name] [function-name (optional)]",
			Short: "return the variables required by a given package to run\n " +
				"if a function name is not specified then variables for all functions are retrieved",
			Long: "return the variables required by a given package to run\n " +
				"if a function name is not specified then variables for all functions are retrieved",
		},
	}
	c.cmd.Flags().StringVarP(&c.buildFilePath, "build-file-path", "b", "", "--build-file-path=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.cmd.Flags().StringVarP(&c.out, "output", "o", "env", "--output yaml or -o yaml; the output format (e.g. env, json, yaml)")
	c.stdout = c.cmd.Flags().Bool("stdout", false, "prints the output to the console")
	c.cmd.Run = c.Run
	return c
}

func (c *EnvPackageCmd) Run(cmd *cobra.Command, args []string) {
	var input *data.Input
	if len(args) > 0 && len(args) < 3 {
		name, err := core.ParseName(args[0])
		core.CheckErr(err, "invalid package name: %s", name)
		local := registry.NewLocalRegistry()
		manifest := local.GetManifest(name)
		if len(args) == 2 {
			fxName := args[1]
			fx := manifest.Fx(fxName)
			input = fx.Input
		} else {
			for i, function := range manifest.Functions {
				if i == 0 {
					input = function.Input
				} else {
					input.Merge(function.Input)
				}
			}
		}
		// add the credentials to download the package
		input.SurveyRegistryCreds(name.Group, name.Name, "", name.Domain, false, true, merge.NewEnVarFromSlice([]string{}))
	} else if len(args) < 2 {
		i18n.Raise(i18n.ERR_INSUFFICIENT_ARGS)
	} else if len(args) > 2 {
		i18n.Raise(i18n.ERR_TOO_MANY_ARGS)
	}

	var (
		output []byte
		err    error
	)
	switch strings.ToLower(c.out) {
	// if the requested format is env
	case "env":
		output = input.ToEnvFile()
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
		err = ioutil.WriteFile(filename, output, os.ModePerm)
		core.CheckErr(err, "cannot write '%s' file", filename)
	}
}
