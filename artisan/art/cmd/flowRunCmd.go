/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/spf13/cobra"
	"os"
)

type FlowRunCmd struct {
	Cmd           *cobra.Command
	envFilename   string
	credentials   string
	interactive   *bool
	flowPath      string
	runnerName    string
	buildFilePath string
	labels        []string
}

func NewFlowRunCmd() *FlowRunCmd {
	c := &FlowRunCmd{
		Cmd: &cobra.Command{
			Use:   "run [flags] [/path/to/flow.yaml] [runner name]",
			Short: "merge and send a flow to a runner for execution",
			Long:  `merge and send a flow to a runner for execution`,
		},
	}
	c.Cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env; the path to a file containing environment variables to use")
	c.Cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.interactive = c.Cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.Cmd.Flags().StringVarP(&c.buildFilePath, "build-file-path", "b", ".", "--build-file-path=. or -b=.; the path to an artisan build.yaml file from which to pick required inputs")
	c.Cmd.Flags().StringSliceVarP(&c.labels, "label", "l", []string{}, "add one or more labels to the flow; -l label1=value1 -l label2=value2")
	c.Cmd.Run = c.Run
	return c
}

func (c *FlowRunCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 2 {
		c.flowPath = core.ToAbsPath(args[0])
		c.runnerName = args[1]
	} else if len(args) < 1 {
		i18n.Raise("", i18n.ERR_INSUFFICIENT_ARGS)
	} else if len(args) > 1 {
		i18n.Raise("", i18n.ERR_TOO_MANY_ARGS)
	}
	// add the build file level environment variables
	env := merge.NewEnVarFromSlice(os.Environ())
	// load vars from file
	env2, err := merge.NewEnVarFromFile(c.envFilename)
	core.CheckErr(err, "failed to load environment file '%s'", c.envFilename)
	// merge with existing environment
	env.Merge(env2)
	// loads a flow from the path
	f, err := flow.NewWithEnv(c.flowPath, c.buildFilePath, env, "")
	core.CheckErr(err, "cannot load flow")
	// add labels to the flow
	f.AddLabels(c.labels)
	err = f.Merge(*c.interactive)
	core.CheckErr(err, "cannot merge flow")
	err = f.Run(c.runnerName, c.credentials, *c.interactive)
	core.CheckErr(err, "cannot run flow")
}
