/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0

  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"testing"
)

func TestParserOK(t *testing.T) {
	p := new(parser)
	l, err := p.parse(confWellFormed)
	if err != nil || len(l) != 2 {
		t.Fail()
	}
}

func TestParserEndNoStart(t *testing.T) {
	p := new(parser)
	_, err := p.parse(confEndNoStart)
	// we should see an error
	if err == nil {
		t.Fail()
	}
}

// test data
const (
	// an example of a well-formed configuration
	confWellFormed = `
+++toml
# the way the configuration is loaded
# possible values are "file", "http", "environment"
Type = "file"

# the location of the file or http endpoint to call (valid for type=file and type=http)
Path = "app.toml"

# the trigger used to reload the configuration
# possible values are: signal, restart, get, post, put
Trigger = "signal:SIGHUP"

# if type=http and trigger=post/put the content type to be passed in the request
ContentType = "application/txt"
+++
[Log]
    Level = "trace"

[Banner]
    Type = "success"
    Message = "hello probare"

+++json
{
	"type": "file",
	"path": "/cfg/secrets",
	"trigger": "put",
	"content_type": "application/txt"
}
+++
User = "gatblau"
Pwd = "FMy2nNzh"
`
	// an example of a configuration containing a front matter end with no start
	confEndNoStart = `
+++toml
Type = "file"
Path = "app.toml"
Trigger = "signal:SIGHUP"
ContentType = "application/txt"
+++
dxsdxdxsdx
+++
# end mark before start mark
`
)
