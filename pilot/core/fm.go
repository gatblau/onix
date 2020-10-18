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
	"net/url"
)

// pilot configuration for a specific application configuration
// informs pilot how to manage the application configuration
type frontMatter struct {
	// the type of configuration to handle (i.e. file, http or environment)
	Type string `json:"type"`
	// the file path to the configuration (only applicable if Type=file)
	Path string `json:"path"`
	// the URI used to issue http requests to manage the configuration (only applicable if Type is not file)
	Uri string `json:"uri"`
	// the Content-Type http header used for POST/PUT requests to the URI
	ContentType string `json:"contentType"`
	// the username to authenticate with the URI (if required)
	User string `json:"user"`
	// the password to authenticate with the URI (if required)
	Pwd string `json:"pwd"`
	// the mechanism used to inform the application to reload the configuration (signal, get, post, put or reload)
	Trigger string `json:"trigger"`
}

// validate the front matter
func (fm *frontMatter) valid() (bool, error) {
	// a configuration type is always required
	if ok, err := fm.exists("type", fm.Type); !ok {
		return ok, err
	}
	// does it contain an implemented configuration type?
	if ok, err := fm.validConfType(fm.Type); !ok {
		return ok, err
	}
	// a path is only required if the configuration type is a file
	if ok, err := fm.existsWhen("path", fm.Path, fm.Type == TypeFile.String()); !ok {
		return ok, err
	}
	// a URI is only required if the configuration type is not a file
	if ok, err := fm.existsWhen("uri", fm.Uri, fm.Type != TypeFile.String()); !ok {
		return ok, err
	}
	// does it contain an valid URI?
	if fm.Type != TypeFile.String() {
		_, err := url.ParseRequestURI(fm.Uri)
		if err != nil {
			return false, err
		}
	}
	// the trigger should always be provided
	if ok, err := fm.exists("trigger", fm.Trigger); !ok {
		return ok, err
	}
	// does it contain an implemented trigger?
	if ok, err := fm.validTrigger(fm.Trigger); !ok {
		return ok, err
	}
	// set a default value for the content type of not provided
	if len(fm.ContentType) == 0 {
		fm.ContentType = "text/plain"
	}
	// validate the content type
	if ok, err := fm.validContentType(fm.ContentType); !ok {
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

// check if a field exists
func (fm *frontMatter) exists(field string, value string) (bool, error) {
	if len(value) == 0 {
		return false, errors.New(fmt.Sprintf("the front matter is missing value for '%s'", field))
	}
	return true, nil
}

// check if a field exists only if the passed in condition is true
func (fm *frontMatter) existsWhen(field string, value string, condition bool) (bool, error) {
	if condition {
		if len(value) == 0 {
			return false, errors.New(fmt.Sprintf("the front matter is missing value for '%s'", field))
		}
		return true, nil
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

// check if the provided configuration type matches one of the valid values
func (fm *frontMatter) validConfType(cType string) (bool, error) {
	t := new(confType)
	valid := t.parse(cType) != TypeUnknown
	if !valid {
		return false, errors.New(fmt.Sprintf("invalid configuration type %s", cType))
	}
	return valid, nil
}

// check if the provided reload trigger matches one of the valid values
func (fm *frontMatter) validTrigger(trig string) (bool, error) {
	t := new(trigger)
	valid := t.parse(trig) != TriggerUnknown
	if !valid {
		return false, errors.New(fmt.Sprintf("invalid reload trigger type %s", trig))
	}
	return valid, nil
}
