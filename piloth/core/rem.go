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
	"net/http"
)

type Rem struct {
	client *Client
	cfg    *ClientConf
	host   *HostInfo
}

func NewRem() (*Rem, error) {
	conf := &Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	cfg := &ClientConf{
		BaseURI:            conf.Get(PilotCtlUri),
		Username:           "_",
		Password:           "_",
		InsecureSkipVerify: false,
		Timeout:            60,
	}
	c, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create PilotCtl http client: %s", err)
	}
	i, err := NewHostInfo()
	if err != nil {
		return nil, err
	}
	return &Rem{client: c, cfg: cfg, host: i}, nil
}

// Register the host
func (r *Rem) Register() error {
	i := r.host
	// set the machine id
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
	resp, err := r.client.Post(uri, reg, r.addToken)
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
	uri := fmt.Sprintf("%s/ping/%s", r.cfg.BaseURI, r.host.HostID)
	resp, err := r.client.Post(uri, nil, r.addToken)
	if err != nil {
		return nil, fmt.Errorf("ping failed ping: %s", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("call to the remote service failed: %d - %s", resp.StatusCode, resp.Status)
	}
	// get the commands to execute from the response body
	bytes, err := ioutil.ReadAll(resp.Body)
	var commands []rem.CmdRequest
	err = json.Unmarshal(bytes, &commands)
	if err != nil {
		return nil, fmt.Errorf("cannot read ping response: %s", err)
	}
	return commands, nil
}

func (r *Rem) addToken(req *http.Request, payload Serializable) error {
	payload = nil
	// add an authentication token to the request
	req.Header.Set("Authorization", newToken(r.host.HostID))
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	return nil
}
