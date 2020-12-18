/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/tkn"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// list local artefacts
type PipeAppCmd struct {
	cmd         *cobra.Command
	profile     string
	sonar       *bool
	envFilename string
}

func NewPipeAppCmd() *PipeAppCmd {
	c := &PipeAppCmd{
		cmd: &cobra.Command{
			Use:   "app [flags] [build-file-path] [template_name]",
			Short: "deploy an artefact pipeline",
			Long:  `deploy a Tekton Pipeline to build an application package using Artisan`,
		},
	}
	c.cmd.Flags().StringVarP(&c.profile, "profile", "p", "", "the build profile to use. if not provided, the default profile defined in the build file is used. if no default profile is found, then the first profile in the build file is used.")
	c.sonar = c.cmd.Flags().BoolP("sonar", "s", false, "--sonar or -s add Sonar quality check step")
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
	c.cmd.Run = c.Run
	return c
}

func (b *PipeAppCmd) Run(cmd *cobra.Command, args []string) {
	var (
		buildFile    = "."
		templateName = "artefact_pipeline.yaml"
	)
	if len(args) == 1 {
		buildFile = args[0]
	}
	if len(args) == 2 {
		buildFile = args[0]
		templateName = args[1]
	} else if len(args) > 2 {
		core.RaiseErr("only two arguments are required")
	}
	if !path.IsAbs(templateName) {
		templateName, err := filepath.Abs(templateName)
		core.CheckErr(err, "cannot convert '%s' to absolute path", templateName)
	}
	// try to load env from file
	core.LoadEnvFromFile(b.envFilename)
	// collects information to assemble the pipeline
	c := tkn.NewAppPipelineConfig(buildFile, b.profile, *b.sonar)
	// assembles the pipeline
	template := tkn.MergeArtPipe(c, *b.sonar)
	core.CheckErr(ioutil.WriteFile(templateName, template.Bytes(), os.ModePerm), "cannot save template")
}
