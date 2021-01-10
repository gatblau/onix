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
)

// executes an exported function
type ExecCmd struct {
	cmd             *cobra.Command
	interactive     *bool
	container       *bool
	credentials     string
	noTLS           *bool
	ignoreSignature *bool
	path            string
	pubPath         string
}

func NewExecCmd() *ExecCmd {
	c := &ExecCmd{
		cmd: &cobra.Command{
			Use:   "exec [package name] [function]",
			Short: "execute an exported function in a package",
			Long:  ``,
		},
	}
	c.interactive = c.cmd.Flags().BoolP("interactive", "i", false, "switches on interactive mode which prompts the user for information if not provided")
	c.container = c.cmd.Flags().BoolP("container", "c", false, "runs the function within a container of the image specified in the manifest runtime property")
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	c.noTLS = c.cmd.Flags().BoolP("no-tls", "t", false, "use -t or --no-tls to connect to a artisan registry over plain HTTP")
	c.ignoreSignature = c.cmd.Flags().BoolP("ignore-sig", "s", false, "-s or --ignore-sig to ignore signature verification")
	c.cmd.Flags().StringVarP(&c.pubPath, "pub", "p", "", "--pub=/path/to/public/key or -p=/path/to/public/key")
	c.cmd.Run = c.Run
	return c
}

func (r *ExecCmd) Run(cmd *cobra.Command, args []string) {
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
	core.CheckErr(err, "invalid package name")
	// run the function on the open package
	builder.Execute(name, function, r.credentials, *r.noTLS, r.pubPath, *r.ignoreSignature, *r.interactive)
}
