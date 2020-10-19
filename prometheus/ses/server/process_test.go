/*
   Onix Config Manager - SeS - Onix Webhook Receiver for AlertManager
   Copyright (c) 2020 by www.gatblau.org

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
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/prometheus/alertmanager/template"
	"io/ioutil"
	"os/user"
	"testing"
)

// store for configuration items
var itemCache map[string]*oxc.Item

func TestProcess2(t *testing.T) {
	// initialises the item map
	itemCache = make(map[string]*oxc.Item)
	alerts, err := load("alerts.json")
	if err != nil {
		t.Fatal(err)
	}
	err = processAlerts(alerts.Alerts, get, put, "INS")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestProcess(t *testing.T) {
	// initialises the item map
	itemCache = make(map[string]*oxc.Item)
	// defines a variable to read incoming payloads
	var payload template.Data
	// unmarshal 1st state: all etcds are up
	err := unmarshal(&payload, "payload_1_all_up.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = processAlerts(payload.Alerts, get, put, "INS")
	if err != nil {
		t.Error(err)
		return
	}
	checkItems(3, 0)
	// unmarshal 2st state: one etcd is down
	err = unmarshal(&payload, "payload_2_one_down.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = processAlerts(payload.Alerts, get, put, "INS")
	if err != nil {
		t.Error(err)
		return
	}
	checkItems(2, 1)
	// unmarshal 3st state: two etcds are down
	err = unmarshal(&payload, "payload_3_two_down.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = processAlerts(payload.Alerts, get, put, "INS")
	if err != nil {
		t.Error(err)
		return
	}
	checkItems(1, 2)
	// unmarshal 4st state: all etcds are up
	err = unmarshal(&payload, "payload_4_all_up.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = processAlerts(payload.Alerts, get, put, "INS")
	if err != nil {
		t.Error(err)
		return
	}
	checkItems(3, 0)
}

// unmarshals an alert manager payload from a file for testing
func unmarshal(payload *template.Data, file string) error {
	// get the current user
	usr, err := user.Current()
	if err != nil {
		return err
	}
	// load alerts from file
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/go/src/github.com/gatblau/onix/prometheus/ses/test/%s", usr.HomeDir, file))
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(dat), payload)
	if err != nil {
		return err
	}
	return nil
}

// stores config items in memory for testing
func put(item *oxc.Item) (*oxc.Result, error) {
	// puts the item in the cache
	itemCache[item.Key] = item
	return &oxc.Result{
		Changed:   false,
		Error:     false,
		Message:   "",
		Operation: "",
		Ref:       "",
	}, nil
}

// get a config items by its natural key
func get(item *oxc.Item) (*oxc.Item, error) {
	// retrieve the item from the cache
	return itemCache[item.Key], nil
}

// check the expected number of up & down events
func checkItems(upCount int, downCount int) error {
	var up, down int
	for _, item := range itemCache {
		if item.Attribute["status"] == "up" {
			up++
		}
		if item.Attribute["status"] == "down" {
			down++
		}
	}
	if up != upCount {
		return errors.New(fmt.Sprintf("Up events: expected %v, instead got %v", upCount, up))
	}
	if up != upCount {
		return errors.New(fmt.Sprintf("Down events: expected %v, instead got %v", downCount, down))
	}
	return nil
}

// load events from a file
func load(file string) (template.Data, error) {
	// defines a variable to read incoming payloads
	var payload template.Data
	// unmarshal 1st state: all etcds are up
	err := unmarshal(&payload, file)
	return payload, err
}
