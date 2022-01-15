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

// SaveCmd save artisan packages for transfer between registries in disconnected scenarios
type SaveCmd struct {
	cmd *cobra.Command
}

func NewSaveCmd() *SaveCmd {
	c := &SaveCmd{
		cmd: &cobra.Command{
			Use:   "save",
			Short: "save artisan packages ready for transfer between registries in disconnected scenarios",
			Long:  ``,
		},
	}
	return c
}
