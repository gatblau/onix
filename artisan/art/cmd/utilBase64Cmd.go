/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
)

// UtilBase64Cmd generates passwords
type UtilBase64Cmd struct {
	cmd    *cobra.Command
	decode *bool
}

func NewUtilBase64Cmd() *UtilBase64Cmd {
	c := &UtilBase64Cmd{
		cmd: &cobra.Command{
			Use:   "b64 [flags] STRING",
			Short: "base 64 encode (or alternatively decode) a string",
			Long:  `base 64 encode (or alternatively decode) a string`,
			Args:  cobra.ExactArgs(1),
		},
	}
	c.cmd.Run = c.Run
	c.decode = c.cmd.Flags().BoolP("decode", "d", false, "if sets, decodes the string instead of encoding it")
	return c
}

func (c *UtilBase64Cmd) Run(_ *cobra.Command, args []string) {
	if *c.decode {
		decoded, err := base64.StdEncoding.DecodeString(args[0])
		core.CheckErr(err, "cannot decode string")
		fmt.Printf("%s", string(decoded[:]))
	} else {
		fmt.Printf("%s", base64.StdEncoding.EncodeToString([]byte(args[0])))
	}
}
