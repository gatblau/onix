//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/dbman/core"
	"github.com/gatblau/onix/dbman/plugin"
	"github.com/spf13/cobra"
	"strings"
)

type DbQueryCmd struct {
	cmd      *cobra.Command
	format   string
	filename string
}

func NewDbQueryCmd() *DbQueryCmd {
	c := &DbQueryCmd{
		cmd: &cobra.Command{
			Use:     "query [name] [args...]",
			Short:   "runs a database query. args if nay, should be in the format key1=value1,key2=value2,...,keyN=valueN",
			Long:    ``,
			Example: "dbman db query db-version '<APP_VERSION>=0.0.4'",
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.format, "output", "o", "json", "the format of the output - yaml, json, csv")
	c.cmd.Flags().StringVarP(&c.filename, "filename", "f", "", `if a filename is specified, the output will be written to the file. The file name should not include extension.`)
	return c
}

func (c *DbQueryCmd) Run(cmd *cobra.Command, args []string) {
	var (
		queryName string
		params    = make(map[string]string)
	)
	// check the query name has been passed in
	if len(args) == 0 {
		fmt.Printf("!!! You forgot to tell me the name of the query you want to run\n")
		return
	}
	// get the query parameters
	if len(args) < 3 {
		queryName = args[0]
		// check if we have parameters
		if len(args) == 2 {
			parts := strings.Split(args[1], ",")
			for _, part := range parts {
				subPart := strings.Split(part, "=")
				if len(subPart) != 2 {
					fmt.Printf("!!! I cannot break down query parameter '%s': format should be 'key=value'\n", subPart)
					return
				}
				params[strings.Trim(subPart[0], " ")] = strings.Trim(subPart[1], " ")
			}
		}
	} else {
		fmt.Printf("!!! Too many parameters\n")
		return
	}
	// execute the query
	result, _, err := core.DM.Query(queryName, params)
	if err != nil {
		fmt.Printf("!!! I cannot run query '%s': %s\n", queryName, err)
		return
	}
	// if a filename has been specified
	if len(c.filename) > 0 {
		// save to disk
		result.Save(c.format, c.filename)
	} else {
		// print to stdout
		result.Print(c.format)
	}
}

func varsToString(vars []plugin.Var) string {
	buffer := bytes.Buffer{}
	for i, v := range vars {
		buffer.WriteString(v.Name)
		if i < len(vars)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}
