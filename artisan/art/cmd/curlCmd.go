/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
	"time"
)

// CurlCmd client url issues http requests within a retry framework
type CurlCmd struct {
	cmd           *cobra.Command
	maxAttempts   int
	creds         string
	method        string
	payload       string
	file          string
	validCodes    []int
	addValidCodes []int
	delaySecs     time.Duration
	timeout       time.Duration
}

func NewCurlCmd() *CurlCmd {
	c := &CurlCmd{
		cmd: &cobra.Command{
			Use:   "curl [flags] URI",
			Short: "issues an HTTP request and retry if a failure occurs",
			Long:  `issues an HTTP request and retry if a failure occurs`,
			Args:  cobra.ExactArgs(1),
		},
	}
	c.cmd.Flags().StringVarP(&c.creds, "creds", "u", "", "-u user:password")
	c.cmd.Flags().IntVarP(&c.maxAttempts, "max-attempts", "a", 5, "number of attempts before it stops retrying")
	c.cmd.Flags().IntSliceVarP(&c.validCodes, "valid-codes", "c", []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}, "comma separated list of HTTP status codes considered valid (e.g. no retry will be triggered)")
	c.cmd.Flags().IntSliceVarP(&c.validCodes, "add-valid-codes", "C", []int{}, "comma separated list of additional HTTP status codes considered valid (e.g. no retry will be triggered)")
	c.cmd.Flags().StringVarP(&c.method, "method", "X", "GET", "the http method to use (i.e. POST, PUT, GET, DELETE)")
	c.cmd.Flags().StringVarP(&c.payload, "payload", "P", "", "a string with the payload to be sent in the body of the http request")
	c.cmd.Flags().StringVarP(&c.file, "file", "f", "", "the location of a file which content is to be sent in the body of the http request")
	c.cmd.Flags().DurationVarP(&c.delaySecs, "delay", "d", 15*time.Second, "the delay in seconds between each retry interval")
	c.cmd.Flags().DurationVarP(&c.timeout, "timeout", "t", 30*time.Second, "the period after which the http request will timeout if not response is received from the server")
	c.cmd.Run = c.Run
	return c
}

func (c *CurlCmd) Run(cmd *cobra.Command, args []string) {
	uri := args[0]
	token := ""
	if len(c.creds) > 0 {
		uname, pwd := core.UserPwd(c.creds)
		token = core.BasicToken(uname, pwd)
	}
	// Curl(uri, method, token, validCodes, payload, file, maxAttempts, delaySecs)
	core.Curl(uri, c.method, token, append(c.validCodes, c.addValidCodes...), c.payload, c.file, c.maxAttempts, c.delaySecs, c.timeout)
}
