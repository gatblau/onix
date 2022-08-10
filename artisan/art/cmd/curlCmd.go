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

// CurlCmd client url issues http requests within a retry framework
type CurlCmd struct {
	Cmd           *cobra.Command
	maxAttempts   int
	creds         string
	method        string
	payload       string
	file          string
	validCodes    []int
	addValidCodes []int
	delaySecs     int
	timeoutSecs   int
	headers       []string
	outFile       string
	response      bool
}

func NewCurlCmd() *CurlCmd {
	c := &CurlCmd{
		Cmd: &cobra.Command{
			Use:   "curl [flags] URI",
			Short: "issues an HTTP request and retry if a failure occurs",
			Long:  `issues an HTTP request and retry if a failure occurs`,
			Args:  cobra.ExactArgs(1),
		},
	}
	c.Cmd.Flags().StringVarP(&c.creds, "creds", "u", "", "-u user:password")
	c.Cmd.Flags().IntVarP(&c.maxAttempts, "max-attempts", "a", 5, "number of attempts before it stops retrying")
	c.Cmd.Flags().IntSliceVarP(&c.validCodes, "valid-codes", "c", []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}, "comma separated list of HTTP status codes considered valid (e.g. no retry will be triggered)")
	c.Cmd.Flags().IntSliceVarP(&c.addValidCodes, "add-valid-codes", "C", []int{}, "comma separated list of additional HTTP status codes considered valid (e.g. no retry will be triggered)")
	c.Cmd.Flags().StringVarP(&c.method, "method", "X", "GET", "the http method to use (i.e. POST, PUT, GET, DELETE)")
	c.Cmd.Flags().StringVarP(&c.outFile, "out-file", "o", "", "the name of the file where the http response body should be saved; if not set, the response is not saved but printed to stdout")
	c.Cmd.Flags().StringVarP(&c.payload, "payload", "d", "", "a string with the payload to be sent in the body of the http request")
	c.Cmd.Flags().StringVarP(&c.file, "file", "f", "", "the location of a file which content is to be sent in the body of the http request")
	c.Cmd.Flags().StringSliceVarP(&c.headers, "headers", "H", nil, "a comma separated list of http headers (format 'key1:value1','key2:value2,...,'keyN:valueN')")
	c.Cmd.Flags().IntVarP(&c.delaySecs, "delay", "r", 5, "the retry delay (in seconds)")
	c.Cmd.Flags().IntVarP(&c.timeoutSecs, "timeout", "t", 30, "the period (in seconds) after which the http request will timeout if not response is received from the server")
	c.Cmd.Flags().BoolVarP(&c.response, "response", "v", false, "if set, shows additional response information such as status code and headers")
	c.Cmd.Run = c.Run
	return c
}

func (c *CurlCmd) Run(cmd *cobra.Command, args []string) {
	uri := args[0]
	token := ""
	if len(c.creds) > 0 {
		uname, pwd := core.UserPwd(c.creds)
		token = httpserver.BasicToken(uname, pwd)
	}
	core.Curl(uri, c.method, token, append(c.validCodes, c.addValidCodes...), c.payload, c.file, c.maxAttempts, c.delaySecs, c.timeoutSecs, c.headers, c.outFile, c.response)
}
