/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0

  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import "strings"

// the type of configuration being handled
type confType int

// the values for confType
const (
	TypeFile confType = iota
	TypeHttp
	TypeEnvironment
	TypeUnknown
)

// how to turn int to string representations
func (c confType) String() string {
	switch c {
	case TypeFile:
		return "file"
	case TypeHttp:
		return "http"
	case TypeEnvironment:
		return "environment"
	default:
		return "??unknown??"
	}
}

// return a confType from its string representation
func (c confType) parse(value string) confType {
	switch value {
	case "file":
		return TypeFile
	case "http":
		return TypeHttp
	case "environment":
		return TypeEnvironment
	default:
		return TypeUnknown
	}
}

// the trigger for reloading a configuration
type trigger int

const (
	TriggerRestart trigger = iota
	TriggerGet
	TriggerPost
	TriggerPut
	TriggerSignal
	TriggerUnknown
)

// parse the passed string into a trigger
func (t trigger) parse(value string) trigger {
	switch strings.ToLower(value) {
	case "restart":
		return TriggerRestart
	case "get":
		return TriggerGet
	case "post":
		return TriggerPost
	case "put":
		return TriggerPut
	default:
		if strings.HasPrefix(value, "signal:") {
			return TriggerSignal
		}
		return TriggerUnknown
	}
}

// how to turn int to string representations
func (t trigger) String() string {
	switch t {
	case TriggerRestart:
		return "restart"
	case TriggerGet:
		return "get"
	case TriggerPost:
		return "post"
	case TriggerPut:
		return "put"
	case TriggerSignal:
		return "signal"
	default:
		return "??unknown??"
	}
}
