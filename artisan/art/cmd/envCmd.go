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

// list local packages
type EnvCmd struct {
	cmd *cobra.Command
}

func NewEnvCmd() *EnvCmd {
	c := &EnvCmd{
		cmd: &cobra.Command{
			Use:   "env",
			Short: "extract environment information from packages and flows",
			Long:  `extract environment information from packages and flows`,
		},
	}
	return c
}
