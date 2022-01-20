/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/spf13/cobra"
)

// WaitCmd wait until the response payload contains a specific value
type WaitCmd struct {
	cmd      *cobra.Command
	attempts int
	filter   string
	creds    string
}

func NewWaitCmd() *WaitCmd {
	c := &WaitCmd{
		cmd: &cobra.Command{
			Use:   "wait [flags] URI",
			Short: "wait until either the an HTTP GET returns a value or the maximum attempts have been reached",
			Long:  `wait until either the an HTTP GET returns a value or the maximum attempts have been reached`,
			Args:  cobra.ExactArgs(1),
		},
	}
	c.cmd.Flags().StringVarP(&c.filter, "filter", "f", "", "-f json/path/expression")
	c.cmd.Flags().StringVarP(&c.creds, "creds", "u", "", "-u user:password")
	c.cmd.Flags().IntVarP(&c.attempts, "attempts", "a", 5, "-a 10 (number of attempts before it fails)")
	c.cmd.Run = c.Run
	return c
}

func (c *WaitCmd) Run(cmd *cobra.Command, args []string) {
	token := ""
	if len(c.creds) > 0 {
		uname, pwd := core.UserPwd(c.creds)
		token = httpserver.BasicToken(uname, pwd)
	}
	core.Wait(args[0], c.filter, token, c.attempts)
}
