//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugins

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Conf struct {
	value map[string]interface{}
}

func NewConf(json string) (*Conf, error) {
	r := &Conf{
		value: make(map[string]interface{}),
	}
	err := r.FromJSON(json)
	return r, err
}

func (c *Conf) GetString(key string) (string, bool) {
	k := strings.Split(key, ".")
	if len(k) == 1 {
		v := c.value[strings.ToLower(k[0])]
		if v, ok := v.(string); ok {
			return v, true
		} else {
			fmt.Printf("config key %s not found", k[0])
		}
	}
	if len(k) == 2 {
		if m, ok := c.value[strings.ToLower(k[0])].(map[string]interface{}); ok {
			if v, ok := m[strings.ToLower(k[1])].(string); ok {
				return v, true
			} else {
				fmt.Printf("config key %s not found", k[1])
			}
		} else {
			fmt.Printf("config key %s not found", k[0])
		}
	}
	return "", false
}

func (c *Conf) FromJSON(jsonString string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &m)
	c.value = m
	return err
}
