/*
  Onix Config Manager - Build Manager
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func NewPolicyConfig() (*policyConfig, error) {
	wd := WorkDir()
	b, err := ioutil.ReadFile(path.Join(wd, "policy.json"))
	if err != nil {
		return nil, fmt.Errorf("cannot read policy file: %s", err)
	}
	p := new(policyConfig)
	err = json.Unmarshal(b, p)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal policy file: %s", err)
	}
	return p, nil
}

type policyConfig struct {
	Interval int `json:"interval"`
	Policies []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Base        string `json:"base"`
		User        string `json:"user"`
		Pwd         string `json:"pwd"`
		Namespace   string `json:"namespace"`
		PollBase    bool   `json:"pollBase"`
	} `json:"policies"`
}

func WorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return wd
}
