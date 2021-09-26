package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Event a pilot event to be sent to piloctl service
type Event struct {
	// internal time id
	timeId            string
	EventID           string `json:"event_id,omitempty" yaml:"event_id,omitempty" bson:"event_id,omitempty"`
	Client            string `json:"client,omitempty" yaml:"client,omitempty" bson:"client,omitempty"`
	Hostname          string `json:"hostname,omitempty" yaml:"hostname,omitempty" bson:"hostname,omitempty"`
	MachineId         string `json:"machine_id" yaml:"machine_id" bson:"machine_id"`
	HostAddress       string `json:"host_address,omitempty" yaml:"host_address,omitempty" bson:"host_address,omitempty"`
	Organisation      string `json:"organisation,omitempty" yaml:"organisation,omitempty" bson:"organisation,omitempty"`
	OrganisationGroup string `json:"organisation_group,omitempty" yaml:"organisation_group,omitempty" bson:"organisation_group,omitempty"`
	Area              string `json:"area,omitempty" yaml:"area,omitempty" bson:"area,omitempty"`
	Location          string `json:"location,omitempty" yaml:"location,omitempty" bson:"location,omitempty"`
	Facility          int    `json:"facility,omitempty" yaml:"facility,omitempty" bson:"facility,omitempty"`
	Priority          int    `json:"priority,omitempty" yaml:"priority,omitempty" bson:"priority,omitempty"`
	Severity          int    `json:"severity,omitempty" yaml:"severity,omitempty" bson:"severity,omitempty"`
	Time              string `json:"time,omitempty" yaml:"time,omitempty" bson:"time,omitempty"`
	TLSPeer           string `json:"tls_peer,omitempty" yaml:"tls_peer,omitempty" bson:"tls_peer,omitempty"`
	BootTime          string `json:"boot_time,omitempty" yaml:"boot_time,omitempty" bson:"boot_time,omitempty"`
	Content           string `json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
	Tag               string `json:"tag,omitempty" yaml:"tag,omitempty" bson:"tag,omitempty"`
}

// NewEvent create a new serializable event from a syslog entry in RFC 3164
func NewEvent(logPart format.LogParts, info HostInfo) (*Event, error) {
	tId := timeBasedId()
	entry := &Event{timeId: tId}
	entry.Priority = logPart["priority"].(int)
	entry.Severity = logPart["severity"].(int)
	entry.Hostname = logPart["hostname"].(string)
	entry.Client = logPart["client"].(string)
	entry.TLSPeer = logPart["tls_peer"].(string)
	entry.Facility = logPart["facility"].(int)
	entry.Time = logPart["timestamp"].(time.Time).Format(TimeLayout)
	entry.EventID = fmt.Sprintf("%s:%s", info.MachineId, tId)
	entry.BootTime = info.BootTime
	entry.MachineId = info.MachineId
	entry.HostAddress = info.HostIP
	entry.Content = logPart["content"].(string)
	entry.Tag = logPart["tag"].(string)
	return entry, nil
}

// Save the event to the file system
// path: is the folder where events will be saved
func (e *Event) Save(path string) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	filename := filepath.Join(path, fmt.Sprintf("%s.ev", timeBasedId()))
	return ioutil.WriteFile(filename, bytes, os.ModePerm)
}

// timeBasedId generates a time based event timeBasedId
func timeBasedId() string {
	t := time.Now()
	return fmt.Sprintf("%02d%02d%02s%02d%02d%02d%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[2:], t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:5])
}

type Events struct {
	Events []Event `json:"events"`
}

func (r *Events) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

func (r *Events) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}
