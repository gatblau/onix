//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"fmt"
	. "github.com/gatblau/onix/dbman/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// decorator for the release plan cobra command
type ReleasePlanCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewReleasePlanCmd() *ReleasePlanCmd {
	c := &ReleasePlanCmd{
		cmd: &cobra.Command{
			Use:   "plan",
			Short: "displays the release plan",
			Long:  `A release plan is the list of all releases available`,
		},
	}
	c.cmd.Run = c.run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml or json")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", "if the plan filename (without extension) is specified, the output will be written to the file")
	return c
}

func (c *ReleasePlanCmd) run(cmd *cobra.Command, args []string) {
	// fetch the release plan
	plan, err := DM.GetReleasePlan()
	if err != nil {
		fmt.Printf("oops! cannot get release plan: %v", err)
		return
	}
	// if an output filename is provided
	if len(c.filename) > 0 {
		// get the path of the current executing process
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		// create a file with the release plan
		f, err := os.Create(fmt.Sprintf("%v/%v.%v", exPath, c.filename, c.format))
		if err != nil {
			fmt.Printf("failed to create plan file: %s\n", err)
		}
		f.WriteString(plan.Format(c.format))
		f.Close()
	} else {
		// print the plan
		fmt.Println(plan.Format(c.format))
	}
}
