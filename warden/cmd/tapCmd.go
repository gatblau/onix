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
type TapCmd struct {
	cmd                *cobra.Command
	verbose            bool   // enable verbose output
	address            string // port number
	uri                string
	credential         string
	bearertoken        string
	insecureskipverify bool
}

func NewTapCmd() *TapCmd {
	c := &TapCmd{
		cmd: &cobra.Command{
			Use:   "tap [flags]",
			Short: "launches warden http proxy in tap mode",
			Long:  `launches warden http proxy with tap configuration`,
		},
	}
	c.cmd.Flags().BoolVarP(&c.verbose, "verbose", "v", false, "enables verbose output")
	c.cmd.Flags().StringVarP(&c.address, "port", "p", ":8080", "port at which proxy should listen, expected format is :port_number")
	c.cmd.Flags().StringVarP(&c.uri, "uri", "l", "", "uri to which request body has to be forwarded")
	c.cmd.Flags().StringVarP(&c.credential, "credential", "c", "", "credentials, username:password to connect to uri")
	c.cmd.Flags().StringVarP(&c.bearertoken, "bearertoken", "t", "", "token to connect to uri")
	c.cmd.Flags().BoolVarP(&c.insecureskipverify, "insecureskipverify", "i", false, "skip insecure ssl certificate verification when connecting to uri")
	c.cmd.Run = c.Run
	return c
}

func (c *TapCmd) Run(_ *cobra.Command, _ []string) {
	mode.Tap(c.uri, c.credential, c.address, c.bearertoken, c.verbose, c.insecureskipverify)
}
