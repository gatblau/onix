/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	ctlCore "github.com/gatblau/onix/pilotctl/core"
	ctl "github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"net/http"
	"strings"

	// "strings"
	"time"
)

type PilotCtl struct {
	client *ctlCore.Client
	cfg    *ctlCore.ClientConf
	host   *ctl.HostInfo
	worker *Worker
}

func NewPilotCtl(worker *Worker, options PilotOptions) (*PilotCtl, error) {
	defer TRA(CE())
	conf := &Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	// breaks URI in activation key based on CSV format
	// note: more than one URI can be configured using a comma separated value list
	// pilot will attempt to connect to an URI in the list following the defined order
	// if a connection is successful, the URI will be elected as the pilotctl URI
	// this process only when pilot re-starts
	ctlURIs := strings.Split(A.CtlURI, ",")
	// start a loop to probe a resolvable URI and return the successful http client
	for _, uri := range ctlURIs {
		// create a new http client for the uri
		cfg := &ctlCore.ClientConf{
			BaseURI:            uri,
			Username:           "_",
			Password:           "_",
			InsecureSkipVerify: options.InsecureSkipVerify,
			Timeout:            5 * time.Minute,
			Proxy:              Proxy(),
		}
		client, clientErr := ctlCore.NewClient(cfg)
		if clientErr != nil {
			return nil, fmt.Errorf("failed to create PilotCtl http client: %s", clientErr)
		}
		core.InfoLogger.Printf("trying to connect to control URI %s\n", uri)
		// issue a http get to the unauthenticated root to check for a valid response
		resp, _ := client.Get(uri, nil)
		// if there is a response
		if resp != nil {
			// and the response is OK
			if resp.StatusCode == 200 {
				// return a client ready  to connect to such endpoint
				core.InfoLogger.Printf("connected to control URI %s\n", uri)
				return &PilotCtl{client: client, cfg: cfg, host: options.Info, worker: worker}, nil
			} else {
				// otherwise, return the error
				return nil, fmt.Errorf("endpoint found but could not connect, reason: %s", resp.Status)
			}
		}
	}
	// if no endpoint was found return error
	return nil, fmt.Errorf("cannot resolve a valid endpoint for pilotctl")
}

// Register the host
func (r *PilotCtl) Register() (*ctl.RegistrationResponse, error) {
	defer TRA(CE())
	i := r.host
	// set the machine id
	reg := &ctl.RegistrationRequest{
		Hostname:    i.HostName,
		MachineId:   i.HostUUID,
		OS:          i.OS,
		Platform:    fmt.Sprintf("%s, %s, %s", i.Platform, i.PlatformFamily, i.PlatformVersion),
		Virtual:     i.Virtual,
		TotalMemory: i.TotalMemory,
		CPUs:        i.CPUs,
		HostIP:      i.HostIP,
		MacAddress:  i.MacAddress,
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
	defer TRA(CE())
	// is there a result from a job ready?
	var (
		payload ctlCore.Serializable
		result  *ctl.JobResult
		events  *ctl.Events
	)
	// check if the worker has a job result to be sent to pilot control
	result, err := r.worker.Result()
	if err != nil {
		return ctl.PingResponse{}, err
	}
	if result != nil {
		// send the job result in the ping request
		payload = &ctl.PingRequest{Result: result}
	} else {
		// if we do not have any job result to post, can post event information
		// try and get up to a maximum of 5 events
		events, err = getEvents(5)
		// if there is an error retrieving events
		if err != nil {
			// return the error
			return ctl.PingResponse{}, err
		}
		// if there are events to send
		if events != nil && len(events.Events) > 0 {
			// send syslog events in the ping request
			payload = &ctl.PingRequest{Events: events}
		}
	}
	uri := fmt.Sprintf("%s/ping", r.cfg.BaseURI)
	resp, err := r.client.Post(uri, payload, r.addToken)
	if err != nil {
		return ctl.PingResponse{}, fmt.Errorf("ping failed ping: %s", err)
	}
	if resp.StatusCode > 299 {
		return ctl.PingResponse{}, fmt.Errorf("call to the remote service failed: %d - %s", resp.StatusCode, resp.Status)
	}
	// if a result was posted to control, remove it from the local cache
	if result != nil {
		err = r.worker.RemoveResult(result)
		if err != nil {
			ErrorLogger.Printf("failed to remove job result from local queue: %s\n", err)
		}
	}
	// if syslog events were posted to control, remove the marker from the local cache
	if events != nil {
		err = removeEvents()
		if err != nil {
			ErrorLogger.Printf("failed to remove events marker from local cache: %s\n", err)
		}
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

func (r *PilotCtl) addToken(req *http.Request, payload ctlCore.Serializable) error {
	defer TRA(CE())
	payload = nil
	// add an authentication token to the request
	req.Header.Set("Authorization", newToken(r.host.HostUUID, r.host.HostIP, r.host.HostName))
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (r *PilotCtl) SubmitCveReport(report []byte) error {
	var payload ctlCore.Serializable
	payload = &ctl.CveRequest{
		HostUUID: r.host.HostUUID,
		Report:   report,
	}
	uri := fmt.Sprintf("%s/cve/upload", r.cfg.BaseURI)
	resp, err := r.client.Post(uri, payload, r.addToken)
	if err != nil {
		return fmt.Errorf("cannot submit CVE report: %s", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("cannot submit CVE report: call to the remote service failed, %d - %s", resp.StatusCode, resp.Status)
	}
	return nil
}
