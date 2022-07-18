/*
  Onix Config Manager - Warden
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/warden/mode"
	"github.com/spf13/cobra"
)

// LaunchCmd launches host pilot
type BasicCmd struct {
	cmd     *cobra.Command
	verbose bool // enable verbose output
	address string
}

func NewBasicCmd() *BasicCmd {
	c := &BasicCmd{
		cmd: &cobra.Command{
			Use:   "basic [flags]",
			Short: "launches warden http proxy in basic mode",
			Long:  `launches warden http proxy with basic configuration`,
		},
	}
	c.cmd.Flags().BoolVarP(&c.verbose, "verbose", "v", false, "enables verbose output")
	c.cmd.Flags().StringVarP(&c.address, "port", "a", ":8080", "port at which proxy should listen")
	c.cmd.Run = c.Run
	return c
}

func (c *BasicCmd) Run(_ *cobra.Command, _ []string) {
	mode.Basic(c.address, c.verbose)
}
