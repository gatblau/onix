/*
   Onix Config Manager - OxTerra - Terraform Http Backend for Onix
   Copyright (c) 2018-2020 by www.gatblau.org

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
	"github.com/gatblau/oxc"
	"github.com/hashicorp/terraform/states/statemgr"
	"net/http"
)

type Lock struct {
	stateKey string
	oxc      *oxc.Client
}

// initialise locking functions for a specified state
func NewLock(stateKey string, oxc *oxc.Client) *Lock {
	return &Lock{stateKey: stateKey, oxc: oxc}
}

// check if the struct has been initialised
func (l *Lock) ok() bool {
	return l.stateKey == "" && l.oxc == nil
}

// check if the lock exists returning the lock Id if exists or an empty string if it does not
func (l *Lock) exist(stateKey string) (string, error) {
	lock, err := l.oxc.GetItem(lockKeyItem(stateKey))
	if err != nil {
		return "", err
	}
	return lock.Attribute["id"].(string), err
}

// acquires a new lock
func (l *Lock) lock(info *statemgr.LockInfo) (error, int) {
	// try and retrieve lock if exists
	lock, err := l.oxc.GetItem(lockKeyItem(l.stateKey))
	// if the lock already exists
	if err == nil && lock != nil {
		// if the id of the lock in the database is different from the id passed-in
		if lock.Attribute["id"] != info.ID {
			// someone else has the lock, then it returns 409
			return err, http.StatusLocked
		}
	}
	// can go ahead and acquire the lock
	attrs := make(map[string]interface{})
	attrs["id"] = info.ID
	attrs["version"] = info.Version
	attrs["path"] = info.Path
	attrs["info"] = info.Info
	attrs["who"] = info.Who
	result, err := l.oxc.PutItem(&oxc.Item{
		Key:         lockKey(l.stateKey),
		Name:        fmt.Sprintf("Lock for %s", lockKey(l.stateKey)),
		Description: "",
		Type:        TfLockType,
		Attribute:   attrs,
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if result.Error {
		return errors.New(result.Message), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// releases an existing lock
func (l *Lock) unlock() (error, int) {
	// try and retrieve lock if exists
	result, err := l.oxc.DeleteItem(lockKeyItem(l.stateKey))
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if result.Error {
		return errors.New(result.Message), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// gets a keyed TF Lock Item
func lockKeyItem(sKey string) *oxc.Item {
	return &oxc.Item{Key: lockKey(sKey)}
}

// gets the TF log key
func lockKey(key string) string {
	return fmt.Sprintf("TF_STATE_LOCK_%s", key)
}
