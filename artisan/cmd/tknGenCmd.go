/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/tkn"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

type TknGenCmd struct {
	cmd    *cobra.Command
	pkPath string
}

func NewTknGenCmd() *TknGenCmd {
	c := &TknGenCmd{
		cmd: &cobra.Command{
			Use:   "gen [flags] [/path/to/flow.yaml]",
			Short: "generates a tekton pipeline",
			Long:  `generates a tekton pipeline`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.pkPath, "key", "k", "", "--key=/path/to/private/key or -k=/path/to/private/key")
	return c
}

func (b *TknGenCmd) Run(cmd *cobra.Command, args []string) {
	var flowPath string
	switch len(args) {
	case 0:
		flowPath = ""
	case 1:
		flowPath = args[0]
		if !path.IsAbs(flowPath) {
			abs, err := filepath.Abs(flowPath)
			core.CheckErr(err, "cannot convert '%s' to absolute path", flowPath)
			flowPath = abs
		}
	default:
		core.RaiseErr("too many arguments")
	}
	if len(b.pkPath) == 0 {
		core.RaiseErr("decryption key path is required")
	}
	f, err := flow.LoadFlow(flowPath)
	core.CheckErr(err, "cannot load flow")
	pk, err := crypto.LoadPGP(b.pkPath)
	core.CheckErr(err, "cannot load decryption key")
	builder := tkn.NewBuilder(f, pk)
	buf := builder.Create()
	fmt.Println(buf.String())
}
