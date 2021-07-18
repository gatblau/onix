package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TemplMerger merge artisan templates using artisan inputs
type TemplMerger struct {
	regex    *regexp.Regexp
	template map[string][]byte
	file     map[string][]byte
}

// NewTemplMerger create a new instance of the template merger to merge files
func NewTemplMerger() (*TemplMerger, error) {
	regex, err := regexp.Compile("\\${(?P<NAME>[^}]*)}")
	if err != nil {
		return nil, fmt.Errorf("cannot compile regex: %s\n", err)
	}
	return &TemplMerger{
		regex: regex,
	}, nil
}

// LoadTemplates load the template files to use
func (t *TemplMerger) LoadTemplates(files []string) error {
	m := make(map[string][]byte)
	for _, file := range files {
		// check the file is a template
		if !(filepath.Ext(file) == ".tem" || filepath.Ext(file) == ".t") {
			return fmt.Errorf("file '%s' is not a template file, artisan templates are either .tem or .t files\n", file)
		}
		// ensure the template path is absolute
		path, err := AbsPath(file)
		if err != nil {
			return fmt.Errorf("path '%s' cannot be converted to absolute path: %s\n", file, err)
		}
		// read the file content
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %s\n", file, err)
		}
		m[path] = bytes
	}
	t.template = m
	return nil
}

// Merge merge templates with the passed in environment
func (t *TemplMerger) Merge(env *Envar) error {
	t.file = make(map[string][]byte)
	for path, file := range t.template {
		var (
			merged []byte
			err    error
		)
		// if the template is in simple tem format
		if strings.HasSuffix(path, "tem") {
			merged, err = t.mergeTem(file, *env)
			if err != nil {
				return fmt.Errorf("cannot merge template '%s': %s\n", path, err)
			}
			t.file[path[0:len(path)-len(".tem")]] = merged
		} else {
			merged, err = t.mergeT(file, *env)
			if err != nil {
				return fmt.Errorf("cannot merge template '%s': %s\n", path, err)
			}
			t.file[path[0:len(path)-len(".t")]] = merged
		}
	}
	return nil
}

// mergeTem merges a single template file using tem format and the passed in variables
func (t *TemplMerger) mergeTem(tem []byte, env Envar) ([]byte, error) {
	content := string(tem)
	// find all environment variable placeholders in the content
	vars := t.regex.FindAll(tem, -1)
	// loop though the found vars to merge
	for _, v := range vars {
		defValue := ""
		// removes placeholder marks: ${...}
		vname := strings.TrimSuffix(strings.TrimPrefix(string(v), "${"), "}")
		// is a default value defined?
		cut := strings.Index(vname, ":")
		// split default value and var name
		if cut > 0 {
			// get the default value
			defValue = vname[cut+1:]
			// get the name of the var without the default value
			vname = vname[0:cut]
		}
		// check the name of the env variable is not "PWD" as it can return the current directory in some OSs
		if vname == "PWD" {
			fmt.Errorf("environment variable cannot be PWD, choose a different name\n")
		}
		// fetch the env variable value
		ev := env.Vars[vname]
		// if the variable is not defined in the environment
		if len(ev) == 0 {
			// if no default value has been defined
			if len(defValue) == 0 {
				return nil, fmt.Errorf("environment variable '%s' required and not defined, cannot merge\n", vname)
			} else {
				// merge with the default value
				content = strings.Replace(content, string(v), defValue, -1)
			}
		} else {
			// merge with the env variable value
			content = strings.Replace(content, string(v), ev, -1)
		}
	}
	return []byte(content), nil
}

// mergeT merges a single template file using go template format and the passed in variables
func (t *TemplMerger) mergeT(tem []byte, env Envar) ([]byte, error) {
	tt, err := template.New("t").Funcs(template.FuncMap{
		"group": env.Group,
	}).Parse(string(tem))
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	err = tt.Execute(&tpl, env)
	if err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}
