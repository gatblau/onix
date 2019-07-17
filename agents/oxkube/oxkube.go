/*
   Onix Kube - Copyright (c) 2019 by www.gatblau.org

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
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type OxKube struct {
	config *Config
	log    *logrus.Entry
	client *Client
	ready  bool
}

func (k *OxKube) start() error {
	var err error
	// load the configuration file
	err = k.loadConfig()
	if err != nil {
		return err
	}
	// initialises the Onix REST client
	k.client, err = NewClient(k.log, k.config)
	if err != nil {
		return err
	}
	// checks if a meta model for K8S is defined in Onix
	k.log.Tracef("Checking if the KUBE meta-model is defined in Onix.")
	var (
		exist    bool
		attempts int
		interval time.Duration = 30 // the interval to wait for reconnection
	)
	for {
		exist, err = k.client.modelExists()
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "500") {
			// there is a CMDB error so exit
			k.log.Errorf("Can't connect to Onix: %s.", err)
			return err
		} else {
			attempts = attempts + 1
			k.log.Warnf("Can't connect to Onix: %s. "+
				"Attempt %s, waiting before attempting to connect again.", err, strconv.Itoa(attempts))
			time.Sleep(interval * time.Second)
		}
	}
	// if not...
	if !exist {
		// creates a meta model
		k.log.Tracef("The KUBE meta-model is not yet defined in Onix, proceeding to create it.")
		result := k.client.putModel()
		if result.Error {
			k.log.Errorf("Can't create KUBE meta-model: %s", result.Message)
			return errors.New(result.Message)
		}
	} else {
		k.log.Tracef("KUBE meta-model found in Onix.")
	}
	// the webhook is ready to receive incoming connections
	k.ready = true
	// start the configured consumer
	switch k.config.Consumers.Consumer {
	case "webhook":
		k.log.Tracef("Webhook consumer has been selected.")
		wh := Webhook{
			log:    k.log,
			config: k.config.Consumers.Webhook,
			ready:  k.ready,
		}
		k.log.Tracef("Starting the webhook consumer.")
		wh.Start(k.client)
	case "broker":
		k.log.Tracef("Broker consumer has been selected.")
		panic("Broker consumer is not implemented.")
	default:
		k.log.Tracef("No consumer has been selected.")
		panic(fmt.Sprintf("Mode '%s' is not implemented.", k.config.Consumers.Consumer))
	}
	return nil
}

// load the configuration file
func (k *OxKube) loadConfig() error {
	// loads the configuration
	c, err := NewConfig()
	if err == nil {
		k.config = &c
	} else {
		return err
	}

	// adds the platform field to the logger
	k.log = logrus.WithFields(logrus.Fields{
		"Id": k.config.Id,
	})

	// try and parse the logging level in the configuration
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		// if the value was not recognised then return the error
		k.log.Errorf("Failed to recognise value LogLevel entry in the configuration: %s.", err)
		return err
	}
	// otherwise sets the logging level for the entire system
	logrus.SetLevel(level)
	k.log.Infof("%s has been set as the logger level.", strings.ToUpper(c.LogLevel))
	return nil
}
