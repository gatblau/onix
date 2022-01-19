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
	"github.com/spf13/cobra"
	"net/http"
	"path/filepath"
)

// ServeCmd serve static files over http
type ServeCmd struct {
	cmd  *cobra.Command
	port int
}

func NewServeCmd() *ServeCmd {
	c := &ServeCmd{
		cmd: &cobra.Command{
			Use:   "serve [flags] PATH",
			Short: "serves static files over an http endpoint",
			Long:  `serves static files over an http endpoint`,
		},
	}
	c.cmd.Flags().IntVarP(&c.port, "port", "P", 8100, "the http port on which the server listens for connections")
	c.cmd.Run = c.Run
	return c
}

func (c *ServeCmd) Run(cmd *cobra.Command, args []string) {
	var path string
	if len(args) == 0 {
		path = "."
	} else {
		path = args[0]
	}
	path, err := filepath.Abs(path)
	core.CheckErr(err, "cannot resolve absolute path")
	// create file server handler
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)
	// start HTTP server with `http.DefaultServeMux` handler
	core.InfoLogger.Printf("serving the contents of '%s'", path)
	core.CheckErr(http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil), "cannot start http server")
}
