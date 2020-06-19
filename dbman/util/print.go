//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// toString a passed-in object to the console output or to a file
// filename: if set, output to a file
// obj: the object to toString
// format: json or yaml/yml
func Print(obj interface{}, format string, filename string) {
	// if an output filename is provided
	if len(filename) > 0 {
		// get the path of the current executing process
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		// create a file with the getReleaseInfo getPlan
		f, err := os.Create(fmt.Sprintf("%v/%v.%v", exPath, filename, format))
		if err != nil {
			fmt.Printf("failed to create getPlan file: %s\n", err)
		}
		f.WriteString(toString(obj, format))
		f.Close()
	} else {
		// toString the getPlan
		fmt.Println(toString(obj, format))
	}
}

// returns a string representation of the passed-in object in the requested format (i.e. json or yaml)
// obj: the object to be formatted
// format: the format to use for the output, either json or yaml/yml
func toString(obj interface{}, format string) string {
	table, isTable := obj.(Table)

	switch strings.ToLower(format) {
	case "yml":
		fallthrough
	case "yaml":
		o, err := yaml.Marshal(obj)
		if err != nil {
			fmt.Printf("!!! cannot convert output to yaml: %v", err)
		}
		return string(o)
	case "json":
		o, err := json.MarshalIndent(obj, "", " ")
		if err != nil {
			fmt.Printf("!!! cannot convert output to json: %v", err)
		}
		return string(o)
	case "csv":
		if isTable {
			o, err := tableTo(table, format)
			if err != nil {
				fmt.Printf("!!! cannot convert table to CVS: %v", err)
			}
			return o
		} else {
			fmt.Printf("!!! output format %v is only supported with queries", format)
		}
	default:
		fmt.Printf("!!! output format %v not supported, try yaml or json", format)
	}
	return ""
}

func tableTo(table Table, format string) (string, error) {
	switch strings.ToLower(format) {
	case "csv":
		{
			buffer := bytes.Buffer{}
			for i := 0; i < len(table.Header); i++ {
				buffer.WriteString(table.Header[i])
				if i < len(table.Header)-1 {
					buffer.WriteString(",")
				}
			}
			buffer.WriteString("\n")
			for _, row := range table.Rows {
				for i := 0; i < len(row); i++ {
					buffer.WriteString(row[i])
					if i < len(row)-1 {
						buffer.WriteString(",")
					}
				}
				buffer.WriteString("\n")
			}
			out := buffer.String()
			return out[:len(out)-1], nil
		}
	default:
		return "", errors.New(fmt.Sprintf("!!! I cannot recognise the format '%s'", format))
	}
}
