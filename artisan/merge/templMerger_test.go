package merge

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"testing"
)

func TestMergeUsingFunctions(t *testing.T) {
	e := &Envar{Vars: map[string]string{}}
	// this is an ordinary variable
	e.Vars["TITLE"] = "Example of merging Grouped Variables"

	// these are grouped variables
	// note the naming convention: GROUP-NAME__VARIABLE-NAME__VARIABLE-INDEX
	e.Vars["PORT__NAME__1"] = "Standard TCP"
	e.Vars["PORT__DESC__1"] = "The standard port"
	e.Vars["PORT__VALUE__1"] = "80"

	e.Vars["PORT__NAME__2"] = "Alternative TCP"
	e.Vars["PORT__DESC__2"] = "An alternative http port"
	e.Vars["PORT__VALUE__2"] = "8080"

	e.Vars["PORT__NAME__3"] = "Standard Encrypted"
	e.Vars["PORT__DESC__3"] = "HTTPS port"
	e.Vars["PORT__VALUE__3"] = "443"

	e.Vars["URI__NAME__1"] = "URI 1"
	e.Vars["URI__VALUE__1"] = "www.hhhh.com"

	e.Vars["URI__NAME__2"] = "URI 2"
	e.Vars["URI__VALUE__2"] = "www.hhwedwdehh.com"

	m, _ := NewTemplMerger()
	err := m.LoadTemplates([]string{"test/sample_using_functions.yaml.art"})
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = m.Merge(e)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, bytes := range m.file {
		fmt.Println(string(bytes))
	}
}

func TestMergeUsingOperators(t *testing.T) {
	e := &Envar{Vars: map[string]string{}}
	// this is an ordinary variable
	e.Vars["TITLE"] = "Example of merging Grouped Variables"

	// these are grouped variables
	// note the naming convention: GROUP-NAME__VARIABLE-NAME__VARIABLE-INDEX
	e.Vars["PORT__NAME__1"] = "Standard TCP"
	e.Vars["PORT__DESC__1"] = "The standard port"
	e.Vars["PORT__VALUE__1"] = "80"

	e.Vars["PORT__NAME__2"] = "Alternative TCP"
	e.Vars["PORT__DESC__2"] = "An alternative http port"
	e.Vars["PORT__VALUE__2"] = "8080"

	e.Vars["PORT__NAME__3"] = "Standard Encrypted"
	e.Vars["PORT__DESC__3"] = "HTTPS port"
	e.Vars["PORT__VALUE__3"] = "443"

	e.Vars["URI__NAME__1"] = "URI 1"
	e.Vars["URI__VALUE__1"] = "www.hhhh.com"

	e.Vars["URI__NAME__2"] = "URI 2"
	e.Vars["URI__VALUE__2"] = "www.hhwedwdehh.com"

	m, _ := NewTemplMerger()
	err := m.LoadTemplates([]string{"test/sample_using_operators.yaml.art"})
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = m.Merge(e)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, bytes := range m.file {
		fmt.Println(string(bytes))
	}
}
