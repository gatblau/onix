//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"encoding/json"
	"github.com/gatblau/oxc"
	"net/http"
)

// database initialisation information
type DbInit struct {
	Items []Item `json:"items"`
}

type Item struct {
	Action string `json:"action"`
	Script string `json:"script"`
	Admin  bool   `json:"admin"`
	Db     bool   `json:db`
	Vars   []Var  `json:"vars"`
}

type Var struct {
	Name string `json:"name"`
	From string `json:"from"`
}

// get a JSON bytes reader for the Plan
func (init *DbInit) json() (*bytes.Reader, error) {
	jsonBytes, err := init.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Plan
func (init *DbInit) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(init)
	return &b, err
}

// get the Plan in the http Response
func (init *DbInit) decode(response *http.Response) (*DbInit, error) {
	result := new(DbInit)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}
