/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
)

// SpecInfoCmd show spec information
type SpecInfoCmd struct {
	Cmd   *cobra.Command
	creds string
}

func NewSpecInfoCmd() *SpecInfoCmd {
	c := &SpecInfoCmd{
		Cmd: &cobra.Command{
			Use:   "info [OPTIONS] SPEC-FILE-PATH",
			Short: "displays release information defined in a spec file",
			Long: `displays release information defined in a spec file
Usage: art spec info [OPTIONS] SPEC-FILE-PATH

If the path to the spec.yaml file is not specified, the current folder is assumed.

Example:

given a spec.yaml as follows:

spec.yaml
---
version: 1.0
info: |
	This is the  information for Release 1.0
    More info here...
...
   
# running the command below:
art spec info .

# should display the following:
version: 1.0
This is the  information for Release 1.0
More info here...

`,
		},
	}
	c.Cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to retrieve the specification from an endpoint")
	c.Cmd.Run = c.Run
	return c
}

func (c *SpecInfoCmd) Run(cmd *cobra.Command, args []string) {
	// if no path to the spec.yaml has been provided
	if args == nil || len(args) == 0 {
		// assume current path
		args = []string{"."}
	}
	spec, err := release.NewSpec(args[0], c.creds)
	core.CheckErr(err, "cannot load spec")
	if len(spec.Version) > 0 {
		fmt.Printf("version: %s\n", spec.Version)
	}
	if len(spec.Info) > 0 {
		fmt.Println(spec.Info)
	} else {
		fmt.Println("no information available")
	}
}
