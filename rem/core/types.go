package core

import (
	"bytes"
	"encoding/json"
	"time"
)

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// Cmd command information for remote host execution
type Cmd struct {
	Id       int               `json:"id"`
	Package  string            `json:"package"`
	Function string            `json:"function"`
	Input    map[string]string `json:"input"`
}

// CmdRequest a command for execution with a job reference
type CmdRequest struct {
	cmd *Cmd
	ref string
}

// Host  host monitoring information
type Host struct {
	Name      string `json:"name"`
	Customer  string `json:"customer"`
	Region    string `json:"region"`
	Location  string `json:"location"`
	Connected bool   `json:"connected"`
	LastSeen  string `json:"last_seen"`
}

// Registration information for host registration
type Registration struct {
	Hostname string `json:"hostname"`
	// github.com/denisbrodbeck/machineid
	MachineId   string  `json:"machine_id"`
	OS          string  `json:"os"`
	Platform    string  `json:"platform"`
	Virtual     bool    `json:"virtual"`
	TotalMemory float64 `json:"total_memory"`
	CPUs        int     `json:"cpus"`
}

// Reader Get a JSON bytes reader for the Serializable
func (r *Registration) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Bytes Get a []byte representing the Serializable
func (r *Registration) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}

// Event host events
type Event struct {
	// 0: host up, 1: host down, 2: network up, 3: network down
	Type int       `json:"type"`
	Time time.Time `json:"time"`
}

type Events []Event

// Job a job to be executed on one or more hosts
type Job struct {
	Id        int      `json:"id"`
	HostId    []string `json:"host_id"`
	CmdId     string   `json:"cmd_id"`
	Created   string   `json:"created,omitempty"`
	Started   string   `json:"started,omitempty"`
	Completed string   `json:"completed,omitempty"`
}

// ToJson convert the passed-in object to a JSON byte slice
// NOTE: json.Marshal is purposely not used as it will escape any < > characters
func ToJson(object interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// switch off the escaping!
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}

type Region struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Location struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	RegionKey string `json:"region_key"`
}

type Admission struct {
	Key    string   `json:"key"`
	Active bool     `json:"active"`
	Tag    []string `json:"tag"`
}
