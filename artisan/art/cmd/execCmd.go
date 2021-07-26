package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/runner"
	"github.com/spf13/cobra"
)

type ExeCCmd struct {
	cmd         *cobra.Command
	interactive *bool
	credentials string
	path        string
	envFilename string
}

func NewExeCCmd() *ExeCCmd {
	c := &ExeCCmd{
		cmd: &cobra.Command{
			Use:   "exec [flags] [package-name] [function-name]",
			Short: "runs a function within a package using an artisan runtime",
			Long: `runs a function within a package using an artisan runtime
* package-name: 
   mandatory - the fully qualified name of the package containing the function to execute
* function-name: 
   mandatory - the name of the function exported by the package that should be executed

NOTE: exec always pulls the package from its registry as it is done within the runtime and that is its only behaviour
   if the package is in a secure registry, then credentials must be specified via -u / --credentials flag
   if running in a linux host, ensure the user executing the exec command has UID/GID = 100000000 
   to avoid read / write issues from / to the host - e.g. public PGP key in the host artisan registry is required to open
     the package within the runtime - keys are accessible within the runtime using bind mounts
`,
		},
	}
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD artisan registry user and password")
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
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
	)
	// create an instance of the runner
	run, err := runner.New()
	core.CheckErr(err, "cannot initialise runner")
	// load environment variables from file
	// NOTE: do not load from host environment to prevent clashes in the container
	env, err := merge.NewEnVarFromFile(c.envFilename)
	core.CheckErr(err, "failed to load environment file '%s'", c.envFilename)
	if len(c.credentials) == 0 {
		core.Msg("no credentials have been provided, if you are connecting to a authenticated registry, you need to pass the -u flag")
	}
	// launch a runtime to execute the function
	err = run.ExeC(packageName, fxName, c.credentials, *c.interactive, env)
	i18n.Err(err, i18n.ERR_CANT_EXEC_FUNC_IN_PACKAGE, fxName, packageName)
}
