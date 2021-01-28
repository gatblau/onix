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

type ExeCCmd struct {
	cmd         *cobra.Command
	interactive *bool
	credentials string
	path        string
}

func NewExeCCmd() *ExeCCmd {
	c := &ExeCCmd{
		cmd: &cobra.Command{
			Use:   "exec [flags] [package-name] [function-name] [path/to/context]",
			Short: "runs a function within a package using an artisan runtime",
			Long: `runs a function within a package using an artisan runtime
* package-name: 
   mandatory - the fully qualified name of the package containing the function to execute
* function-name: 
   mandatory - the name of the function exported by the package that should be executed
* path/to/context: 
   optional - the path to the folder in the host where context files are located
   if not specified then current path is assumed.

NOTE: exec always pulls the package from its registry as it is done within the runtime and that is its only behaviour
   if the package is in a secure registry, then credentials must be specified via -u / --credentials flag
   if running in a linux host, ensure the user executing the exec command has UID/GID = 100000000 
   to avoid read / write issues from / to the host - e.g. public PGP key in the host artisan registry is required to open
     the package within the runtime - keys are accessible in the runtime using bind mounts
`,
		},
	}
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD artisan registry user and password")
	c.cmd.Run = c.Run
	return c
}

func (c *ExeCCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		core.RaiseErr("insufficient arguments")
	} else if len(args) > 2 {
		core.RaiseErr("too many arguments")
	}
	var (
		packageName = args[0]
		fxName      = args[1]
		path        = "."
	)
	if len(args) > 1 {
		path = args[1]
	}
	// create an instance of the runner
	run, err := runner.NewFromPath(path)
	core.CheckErr(err, "cannot initialise runner")
	// launch a runtime to execute the function
	err = run.ExeC(packageName, fxName, c.credentials, path, *c.interactive)
	core.CheckErr(err, "cannot execute function '%s' in package '%s'", fxName, packageName)
}
