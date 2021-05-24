package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"encoding/json"
	"fmt"
	rem "github.com/gatblau/onix/rem/core"
	"io/ioutil"
)

type Rem struct {
	client *Client
	cfg    *ClientConf
	host   string
}

func NewRem() (*Rem, error) {
	conf := &Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	cfg := &ClientConf{
		BaseURI:            conf.Get(PilotRemUri),
		InsecureSkipVerify: false,
		Username:           conf.Get(PilotRemUsername),
		Password:           conf.Get(PilotRemPassword),
		Timeout:            60,
	}
	c, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Rem{client: c, cfg: cfg}, nil
}

// Register the host
func (r *Rem) Register() error {
	i, err := NewHostInfo()
	if err != nil {
		return err
	}
	// set the machine id
	r.host = i.HostID
	reg := &rem.Registration{
		Hostname:    i.HostName,
		MachineId:   i.HostID,
		OS:          i.OS,
		Platform:    fmt.Sprintf("%s, %s, %s", i.Platform, i.PlatformFamily, i.PlatformVersion),
		Virtual:     i.Virtual,
		TotalMemory: i.TotalMemory,
		CPUs:        i.CPUs,
	}
	uri := fmt.Sprintf("%s/register", r.cfg.BaseURI)
	resp, err := r.client.Post(uri, reg, r.client.addHttpHeaders)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("the request failed with error: %d - %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// Ping send a ping to the remote server
func (r *Rem) Ping() ([]rem.CmdRequest, error) {
	// check teh host has been registered
	if len(r.host) == 0 {
		return nil, fmt.Errorf("can't ping if not registered")
	}
	uri := fmt.Sprintf("%s/ping/%s", r.cfg.BaseURI, r.host)
	resp, err := r.client.Post(uri, nil, r.client.addHttpHeaders)
	if err != nil {
		return nil, fmt.Errorf("cannot execute ping: %s", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("the remote service error: %d - %s", resp.StatusCode, resp.Status)
	}
	// get the response body
	bytes, err := ioutil.ReadAll(resp.Body)
	commands := make([]rem.CmdRequest, 0)
	err = json.Unmarshal(bytes, commands)
	if err != nil {
		return nil, fmt.Errorf("cannot read ping response: %s", err)
	}
	return commands, nil
}
