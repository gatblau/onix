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

// ImportCmd Import the contents from a tarball to create an artisan package in the local registry
type ImportCmd struct {
	cmd *cobra.Command
}

func NewImportCmd() *ImportCmd {
	c := &ImportCmd{
		cmd: &cobra.Command{
			Use:   "import",
			Short: "imports the content from one or more tarball files to create one or more packages in the local registry",
			Long: `Imports the content from one or more tarball files to create one or more packages in the local registry.
`,
		},
	}
	return c
}
