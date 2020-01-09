/*
   Terraform Http Backend - Onix - Copyright (c) 2018 by www.gatblau.org

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
	"errors"
	. "gatblau.org/onix/wapic"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Backend struct {
	config *Config
	log    *log.Entry
	client *Client
	ready  bool
}

func (b *Backend) start() error {
	var err error

	// load the configuration file
	c, err := NewConfig()
	if err != nil {
		return err
	} else {
		// sets the configuration
		b.config = &c
		// sets the logger
		b.log = log.WithFields(log.Fields{"Id": c.Id})
	}

	// initialises the Onix REST client
	b.client, err = New(b.log, b.config.Onix)
	if err != nil {
		return err
	}
	// checks if a meta model for Terraform is defined in Onix
	b.log.Tracef("Checking if the TERRAFORM meta-model is defined in Onix.")
	var (
		mmodel   *MetaModel = NewModel(b.log, b.client)
		exist    bool
		attempts int
		interval time.Duration = 30 // the interval to wait for reconnection
	)
	for {
		exist, err = mmodel.exists()
		if err == nil {
			break
		}
		attempts = attempts + 1
		b.log.Warnf("Can't connect to Onix: %s. "+
			"Attempt %s, waiting before attempting to connect again.", err, strconv.Itoa(attempts))
		time.Sleep(interval * time.Second)
	}
	// if not...
	if !exist {
		// creates a meta model
		b.log.Tracef("The TERRAFORM meta-model is not yet defined in Onix, proceeding to create it.")
		result := mmodel.create()
		if result.Error {
			b.log.Errorf("Can't create TERRAFORM meta-model: %s", result.Message)
			return errors.New(result.Message)
		}
	} else {
		b.log.Tracef("TERRAFORM meta-model found in Onix.")
	}
	// the backend is ready to receive http connections
	b.ready = true
	// start the service listener
	svc := NewService(*b)
	svc.Start()
	return nil
}
