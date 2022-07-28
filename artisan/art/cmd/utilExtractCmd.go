/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"bufio"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// UtilExtractCmd client url issues http requests within a retry framework
type UtilExtractCmd struct {
	Cmd     *cobra.Command
	matches int
	prefix  string
	suffix  string
}

func NewUtilExtractCmd() *UtilExtractCmd {
	c := &UtilExtractCmd{
		Cmd: &cobra.Command{
			Use:     "extract [flags]",
			Short:   "extracts text between specified prefix and suffix, it should be used only with pipes",
			Long:    `extracts text between specified prefix and suffix, it should be used only with pipes`,
			Example: "cat your-file.txt | extract --prefix AAA --suffix $ -n 1",
			Args:    cobra.ExactArgs(0),
		},
	}
	c.Cmd.Flags().StringVarP(&c.prefix, "prefix", "p", "", "-p \"the prefix\"")
	c.Cmd.Flags().StringVarP(&c.suffix, "suffix", "s", "$", "-s \"the suffix\", if not specified an end of line marker is assumed")
	c.Cmd.Flags().IntVarP(&c.matches, "matches", "n", -1, "the maximum number of matches to retrieve")
	c.Cmd.Run = c.Run
	return c
}

func (c *UtilExtractCmd) Run(cmd *cobra.Command, args []string) {
	// captures information from the standard input
	info, _ := os.Stdin.Stat()
	// check that the standard input is not a character device file - i.e. one with which the Driver communicates
	// by sending and receiving single characters (bytes, octets)
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage:")
		fmt.Println("  cat your-file.txt | extract --prefix AAA --suffix $")
	} else if info.Size() > 0 {
		input := new(strings.Builder)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
		}
		output := core.Extract(input.String(), c.prefix, c.suffix, c.matches)
		if len(output) == 1 {
			fmt.Printf(output[0])
		} else {
			fmt.Printf(strings.Join(output, ","))
		}
	}
}
