package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
)

// CmdInfo all the information required by pilot to execute a command
type CmdInfo struct {
	JobId         int64       `json:"job_id"`
	Package       string      `json:"package"`
	Function      string      `json:"function"`
	User          string      `json:"user"`
	Pwd           string      `json:"pwd"`
	Verbose       bool        `json:"verbose"`
	Containerised bool        `json:"containerised"`
	Input         *data.Input `json:"input,omitempty"`
}

func (c *CmdInfo) Value() string {
	var artCmd string
	// if command is to run in a runtime
	if c.Containerised {
		// use art exec
		artCmd = "exec"
	} else {
		// otherwise, use art exe
		artCmd = "exe"
	}
	// if user credentials for the Artisan registry were provided
	if len(c.User) > 0 && len(c.Pwd) > 0 {
		// pass the credentials to the art cli
		return fmt.Sprintf("art %s -u %s:%s %s %s", artCmd, c.User, c.Pwd, c.Package, c.Function)
	}
	// otherwise run the command without credentials (assume public registry)
	return fmt.Sprintf("art %s %s %s", artCmd, c.Package, c.Function)
}

func (c *CmdInfo) Env() []string {
	var vars []string
	// append vars
	for _, v := range c.Input.Var {
		vars = append(vars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	// append secrets
	for _, s := range c.Input.Secret {
		vars = append(vars, fmt.Sprintf("%s=%s", s.Name, s.Value))
	}
	return vars
}

func (c *CmdInfo) Envar() *merge.Envar {
	return merge.NewEnVarFromSlice(c.Env())
}

func (c *CmdInfo) PrintEnv() string {
	var vars bytes.Buffer
	vars.WriteString("printing variables passed to the shell\n{\n")
	for _, v := range c.Input.Var {
		vars.WriteString(fmt.Sprintf("%s=%s\n", v.Name, v.Value))
	}
	vars.WriteString("}\n")
	return vars.String()
}
