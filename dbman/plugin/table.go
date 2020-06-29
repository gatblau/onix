//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// generic table used as a serializable result set for queries
type Table struct {
	Header Row   `json:"header,omitempty"`
	Rows   []Row `json:"row,omitempty"`
}

// a row in the table
type Row []string

// save the table to a file with the specified format
//   - filename: the filename with no extension
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Save(format string, filename string) {
	// get the path of the current executing process
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("!!! I cannot find the path to the current process: %s\n", err)
	}
	exPath := filepath.Dir(ex)
	// create a file with the getReleaseInfo getPlan
	f, err := os.Create(fmt.Sprintf("%v/%v.%v", exPath, filename, format))
	if err != nil {
		fmt.Printf("!!! I cannot create the result file: %s\n", err)
	}
	f.WriteString(table.Sprint(format))
	f.Close()
}

// return the table as a string of the specified format
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Sprint(format string) string {
	switch strings.ToLower(format) {
	case "yml":
		fallthrough
	case "yaml":
		return table.AsYAML()
	case "json":
		return table.AsJSON()
	case "csv":
		return table.AsCSV()
	default:
		fmt.Printf("!!! output format %v not supported, try YAML, JSON or CSV", format)
	}
	return ""
}

// return the table as a JSON string
func (table *Table) AsJSON() string {
	o, err := json.MarshalIndent(table, "", " ")
	if err != nil {
		fmt.Printf("!!! cannot convert output to JSON: %v", err)
	}
	return string(o)
}

// return the table as a YAML string
func (table *Table) AsYAML() string {
	o, err := yaml.Marshal(table)
	if err != nil {
		fmt.Printf("!!! cannot convert output to YAML: %v", err)
	}
	return string(o)
}

// return the table as a CSV string
func (table *Table) AsCSV() string {
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
	return out[:len(out)-1]
}

// print the table content with the specified format to the stdout
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Print(format string) {
	fmt.Println(table.Sprint(format))
}
