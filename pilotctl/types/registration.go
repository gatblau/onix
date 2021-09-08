package types

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import "bytes"

// RegistrationRequest information sent by pilot upon host registration
type RegistrationRequest struct {
	Hostname    string  `json:"hostname"`
	HostIP      string  `json:"host_ip"`
	MachineId   string  `json:"machine_id"`
	OS          string  `json:"os"`
	Platform    string  `json:"platform"`
	Virtual     bool    `json:"virtual"`
	TotalMemory float64 `json:"total_memory"`
	CPUs        int     `json:"cpus"`
}

// Reader Get a JSON bytes reader for the Serializable
func (r *RegistrationRequest) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Bytes Get a []byte representing the Serializable
func (r *RegistrationRequest) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}

// RegistrationResponse data returned to pilot upon registration
type RegistrationResponse struct {
	// the status of the registration - I: created, U: updated, N: already exist
	Operation string `json:"operation"`
}
