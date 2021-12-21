package main

/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
)

type Conf struct{}

func NewConf() *Conf {
	return new(Conf)
}

func (c *Conf) getValue(key string) (string, error) {
	val := os.Getenv(key)
	if len(val) == 0 {
		return val, fmt.Errorf("configuration variable '%s' is undefined", key)
	}
	return val, nil
}

func (c *Conf) getOnixWAPIURI() (string, error) {
	return c.getValue("OX_WAPI_URI")
}

func (c *Conf) getOnixWAPIUser() (string, error) {
	return c.getValue("OX_WAPI_UNAME")
}

func (c *Conf) getOnixWAPIPwd() (string, error) {
	return c.getValue("OX_WAPI_PWD")
}
