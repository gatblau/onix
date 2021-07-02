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
	ctl "github.com/gatblau/onix/pilotctl/core"
	"github.com/gatblau/onix/piloth/cmd"
	"io/ioutil"
	"net/http"
)

type PilotCtl struct {
	client *Client
	cfg    *ClientConf
	host   *HostInfo
	worker *cmd.Worker
}

func NewPilotCtl(worker *cmd.Worker) (*PilotCtl, error) {
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
	return &PilotCtl{client: c, cfg: cfg, host: i, worker: worker}, nil
}

// Register the host
func (r *PilotCtl) Register() error {
	i := r.host
	// set the machine id
	reg := &ctl.Registration{
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
func (r *PilotCtl) Ping() (ctl.CmdRequest, error) {
	// is there a result from a job ready?
	var (
		payload Serializable
		result  cmd.Result
	)
	select {
	case result = <-r.worker.Queue.Results:
		// a result was received from the channel
		payload = &result
	default:
		// no result is available
		payload = nil
	}
	uri := fmt.Sprintf("%s/ping/%s", r.cfg.BaseURI, r.host.HostID)
	resp, err := r.client.Post(uri, payload, r.addToken)
	if err != nil {
		return ctl.CmdRequest{}, fmt.Errorf("ping failed ping: %s", err)
	}
	if resp.StatusCode > 299 {
		return ctl.CmdRequest{}, fmt.Errorf("call to the remote service failed: %d - %s", resp.StatusCode, resp.Status)
	}
	// get the commands to execute from the response body
	bytes, err := ioutil.ReadAll(resp.Body)
	var command ctl.CmdRequest
	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return ctl.CmdRequest{}, fmt.Errorf("cannot read ping response: %s", err)
	}
	return command, nil
}

func (r *PilotCtl) addToken(req *http.Request, payload Serializable) error {
	payload = nil
	// add an authentication token to the request
	req.Header.Set("Authorization", newToken(r.host.HostID))
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	return nil
}
