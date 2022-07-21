/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "bytes"

type CveRequest struct {
	HostUUID string `json:"host_uuid"`
	Report   []byte `json:"report"`
}

func (r *CveRequest) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

func (r *CveRequest) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}
