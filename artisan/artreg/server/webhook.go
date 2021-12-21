/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/hashicorp/go-uuid"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// the type of action on an package that is published by the webhook
type WebHookAction int

func NewWhAction(value string) *WebHookAction {
	action := new(WebHookAction)
	action.FromString(value)
	return action
}

// actions on registry packages
const (
	// package has been uploaded
	WhUploaded WebHookAction = iota
	// tag has been added to the package tag list
	WhTagged
	// tag has been removed from the package tag list
	WhUnTagged
	// the package has been removed from the registry
	WhRemoved
	// unknown action
	WhUnknown
)

// returns a String representation of the action
func (action WebHookAction) String() string {
	switch action {
	case WhUploaded:
		return "UPLOAD"
	case WhTagged:
		return "TAGGED"
	case WhUnTagged:
		return "UNTAGGED"
	case WhRemoved:
		return "REMOVED"
	default:
		return "UNKNOWN"
	}
}

// load the action from a string
func (action WebHookAction) FromString(actionStrValue string) WebHookAction {
	switch strings.ToUpper(actionStrValue) {
	case "UPLOAD":
		return WhUploaded
	case "TAGGED":
		return WhUnTagged
	case "UNTAGGED":
		return WhUnTagged
	case "REMOVED":
		return WhRemoved
	default:
		return WhUnknown
	}
}

// web hook information
type WebHook struct {
	// action on the package
	action WebHookAction `json:"action"`
	// the repository group
	group string `json:"group"`
	// the repository name
	name string `json:"name"`
	// the package tag
	tag string `json:"tag"`
	// the package unique identifier
	artieId string `json:"artie_id"`
}

// configuration for a web hook
type WebHookConfig struct {
	// the unique webhook identifier
	Id string `json:"id"`
	// Actions that should trigger the webhook
	Actions []WebHookAction `json:"actions"`
	// the repository Group for the webhook
	Group string `json:"group"`
	// the repository Name for the webhook
	Name string `json:"name"`
	// the webhook URI
	Uri string `json:"uri"`
	// the webhook URI user
	Uname string `json:"uname"`
	// the webhook URI password
	Pwd string `json:"pwd"`
}

func (c WebHookConfig) equals(h *WebHookConfig) bool {
	return c.Name == h.Name && c.Group == h.Group && c.Uri == h.Uri
}

type WebHooks struct {
	Hooks []*WebHookConfig
}

func NewWebHooks() *WebHooks {
	return &WebHooks{
		Hooks: []*WebHookConfig{},
	}
}

func (wh *WebHooks) save() error {
	b, err := json.Marshal(wh)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(wh.file(), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (wh *WebHooks) load() error {
	_, err := os.Stat(wh.file())
	if os.IsNotExist(err) {
		err := os.MkdirAll(wh.path(), os.ModePerm)
		if err != nil {
			fmt.Printf("cannot create webhooks config path: %s", err)
		}
		return wh.save()
	}
	b, err := ioutil.ReadFile(wh.file())
	if err != nil {
		return err
	}
	return json.Unmarshal(b, wh)
}

func (wh *WebHooks) update(h *WebHookConfig) {
	// if the webhook configuration does not exist
	contained, position := wh.contains(h)
	if !contained {
		// adds it
		wh.Hooks = append(wh.Hooks, h)
	} else {
		// updates it
		wh.Hooks[position] = h
	}
}

func (wh *WebHooks) contains(h *WebHookConfig) (bool, int) {
	for ix, hook := range wh.Hooks {
		if hook.equals(h) {
			return true, ix
		}
	}
	return false, -1
}

func (wh *WebHooks) path() string {
	return path.Join(core.RegistryPath(), "hooks")
}

func (wh *WebHooks) file() string {
	return path.Join(wh.path(), "config.json")
}

func (wh *WebHooks) GetList(group string, name string) []*WebHookConfig {
	var hooks = make([]*WebHookConfig, 0)
	for _, hook := range wh.Hooks {
		if hook.Group == group && hook.Name == name {
			hooks = append(hooks, hook)
		}
	}
	return hooks
}

func (wh *WebHooks) Get(group string, name string, whId string) *WebHookConfig {
	for _, hook := range wh.Hooks {
		if hook.Group == group && hook.Name == name && hook.Id == whId {
			return hook
		}
	}
	return nil
}

func (wh *WebHooks) Add(config *WebHookConfig) (string, error) {
	confs := wh.GetList(config.Group, config.Name)
	// check if the configuration for the uri already exists
	var (
		position = -1
		id       string
	)
	// generate a new id
	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("cannot generate configuratio Id: %s", err)
	}
	for ix, conf := range confs {
		if conf.Uri == config.Uri {
			position = ix
			break
		}
	}
	// if the configuration exists
	if position != -1 {
		// updates it
		// preserve the id from the existing configuration
		id = wh.Hooks[position].Id
		// set the id in the new configuration
		config.Id = id
		// update the configuration
		wh.Hooks[position] = config
	} else {
		// add it
		// set the id in the new configuration
		config.Id = id
		// appends the configuration
		wh.Hooks = append(wh.Hooks, config)
	}
	err = wh.save()
	if err != nil {
		return "", fmt.Errorf("cannot save webhook configuration: %s", err)
	}
	return id, nil
}

func (wh *WebHooks) Remove(group string, name string, id string) bool {
	var position = -1
	for ix, hook := range wh.Hooks {
		if hook.Group == group && hook.Name == name && hook.Id == id {
			position = ix
			break
		}
	}
	// Remove the element at index "position" from a.
	l := len(wh.Hooks)
	wh.Hooks[position] = wh.Hooks[l-1] // Copy last element to index i.
	wh.Hooks[l-1] = nil                // Erase last element (write zero value).
	wh.Hooks = wh.Hooks[:l-1]          // Truncate slice.
	return position != -1
}
