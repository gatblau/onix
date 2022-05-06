/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	_ "embed"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Cmd *cobra.Command
}

// NewRootCmd
// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Broadway%20KB&text=artisan%0A
func NewRootCmd() *RootCmd {
	c := &RootCmd{
		Cmd: &cobra.Command{
			Use:   "art",
			Short: "Artisan: the Onix DevOps CLI",
			Long: fmt.Sprintf(`
+++++++++++++++++| ONIX CONFIG MANAGER |+++++++++++++++++
|         __    ___  _____  _   __    __    _           |
|        / /\  | |_)  | |  | | ( ('  / /\  | |\ |       |
|       /_/--\ |_| \  |_|  |_| _)_) /_/--\ |_| \|       |
|           the DevOps command line interface           |
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++

version: %s`, core.Version),
			Version: core.Version,
		},
	}
	c.Cmd.SetVersionTemplate("version: {{.Version}}\n")
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootCmd) initConfig() {
}
