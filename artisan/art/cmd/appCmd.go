/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/app"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

// AppCmd creates deployment config files for a target platform
type AppCmd struct {
	cmd     *cobra.Command
	creds   string
	path    string
	format  string
	profile string
}

func NewAppCmd() *AppCmd {
	c := &AppCmd{
		cmd: &cobra.Command{
			Use:   "app [flags] [app manifest URI]",
			Short: "generates application deployment configuration files for a target platform",
			Long: `generates application deployment configuration files for a target platform\n
the URI can be either http(s):// or file://`,
			Args: cobra.ExactArgs(1),
		},
	}
	c.cmd.Flags().StringVarP(&c.creds, "creds", "u", "", "-u user:password; the credentials for git authentication")
	c.cmd.Flags().StringVarP(&c.format, "format", "f", "compose", "the target format for the configuration files (i.e. compose or k8s)")
	c.cmd.Flags().StringVarP(&c.path, "path", "p", ".", "the output path where the configuration files will be written")
	c.cmd.Flags().StringVarP(&c.profile, "profile", "r", "", "the application profile to use, "+
		"if not specified, and profiles have been defined in the application manifest, the first profile is used\n"+
		"if not specified, and profiles have not been defined in the application manifest, all services in the manifest are included ")
	c.cmd.Run = c.Run
	return c
}

func (c *AppCmd) Run(cmd *cobra.Command, args []string) {
	// get the app manifest URI
	uri := args[0]
	if len(uri) == 0 {
		core.ErrorLogger.Fatalf("missing application manifest URI\n")
	}
	if err := app.GenerateResources(uri, c.format, c.profile, c.creds, c.path); err != nil {
		core.ErrorLogger.Fatalf(err.Error())
	}
}
