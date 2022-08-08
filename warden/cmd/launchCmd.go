/*
  Onix Config Manager - Warden
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"log"
	"strings"

	"github.com/gatblau/onix/warden/mode"
	"github.com/spf13/cobra"
)

// LaunchCmd launches host pilot
type LaunchCmd struct {
	cmd                *cobra.Command
	mode               string // set the proxy mode of operation
	verbose            bool   // enable verbose output
	address            string
	uri                string
	credential         string
	bearertoken        string
	insecureskipverify bool
}

func NewLaunchCmd() *LaunchCmd {
	c := &LaunchCmd{
		cmd: &cobra.Command{
			Use:   "launch",
			Short: "launches warden http proxy",
			Long:  `launches warden http proxy`,
		},
	}

	c.cmd.Flags().BoolVarP(&c.verbose, "verbose", "v", false, "enables verbose output")
	c.cmd.Flags().StringVarP(&c.mode, "mode", "m", "basic", "tell warden how to setup the proxy based on operation modes, allowed values are basic, tap")
	c.cmd.Flags().StringVarP(&c.address, "port", "a", ":8080", "port at which proxy should listen")
	c.cmd.Flags().StringVarP(&c.uri, "uri", "l", "", "uri to which request body has to be forwarded")
	c.cmd.Flags().StringVarP(&c.credential, "credential", "c", "", "credentials, username:password to connect to uri")
	c.cmd.Flags().StringVarP(&c.bearertoken, "bearertoken", "t", "", "token to connect to uri")
	c.cmd.Flags().BoolVarP(&c.insecureskipverify, "insecureskipverify", "i", false, "skip insecure ssl certificate verification when connecting to uri")

	c.cmd.Run = c.Run
	return c
}

func (c *LaunchCmd) Run(_ *cobra.Command, _ []string) {
	switch strings.ToUpper(c.mode) {
	case "BASIC":
		mode.Basic(c.address, c.verbose)
	case "TAP":
		mode.Tap(c.uri, c.credential, c.address, c.bearertoken, c.verbose, c.insecureskipverify)
	default:
		log.Fatalf("invalid mode: '%s'", c.mode)
	}
}
