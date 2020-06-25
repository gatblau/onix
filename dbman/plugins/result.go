//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Parameter struct {
	value map[string]interface{}
	log   bytes.Buffer
}

func NewParameter() *Parameter {
	return &Parameter{
		value: make(map[string]interface{}),
		log:   bytes.Buffer{},
	}
}

func NewParameterFromJSON(jsonString string) *Parameter {
	r := &Parameter{
		value: make(map[string]interface{}),
		log:   bytes.Buffer{},
	}
	r.FromJSON(jsonString)
	return r
}

func (r *Parameter) GetString(key string) string {
	return r.value[key].(string)
}

func (r *Parameter) Get(key string) interface{} {
	return r.value[key]
}

func (r *Parameter) Set(key string, value interface{}) {
	r.value[key] = value
}

func (r *Parameter) SetTable(table Table) {
	r.value["table"] = table
}

func (r *Parameter) GetTable() *Table {
	if r.value["table"] != nil {
		if m, ok := r.value["table"].(map[string]interface{}); ok {
			// new table
			t := &Table{}
			// marshal the map to json
			bytes, _ := json.Marshal(m)
			// unmarshal the json to Table
			json.Unmarshal(bytes, &t)
			// return
			return t
		}
	}
	return nil
}

func (r *Parameter) SetErrorFromMessage(message string) {
	r.value["error"] = message
}

func (r *Parameter) SetError(err error) {
	r.value["error"] = err.Error()
}

func (r *Parameter) Log(message string) {
	r.log.WriteString(fmt.Sprintf("%s\n", message))
}

func (r *Parameter) ToString() string {
	r.value["log"] = r.log.String()
	bytes, err := json.Marshal(r.value)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (r *Parameter) ToError(err error) string {
	r.SetError(err)
	r.value["log"] = r.log.String()
	bytes, err := json.Marshal(r.value)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func (r *Parameter) FromJSON(jsonString string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &m)
	r.value = m
	return err
}

func (r *Parameter) HasError() bool {
	return r.value["error"] != nil
}

func (r *Parameter) Error() error {
	errMsg := r.value["error"]
	if errMsg != nil {
		return errors.New(r.value["error"].(string))
	}
	return nil
}

func (r *Parameter) PrintLog() {
	if r.value["log"] != nil {
		fmt.Print(r.value["log"])
	}
}

// print the content of the result to a string
func (r *Parameter) Sprint(format string) string {
	switch strings.ToLower(format) {
	case "yml":
		fallthrough
	case "yaml":
		o, err := yaml.Marshal(r.value)
		if err != nil {
			fmt.Printf("!!! cannot convert result to yaml: %v", err)
		}
		return string(o)
	case "json":
		o, err := json.MarshalIndent(r.value, "", " ")
		if err != nil {
			fmt.Printf("!!! cannot convert result to json: %v", err)
		}
		return string(o)
	default:
		fmt.Printf("!!! output format %v not supported, try yaml or json", format)
	}
	return ""
}

// save the content of the result to a file
func (r *Parameter) Save(format string, filename string) error {
	// get the path of the current executing process
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	// create a file with the getReleaseInfo getPlan
	f, err := os.Create(fmt.Sprintf("%v/%v.%v", exPath, filename, format))
	if err != nil {
		fmt.Printf("failed to create file: %s\n", err)
	}
	_, err = f.WriteString(r.Sprint(format))
	if err != nil {
		return err
	}
	return f.Close()
}
