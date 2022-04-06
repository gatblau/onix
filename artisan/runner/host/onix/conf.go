/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package onix

import (
	"fmt"
	"os"
	"strconv"
)

type ConfKey string

const (
	ConfOxWapiUri                ConfKey = "OX_WAPI_URI"
	ConfOxWapiUser               ConfKey = "OX_WAPI_USER"
	ConfOxWapiPwd                ConfKey = "OX_WAPI_PWD"
	ConfOxWapiInsecureSkipVerify ConfKey = "OX_WAPI_INSECURE_SKIP_VERIFY"
)

type Conf struct {
}

func NewConf() *Conf {
	return &Conf{}
}

func (c *Conf) get(key ConfKey) string {
	return os.Getenv(string(key))
}

func (c *Conf) getOxWapiUrl() string {
	return c.getValue(ConfOxWapiUri)
}

func (c *Conf) getOxWapiUsername() string {
	return c.getValue(ConfOxWapiUser)
}

func (c *Conf) getOxWapiPassword() string {
	return c.getValue(ConfOxWapiPwd)
}

func (c *Conf) getValue(key ConfKey) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		fmt.Printf("ERROR: variable %s not defined", key)
		os.Exit(1)
	}
	return value
}

func (c *Conf) getOxWapiInsecureSkipVerify() bool {
	b, err := strconv.ParseBool(c.getValue(ConfOxWapiInsecureSkipVerify))
	if err != nil {
		fmt.Printf("ERROR: invalid value for variable %s", ConfOxWapiInsecureSkipVerify)
		os.Exit(1)
	}
	return b
}
