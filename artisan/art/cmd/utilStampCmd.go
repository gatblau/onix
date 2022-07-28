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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UtilStampCmd generates a timestamp
type UtilStampCmd struct {
	Cmd  *cobra.Command
	path string
}

func NewUtilStampCmd() *UtilStampCmd {
	c := &UtilStampCmd{
		Cmd: &cobra.Command{
			Use:   "stamp [flags]",
			Short: "prints the current timestamp in UTC Unix Nano format",
			Long:  `prints the current timestamp in UTC Unix Nano format`,
			Args:  cobra.ExactArgs(0),
		},
	}
	c.Cmd.Run = c.Run
	c.Cmd.Flags().StringVarP(&c.path, "file-path", "p", "", "if set, writes the timestamp to the file system path")
	return c
}

func (c *UtilStampCmd) Run(_ *cobra.Command, args []string) {
	if len(c.path) > 0 {
		path, _ := filepath.Abs(c.path)
		os.WriteFile(path, []byte(strconv.FormatInt(time.Now().UTC().UnixNano(), 10)), 0755)
	} else {
		fmt.Printf("%d", time.Now().UTC().UnixNano())
	}
}
