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

// UtilCmd miscellaneous utilities
type UtilCmd struct {
	Cmd *cobra.Command
}

func NewUtilCmd() *UtilCmd {
	c := &UtilCmd{
		Cmd: &cobra.Command{
			Use:   "u",
			Short: "provides access to various miscellaneous utilities",
			Long:  `provides access to various miscellaneous utilities`,
		},
	}
	return c
}
