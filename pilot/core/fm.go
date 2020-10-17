/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0

  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"errors"
	"fmt"
	"strings"
)

// pilot configuration stored in Onix along with the application configuration
type frontMatter struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Trigger     string `json:"trigger"`
	ContentType string `json:"contentType"`
}

// validate the front matter
func (fm *frontMatter) valid() (bool, error) {
	// check for required fields
	if ok, err := fm.notEmpty("type", fm.Type); !ok {
		return ok, err
	}
	if ok, err := fm.notEmpty("path", fm.Type); !ok {
		return ok, err
	}
	if ok, err := fm.notEmpty("trigger", fm.Trigger); !ok {
		return ok, err
	}
	// default value for content type
	if len(fm.ContentType) == 0 {
		fm.ContentType = "text/plain"
	}
	// validate the content type
	if ok, err := fm.validContentType(fm.ContentType); !ok {
		return ok, err
	}
	// validates trigger
	if ok, err := fm.validTrigger(fm.Trigger); !ok {
		return ok, err
	}
	return true, nil
}

func (fm *frontMatter) typeVal() confType {
	var t confType
	return t.parse(fm.Type)
}

func (fm *frontMatter) triggerVal() trigger {
	var t trigger
	return t.parse(fm.Trigger)
}

// check a field is not empty
func (fm *frontMatter) notEmpty(field string, value string) (bool, error) {
	if len(value) == 0 {
		return false, errors.New(fmt.Sprintf("the front matter is missing value for '%s'", field))
	}
	return true, nil
}

// check the validity of the content type
func (fm *frontMatter) validContentType(contentType string) (bool, error) {
	switch contentType {
	case "text/plain":
		return true, nil
	case "application/json":
		return true, nil
	case "application/xml":
		return true, nil
	case "application/x-yaml":
		fallthrough
	case "text/yaml":
		return true, nil
	default:
		return false, errors.New(fmt.Sprintf("invalid content type if front matter: %s", contentType))
	}
}

// check validity of trigger
func (fm *frontMatter) validTrigger(trigger string) (bool, error) {
	switch strings.ToLower(trigger) {
	case "":
		fallthrough
	case "restart":
		fallthrough
	case "get":
		fallthrough
	case "post":
		fallthrough
	case "put":
		return true, nil
	default:
		if strings.HasPrefix(strings.ToLower(trigger), "signal:") {
			s := trigger[7:]
			if !(s == "SIGHUP" || s == "SIGUSR1" || s == "SIGUSR2") {
				return false, errors.New(fmt.Sprintf("invalid signal if front matter: %s", s))
			}
			return true, nil
		}
	}
	return false, errors.New(fmt.Sprintf("invalid trigger if front matter: %s", trigger))
}
