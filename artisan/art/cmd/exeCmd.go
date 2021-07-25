package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/spf13/cobra"
	"os"
)

// executes an exported function
type ExeCmd struct {
	cmd             *cobra.Command
	interactive     *bool
	credentials     string
	noTLS           *bool
	ignoreSignature *bool
	path            string
	pubPath         string
	envFilename     string
	preserveFiles   *bool
}

func NewExeCmd() *ExeCmd {
	c := &ExeCmd{
		cmd: &cobra.Command{
			Use:   "exe [package name] [function]",
			Short: "runs a function within a package on the current host",
			Long:  `runs a function within a package on the current host`,
		},
	}
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.noTLS = c.cmd.Flags().BoolP("no-tls", "t", false, "use -t or --no-tls to connect to a artisan registry over plain HTTP")
	c.ignoreSignature = c.cmd.Flags().BoolP("ignore-sig", "s", false, "-s or --ignore-sig to ignore signature verification")
	c.cmd.Flags().StringVarP(&c.pubPath, "pub", "p", "", "--pub=/path/to/public/key or -p=/path/to/public/key - public PGP key to verify package source")
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
	c.cmd.Flags().StringVar(&c.path, "path", "", "--path=/path/to/package/files - specify the location where the Artisan package must be open. If not specified, Artisan opens the package in a temporary folder under a randomly generated name.")
	c.preserveFiles = c.cmd.Flags().BoolP("preserve-files", "f", false, "use -f to preserve the open package files")
	c.cmd.Run = c.Run
	return c
}

func (c *ExeCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		core.RaiseErr("package and function names are required")
	}
	var (
		pack     = args[0]
		function = args[1]
	)
	// get a builder handle
	builder := build.NewBuilder()
	name, err := core.ParseName(pack)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// add the build file level environment variables
	env := merge.NewEnVarFromSlice(os.Environ())
	// load vars from file
	env2, err := merge.NewEnVarFromFile(c.envFilename)
	core.CheckErr(err, "failed to load environment file '%s'", c.envFilename)
	// merge with existing environment
	env.Merge(env2)
	// run the function on the open package
	builder.Execute(name, function, c.credentials, *c.noTLS, c.pubPath, *c.ignoreSignature, *c.interactive, c.path, *c.preserveFiles, env)
}
