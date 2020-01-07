/*
   Terraform Backend for Onix - Copyright (c) 2019 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
package main

import (
	. "gatblau.org/onix/wapic"
	log "github.com/sirupsen/logrus"
)

type Backend struct {
	config *Config
	log    *log.Entry
	client *Client
	ready  bool
}

func (t *Backend) start() error {
	var err error

	// load the configuration file
	c, err := NewConfig()
	if err != nil {
		return err
	} else {
		// sets the configuration
		t.config = &c
		// sets the logger
		t.log = log.WithFields(log.Fields{"Id": c.Id})
	}

	// initialises the Onix REST client
	t.client, err = New(t.log, t.config.Onix)
	if err != nil {
		return err
	}
	// checks if a meta model for K8S is defined in Onix
	t.log.Tracef("Checking if the KUBE meta-model is defined in Onix.")

	return nil
}
