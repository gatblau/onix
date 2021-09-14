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
	ctl "github.com/gatblau/onix/pilotctl/types"
	"github.com/gatblau/onix/piloth/job"
	"io/ioutil"
	"net/http"
	"time"
)

type PilotCtl struct {
	client *Client
	cfg    *ClientConf
	host   *HostInfo
	worker *job.Worker
}

func NewPilotCtl(worker *job.Worker, hostInfo *HostInfo) (*PilotCtl, error) {
	conf := &Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	cfg := &ClientConf{
		BaseURI:            conf.getPilotCtlURI(),
		Username:           "_",
		Password:           "_",
		InsecureSkipVerify: true,
		Timeout:            60 * time.Second,
	}
	c, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create PilotCtl http client: %s", err)
	}
	return &PilotCtl{client: c, cfg: cfg, host: hostInfo, worker: worker}, nil
}

// Register the host
func (r *PilotCtl) Register() (*ctl.RegistrationResponse, error) {
	i := r.host
	// set the machine id
	reg := &ctl.RegistrationRequest{
		Hostname:    i.HostName,
		MachineId:   i.MachineId,
		OS:          i.OS,
		Platform:    fmt.Sprintf("%s, %s, %s", i.Platform, i.PlatformFamily, i.PlatformVersion),
		Virtual:     i.Virtual,
		TotalMemory: i.TotalMemory,
		CPUs:        i.CPUs,
		HostIP:      i.HostIP,
	}
	uri := fmt.Sprintf("%s/register", r.cfg.BaseURI)
	resp, err := r.client.Post(uri, reg, r.addToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("the request failed with error: %d - %s", resp.StatusCode, resp.Status)
	}
	var result ctl.RegistrationResponse
	op, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(op, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Ping send a ping to the remote server
func (r *PilotCtl) Ping() (ctl.PingResponse, error) {
	// is there a result from a job ready?
	var (
		payload Serializable
		result  *ctl.JobResult
	)
	// check if the worker has a job result to be sent to pilot control
	result, ok := r.worker.Result()
	if ok {
		payload = result
	} else {
		payload = nil
	}
	uri := fmt.Sprintf("%s/ping", r.cfg.BaseURI)
	resp, err := r.client.Post(uri, payload, r.addToken)
	if err != nil {
		return ctl.PingResponse{}, fmt.Errorf("ping failed ping: %s", err)
	}
	if resp.StatusCode > 299 {
		return ctl.PingResponse{}, fmt.Errorf("call to the remote service failed: %d - %s", resp.StatusCode, resp.Status)
	}
	// get the commands to execute from the response body
	bytes, err := ioutil.ReadAll(resp.Body)
	var pingResponse ctl.PingResponse
	err = json.Unmarshal(bytes, &pingResponse)
	if err != nil {
		return ctl.PingResponse{}, fmt.Errorf("cannot read ping response: %s", err)
	}
	return pingResponse, nil
}

func (r *PilotCtl) addToken(req *http.Request, payload Serializable) error {
	payload = nil
	// add an authentication token to the request
	req.Header.Set("Authorization", newToken(r.host.HostUUID, r.host.HostIP, r.host.HostName))
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	return nil
}
