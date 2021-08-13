package merge

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
	"github.com/gatblau/onix/artisan/core"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TemplMerger merge artisan templates using artisan inputs
type TemplMerger struct {
	regex    *regexp.Regexp
	rexVar   *regexp.Regexp
	rexRange *regexp.Regexp
	rexItem  *regexp.Regexp
	template map[string][]byte
	file     map[string][]byte
}

// NewTemplMerger create a new instance of the template merger to merge files
func NewTemplMerger() (*TemplMerger, error) {
	// for tem templates:
	// parse ${NAME} vars
	regex, err := regexp.Compile("\\${(?P<NAME>[^}]*)}")
	if err != nil {
		return nil, fmt.Errorf("cannot compile regex: %s\n", err)
	}
	// for art templates:
	// parse {{ $ "NAME"  }} vars
	rexVar, err := regexp.Compile("{{[\\s]*\\$[\\s]*\"(?P<NAME>[\\w]+)\"[\\s]*}}")
	if err != nil {
		return nil, fmt.Errorf("cannot compile regex: %s\n", err)
	}
	// parse {{ range => "GROUP_NAME" }}
	rexRange, err := regexp.Compile("{{[\\s]*range[\\s]*=>[\\s]*\"(?P<GROUP>[\\w]+)\"[\\s]*}}")
	if err != nil {
		return nil, fmt.Errorf("cannot compile regex: %s\n", err)
	}
	// parse {{ % "NAME" }}
	rexItem, err := regexp.Compile("{{[\\s]*\\%[\\s]*\"(?P<ITEM>[\\w]+)\"[\\s]*}}")
	if err != nil {
		return nil, fmt.Errorf("cannot compile regex: %s\n", err)
	}
	return &TemplMerger{
		regex:    regex,
		rexVar:   rexVar,
		rexRange: rexRange,
		rexItem:  rexItem,
	}, nil
}

// LoadTemplates load the template files to use
func (t *TemplMerger) LoadTemplates(files []string) error {
	m := make(map[string][]byte)
	for _, file := range files {
		// check the file is a template
		if !(filepath.Ext(file) == ".tem" || filepath.Ext(file) == ".art") {
			return fmt.Errorf("file '%s' is not a template file, artisan templates are either .tem or .art files\n", file)
		}
		// ensure the template path is absolute
		path, err := core.AbsPath(file)
		if err != nil {
			return fmt.Errorf("path '%s' cannot be converted to absolute path: %s\n", file, err)
		}
		// read the file content
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %s\n", file, err)
		}
		m[path] = t.transpileOperators(bytes)
	}
	t.template = m
	return nil
}

// Merge templates with the passed in environment
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
			merged, err = t.mergeART(path, file, *env)
			if err != nil {
				return fmt.Errorf("cannot merge template: %s\n", err)
			}
			t.file[path[0:len(path)-len(".art")]] = merged
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

// mergeART merges a single template file using go template format and the passed in variables
func (t *TemplMerger) mergeART(path string, temp []byte, env Envar) ([]byte, error) {
	ctx, err := NewContext(env)
	if err != nil {
		return nil, err
	}
	tt, err := template.New(path).Funcs(template.FuncMap{
		"select": ctx.Select,
		"item":   ctx.Item,
		"var":    ctx.Var,
	}).Parse(string(temp))
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	err = tt.Execute(&tpl, ctx)
	if err != nil {
		return nil, err
	}
	return removeEmptyLines(tpl.String())
}

func removeEmptyLines(in string) ([]byte, error) {
	regex, err := regexp.Compile("\n\n")
	if err != nil {
		return nil, err
	}
	return []byte(regex.ReplaceAllString(in, "\n")), nil
}

func (t *TemplMerger) transpileOperators(source []byte) []byte {
	names := t.rexVar.FindAllStringSubmatch(string(source), -1)
	for _, n := range names {
		str := strings.ReplaceAll(string(source), n[0], fmt.Sprintf("{{ var \"%s\" }}", n[1]))
		source = []byte(str)
	}
	names = t.rexRange.FindAllStringSubmatch(string(source), -1)
	for _, n := range names {
		str := strings.ReplaceAll(string(source), n[0], fmt.Sprintf("{{ select \"%s\"}}{{ range .Items }}", n[1]))
		source = []byte(str)
	}
	names = t.rexItem.FindAllStringSubmatch(string(source), -1)
	for _, n := range names {
		str := strings.ReplaceAll(string(source), n[0], fmt.Sprintf("{{ item \"%s\" . }}", n[1]))
		source = []byte(str)
	}
	return source
}

func (t *TemplMerger) Save() error {
	for fileName, bytes := range t.file {
		// override file with merged values
		err := writeToFile(fileName, string(bytes))
		if err != nil {
			return fmt.Errorf("cannot update config file: %s\n", err)
		}
	}
	return nil
}
