/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
	"os"
)

// create a file seal
type RunCmd struct {
	cmd         *cobra.Command
	envFilename string
	interactive *bool
}

func NewRunCmd() *RunCmd {
	c := &RunCmd{
		cmd: &cobra.Command{
			Use:   "run [function name] [project path]",
			Short: "runs the function commands specified in the project's build.yaml file",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.cmd.Run = c.Run
	return c
}

func (c *RunCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		core.RaiseErr("At least function name is required")
	}
	var function = args[0]
	var path = "."
	if len(args) > 1 {
		path = args[1]
	}
	builder := build.NewBuilder()
	// add the build file level environment variables
	env := core.NewEnVarFromSlice(os.Environ())
	// load vars from file
	env2, _ := core.NewEnVarFromFile(c.envFilename)
	// merge with existing environment
	env.Merge(env2)
	// execute the function
	builder.Run(function, path, *c.interactive, env)
}
