/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/runner"
	"github.com/spf13/cobra"
)

// create a file seal
type RunCCmd struct {
	cmd         *cobra.Command
	interactive *bool
	envFilename string
}

func NewRunCCmd() *RunCCmd {
	c := &RunCCmd{
		cmd: &cobra.Command{
			Use:   "runc [function name] [project path]",
			Short: "runs the function commands specified in the project's build.yaml file within an artisan runtime container",
			Long:  ``,
		},
	}
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
	c.cmd.Run = c.Run
	return c
}

func (c *RunCCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		core.RaiseErr("At least function name is required")
	}
	var function = args[0]
	var path = "."
	if len(args) > 1 {
		path = args[1]
	}
	// load environment variables from file, if file not specified then try loading .env
	core.LoadEnvFromFile(c.envFilename)
	// create an instance of the runner
	run, err := runner.NewFromPath(path)
	core.CheckErr(err, "cannot initialise runner")
	// launch a runtime to execute the function
	err = run.RunC(function, *c.interactive)
	core.CheckErr(err, "cannot execute function '%s'", function)
}
