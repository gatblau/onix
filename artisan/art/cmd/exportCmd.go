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
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// ExportCmd exports a package from the local registry to allow copying without using registries
type ExportCmd struct {
	cmd         *cobra.Command
	credentials string
	output      string
}

func NewExportCmd() *ExportCmd {
	c := &ExportCmd{
		cmd: &cobra.Command{
			Use:   "export PACKAGE_NAME[:TAG] [flags]",
			Short: "exports a package",
			Long:  `exports a package`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.output, "output", "o", ".", "-o ./exported; the output where the package will be exported")
	c.cmd.Flags().StringVarP(&c.credentials, "user", "u", "", "USER:PASSWORD server user and password")
	return c
}

func (c *ExportCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if len(args) < 1 {
		log.Fatal("name of the package to export is required")
	}
	// get the name of the package to push
	nameTag := args[0]
	// validate the name
	name, err := core.ParseName(nameTag)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry()
	// export packages into tar bytes
	content, err := local.Export(name, c.credentials)
	core.CheckErr(err, "failed to export package")
	if len(c.output) == 0 {
		fmt.Print(content)
	} else {
		absPath, err := filepath.Abs(c.output)
		core.CheckErr(err, "cannot obtain absolute output")
		core.CheckErr(os.MkdirAll(absPath, 0755), "cannot create target output")
		targetPath := filepath.Join(absPath, fmt.Sprintf("%s.tar", name.NormalString()))
		core.CheckErr(os.WriteFile(targetPath, content, 0755), "cannot write exported package file")
	}
}
