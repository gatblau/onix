/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

type TknCmd struct {
	Cmd *cobra.Command
}

func NewTknCmd() *TknCmd {
	c := &TknCmd{
		Cmd: &cobra.Command{
			Use:   "tkn",
			Short: "provides functions to manage Tekton pipelines",
			Long:  `provides functions to manage Tekton pipelines`,
		},
	}
	return c
}
