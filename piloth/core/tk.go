/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type tenantKey struct {
	Key       tenantKeyString `json:"uk"`
	Signature string          `json:"s"`
}

type tenantKeyString string

type tenantKeyInfo struct {
	Tenant string
	URI    string
	IV     []byte
	SK     []byte
}

func loadTenantKey(path string) (*tenantKey, error) {
	if len(path) == 0 {
		path = ".tenant"
	}
	path = Abs(path)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read tenant key file: %s\n", err)
	}
	d, err := hex.DecodeString(string(b[:]))
	if err != nil {
		return nil, fmt.Errorf("cannot decode tenant key: %s\n", err)
	}
	key := new(tenantKey)
	err = json.Unmarshal(d, key)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal tenant key: %s\n", err)
	}
	return key, nil
}

// ReadTenantKey read the content of an encoded tenant key and verifies its digital signature
func readTenantKey(key tenantKey) (*tenantKeyInfo, error) {
	// check the validity of the key's digital signature
	if valid, err := verify(string(key.Key), key.Signature); !valid {
		return nil, fmt.Errorf("invalid tenant key signature: %s\n", err)
	}
	// decrypt the key information
	d, err := decrypt(sk, string(key.Key), iv)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt tenant key data: %s\n", err)
	}
	db, err := hex.DecodeString(d[:])
	if err != nil {
		return nil, fmt.Errorf("cannot decode tenant key data: %s\n", err)
	}
	parts := strings.Split(string(db[44:]), ",")
	return &tenantKeyInfo{
		Tenant: parts[0],
		URI:    parts[1],
		IV:     db[:12],
		SK:     db[12:44],
	}, nil
}
