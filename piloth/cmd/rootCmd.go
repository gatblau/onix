/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Cmd *cobra.Command
}

// NewRootCmd creates the root command
func NewRootCmd() *RootCmd {
	c := &RootCmd{
		Cmd: &cobra.Command{
			Use:   "piloth",
			Short: "Host Pilot",
			Long: fmt.Sprintf(`
+-----------------| ONIX CONFIG MANAGER |-----------------+
|      _/_/_/    _/_/_/  _/          _/_/    _/_/_/_/_/   |
|     _/    _/    _/    _/        _/    _/      _/        |
|    _/_/_/      _/    _/        _/    _/      _/         |
|   _/          _/    _/        _/    _/      _/          |
|  _/        _/_/_/  _/_/_/_/    _/_/        _/           |
|                     Host Controller                     | 
+---------------------------------------------------------+

version: %s`, core.Version),
			Version: core.Version,
		},
	}
	c.Cmd.SetVersionTemplate("Pilot Host version: {{.Version}}\n")
	return c
}
