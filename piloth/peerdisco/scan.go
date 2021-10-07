package peerdisco

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
	"github.com/gatblau/onix/pilotctl/types"
	"github.com/gatblau/onix/piloth/core"
	"github.com/schollz/peerdiscovery"
	"hash/fnv"
	"time"
)

type ScanJob struct {
	info types.HostInfo
}

func NewScanJob(info types.HostInfo) ScanJob {
	return ScanJob{
		info: info,
	}
}

// Execute Called by the Scheduler when a Trigger fires that is associated with the Job.
func (s ScanJob) Execute() {
	// discover peers
	links, err := discover(s.info)
	if err != nil {
		core.ErrorLogger.Printf("failed to discover pilot peers: %s\n", err)
		return
	}
	existingLinks, err := readLinks()
	if err != nil {
		core.ErrorLogger.Printf("failed to retrieve links list: %s\n", err)
		return
	}
	if !links.Equals(existingLinks) {
		// TODO: raise event
		core.InfoLogger.Printf("Peer network change detected\n")
		// TODO: update saved link list
	}
}

// Description returns a Job description.
func (s ScanJob) Description() string {
	return "Scan the network searching for Pilot instances"
}

// Key returns a Job unique key.
func (s ScanJob) Key() int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s.Description()))
	return int(h.Sum32())
}

func discover(info types.HostInfo) (links Links, err error) {
	// discover peers
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		// unlimited peers to discover
		Limit: -1,
		// broadcast host information
		Payload: linkInfo(info),
		// time between broadcasts
		Delay: 500 * time.Millisecond,
		// time spend discovering
		TimeLimit: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to discover pilot instances: %s\n", err)
	}
	for _, d := range discoveries {
		var link *Link
		err = json.Unmarshal(d.Payload, link)
		if err != nil {
			core.ErrorLogger.Printf("cannot unmarshal peer payload: %s\n", err)
		}
		link.Address = d.Address
		links = append(links, *link)
	}
	return links, nil
}

func linkInfo(info types.HostInfo) []byte {
	infoBytes, err := json.Marshal(NewLink(info))
	if err != nil {
		core.ErrorLogger.Printf("cannot marshal host info: %s\n", err)
	}
	return infoBytes
}

func readLinks() ([]Link, error) {
	panic("not implemented")
}
