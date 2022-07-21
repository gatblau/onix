/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

// Host monitoring information
type Host struct {
	HostUUID       string   `json:"host_uuid"`
	HostMacAddress string   `json:"host_mac_address"`
	OrgGroup       string   `json:"org_group"`
	Org            string   `json:"org"`
	Area           string   `json:"area"`
	Location       string   `json:"location"`
	Connected      bool     `json:"connected"`
	LastSeen       int64    `json:"last_seen"`
	Since          int      `json:"since"`
	SinceType      string   `json:"since_type"`
	Label          []string `json:"label"`
	Critical       int      `json:"critical"`
	High           int      `json:"high"`
	Medium         int      `json:"medium"`
	Low            int      `json:"low"`
}
