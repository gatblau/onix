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
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

// list local artefacts
type PipeGenCmd struct {
	cmd         *cobra.Command
	envFilename string
}

func NewPipeGenCmd() *PipeGenCmd {
	c := &PipeGenCmd{
		cmd: &cobra.Command{
			Use:   "gen [flags] [/path/to/flow.yaml]",
			Short: "generates a pipeline definition",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.envFilename, "env", "e", ".env", "--env=.env or -e=.env")
	c.cmd.Run = c.Run
	return c
}

func (b *PipeGenCmd) Run(cmd *cobra.Command, args []string) {
	var flow string
	switch len(args) {
	case 0:
		flow = ""
	case 1:
		flow = args[0]
		if !path.IsAbs(flow) {
			abs, err := filepath.Abs(flow)
			core.CheckErr(err, "cannot convert '%s' to absolute path", flow)
			flow = abs
		}
	default:
		core.RaiseErr("too many arguments")
	}
	// try to load env from file
	core.LoadEnvFromFile(b.envFilename)

}
