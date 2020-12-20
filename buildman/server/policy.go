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
	b, err := ioutil.ReadFile(PolicyFile())
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
	// the polling interval to check for base image changes
	Interval int `json:"interval"`
	// a list of policies to trigger new image builds
	Policies []*policyConf `json:"policies"`
}

type policyConf struct {
	// the name of the policy that MUST correspond to the name of the image build pipeline
	Name string `json:"name"`
	// the policy description
	Description string `json:"description"`
	// the name of the base image
	Base string `json:"base"`
	// the label on the base image containing the created date
	BaseCreated string `json:"app-base-created-label"`
	// the username to connect to the base image registry
	BaseUser string `json:"base-user"`
	// the password to connect to the base image registry
	BasePwd string `json:"base-pwd"`
	// the name of the application image
	App string `json:"app"`
	// the username to connect to the application image registry
	AppUser string `json:"app-user"`
	// the password to connect to the application image registry
	AppPwd string `json:"app-pwd"`
	// the kubernetes namespace where the image build pipeline is to be launched
	Namespace string `json:"namespace"`
	// a flag to enable the polling of base image information
	PollBase bool `json:"pollBase"`
}

func WorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return wd
}

// gets the policy file FQDN
func PolicyFile() string {
	var filename = path.Join(WorkDir(), "policy.json")
	policyPath := os.Getenv("OXBM_POLICY_PATH")
	if len(policyPath) > 0 {
		filename = path.Join(policyPath, "policy.json")
	}
	return filename
}
