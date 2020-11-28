/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"github.com/gatblau/onix/artie/core"
	"path"
	"strings"
)

// the type of action on an artefact that is published by the webhook
type WebHookAction int

func NewWhAction(value string) *WebHookAction {
	action := new(WebHookAction)
	action.FromString(value)
	return action
}

// actions on registry artefacts
const (
	// artefact has been uploaded
	WhUploaded WebHookAction = iota
	// tag has been added to the artefact tag list
	WhTagged
	// tag has been removed from the artefact tag list
	WhUnTagged
	// the artefact has been removed from the registry
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
	// action on the artefact
	action WebHookAction `json:"action"`
	// the repository group
	group string `json:"group"`
	// the repository name
	name string `json:"name"`
	// the artefact tag
	tag string `json:"tag"`
	// the artefact unique identifier
	artieId string `json:"artie_id"`
}

// configuration for a web hook
type WebHookConfig struct {
	// actions that should trigger the webhook
	actions []WebHookAction
	// the repository group for the webhook
	group string
	// the repository name for the webhook
	name string
	// the webhook URI
	uri string
	// the webhook URI user
	uname string
	// the webhook URI password
	pwd string
}

func (c WebHookConfig) equals(h WebHookConfig) bool {
	return c.name == h.name && c.group == h.group && c.uri == h.uri
}

type WebHooks struct {
	hooks []WebHookConfig
}

func NewWebHooks() *WebHooks {
	return &WebHooks{
		hooks: []WebHookConfig{},
	}
}

func (wh *WebHooks) save() error {
	return nil
}

func (wh *WebHooks) load() error {
	return nil
}

func (wh *WebHooks) update(h WebHookConfig) {
	// if the webhook configuration does not exist
	contained, position := wh.contains(h)
	if !contained {
		// adds it
		wh.hooks = append(wh.hooks, h)
	} else {
		// updates it
		wh.hooks[position] = h
	}
}

func (wh *WebHooks) contains(h WebHookConfig) (bool, int) {
	for ix, hook := range wh.hooks {
		if hook.equals(h) {
			return true, ix
		}
	}
	return false, -1
}

func (wh *WebHooks) path() string {
	return path.Join(core.RegistryPath(), "hooks")
}
