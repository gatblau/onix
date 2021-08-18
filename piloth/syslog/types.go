package syslog

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

type RsyslogLogRFC3164 struct {
	Client    string `json:"client" bson:"client"`
	Content   string `json:"content" bson:"content"`
	Facility  int    `json:"facility" bson:"facility"`
	Hostname  string `json:"hostname" bson:"hostname"`
	Priority  int    `json:"priority" bson:"priority"`
	Severity  int    `json:"severity" bson:"severity"`
	Tag       string `json:"tag" bson:"tag"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
	TLSPeer   string `json:"tls_peer" bson:"tls_peer"`
}

type EventLog struct {
	EventID         string `json:"event_id,omitempty" yaml:"event_id,omitempty" bson:"event_id,omitempty"`
	Client          string `json:"client,omitempty" yaml:"client,omitempty" bson:"client,omitempty"`
	CreateTimeStamp string `json:"create_time_stamp,omitempty" yaml:"create_time_stamp,omitempty" bson:"create_time_stamp,omitempty"`
	Hostname        string `json:"hostname,omitempty" yaml:"hostname,omitempty" bson:"hostname,omitempty"`
	HostID          string `json:"host_id" yaml:"host_id" bson:"host_id"`
	HostAddress     string `json:"host_address,omitempty" yaml:"host_address,omitempty" bson:"host_address,omitempty"`
	Location        string `json:"location" yaml:"location" bson:"location"`
	Facility        int    `json:"facility,omitempty" yaml:"facility,omitempty" bson:"facility,omitempty"`
	Priority        int    `json:"priority,omitempty" yaml:"priority,omitempty" bson:"priority,omitempty"`
	Severity        int    `json:"severity,omitempty" yaml:"severity,omitempty" bson:"severity,omitempty"`
	Tag             string `json:"tag" yaml:"tag" bson:"tag"`
	EventTimestamp  string `json:"event_timestamp,omitempty" yaml:"event_timestamp,omitempty" bson:"event_timestamp,omitempty"`
	Content         string `json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
	Details         string `json:"details" yaml:"details" bson:"details"`
}
