/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=5%20Line%20Oblique&text=pilot%0A
type RootCmd struct {
	*cobra.Command
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		&cobra.Command{
			Use:   "pilot",
			Short: "Onix configuration manager agent",
			Long: `
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++
|                ___     ( ) //  ___    __  ___         |
|              //   ) ) / / // //   ) )  / /            |
|             //___/ / / / // //   / /  / /             |
|            //       / / // ((___/ /  / /              |
|         the Onix Pilot command line interface         |
|   Onix configuration agent for hosts and containers   |
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++`,
		},
	}
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootCmd) initConfig() {
}
