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
	"strconv"
	"strings"
	"time"
)

type userKey struct {
	Key       userKeyString `json:"uk"`
	Signature string        `json:"s"`
}

type userKeyString string

type userKeyInfo struct {
	Username string
	URI      string
	IV       []byte
	SK       []byte
	Expiry   *time.Time
}

func loadUserKey(path string) (*userKey, error) {
	if len(path) == 0 {
		path = ".userkey"
	}
	path = Abs(path)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read user key file: %s\n", err)
	}
	d, err := hex.DecodeString(string(b[:]))
	if err != nil {
		return nil, fmt.Errorf("cannot decode user key: %s\n", err)
	}
	key := new(userKey)
	err = json.Unmarshal(d, key)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal user key: %s\n", err)
	}
	return key, nil
}

// readUserKey read the content of an encoded user key and verifies its digital signature
func readUserKey(key userKey) (*userKeyInfo, error) {
	// check the validity of the key's digital signature
	if valid, err := verify(string(key.Key), key.Signature); !valid {
		return nil, fmt.Errorf("invalid user key signature: %s\n", err)
	}
	// decrypt the key information
	d, err := decrypt(sk, string(key.Key), iv)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt user key data: %s\n", err)
	}
	db, err := hex.DecodeString(d[:])
	if err != nil {
		return nil, fmt.Errorf("cannot decode user key data: %s\n", err)
	}
	parts := strings.Split(string(db[44:]), ",")
	var expiry *time.Time
	// parse the expiry days
	days, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve expiry days from user key: %s\n", err)
	}
	// if days is not zero then create a date stamp
	if days > 0 {
		expiryDate := time.Now().Add(time.Hour * 24 * time.Duration(days))
		expiry = &expiryDate
	}
	return &userKeyInfo{
		Username: parts[0],
		URI:      parts[1],
		IV:       db[:12],
		SK:       db[12:44],
		Expiry:   expiry,
	}, nil
}
