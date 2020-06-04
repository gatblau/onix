//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// print a passed-in object to the console output or to a file
// filename: if set, output to a file
// obj: the object to print
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
		// print the getPlan
		fmt.Println(toString(obj, format))
	}
}

// returns a string representation of the passed-in object in the requested format (i.e. json or yaml)
// obj: the object to be formatted
// format: the format to use for the output, either json or yaml/yml
func toString(obj interface{}, format string) string {
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
	default:
		fmt.Printf("!!! output format %v not supported, try yaml or json", format)
	}
	return ""
}
