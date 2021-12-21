package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import "testing"

func TestParseName(t *testing.T) {
	data := map[string]bool{
		// correct
		"localhost/my-group/my-name":           true,
		"localhost:8082/my-group/my-name":      true,
		"localhost:8082/my-group/my-name:v1-0": true,
		"my-group/my-name":                     true,
		"my-name":                              true,
		"my-name:v1":                           true,
		// invalid character in domain
		"loc%$alhost/my-group/my-name": false,
		// domain cannot start with hyphen
		"-localhost/my-group/my-name": false,
		"localhost/my-group/:ggd":     true,
		":ggd":                        true,
		// domain cannot start with colon
		":localhost/my-group/:ggd": false,
		// missing group and name
		"127.0.0.1:884": false,
		// missing name
		"127.0.0.1:884/my-group":         false,
		"127.0.0.1:884/my-group/my-name": true,
		// no protocol scheme allowed
		"http://127.0.0.1:884/my-group/my-name":  false,
		"https://127.0.0.1:884/my-group/my-name": false,
		"tcp://127.0.0.1:884/my-group/my-name":   false,
		"ws://127.0.0.1:884/my-group/my-name":    false,
		"127.0.0.1:8f84/my-group/my-name":        false,
	}
	for name, valid := range data {
		_, err := ParseName(name)
		if valid && err != nil {
			t.Errorf(err.Error())
		}
		if !valid && err == nil {
			t.Errorf("name %s should be invalid", name)
		}
	}
}
