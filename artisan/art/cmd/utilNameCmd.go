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
)

// UtilNameCmd generates passwords
type UtilNameCmd struct {
	Cmd          *cobra.Command
	number       *int
	specialChars *bool
}

func NewUtilNameCmd() *UtilNameCmd {
	c := &UtilNameCmd{
		Cmd: &cobra.Command{
			Use:   "name [flags]",
			Short: "generates a random name",
			Long:  `generates a random name`,
		},
	}
	c.number = c.Cmd.Flags().IntP("max-number", "n", 0, "adds a random number at the end of the name ranging from 0 to max-number")
	c.Cmd.Run = c.Run
	return c
}

func (c *UtilNameCmd) Run(_ *cobra.Command, _ []string) {
	fmt.Printf("%s", core.RandomName(*c.number))
}
