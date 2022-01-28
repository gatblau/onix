/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ctl "github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AKToken struct {
	info       userKeyInfo
	Username   string    `json:"username"`
	MacAddress string    `json:"mac_address"`
	IpAddress  string    `json:"ip_address"`
	Hostname   string    `json:"hostname"`
	Time       time.Time `json:"time"`
}

func NewAKToken(clientInfo userKeyInfo, hostInfo *ctl.HostInfo) AKToken {
	return AKToken{
		info:       clientInfo,
		Username:   clientInfo.Username,
		MacAddress: hostInfo.PrimaryMAC,
		IpAddress:  hostInfo.HostIP,
		Hostname:   hostInfo.HostName,
		Time:       time.Now(),
	}
}

func (t AKToken) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		ErrorLogger.Printf("cannot create activation key request bearer token: %s\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("Bearer %s %s", t.Username, encrypt(t.info.SK, hex.EncodeToString(b), t.info.IV))
}

func activate(info *ctl.HostInfo) {
	var (
		failures float64 = 0
		interval time.Duration
	)
	// first check for a valid activation key
	if !AkExist() {
		// if no user key exists
		if !UserKeyExist() {
			// cannot continue
			ErrorLogger.Printf("cannot launch pilot, missing activation key\n")
			os.Exit(1)
		}
		// otherwise, it can start the activation process
		InfoLogger.Printf("cannot find activation key, initiating activation protocol\n")
		uKey, err := loadUserKey(UserKeyFile())
		if err != nil {
			// cannot continue
			ErrorLogger.Printf("cannot launch pilot, cannot load user key: %s\n", err)
			os.Exit(1)
		}
		tenant, err := readUserKey(*uKey)
		if err != nil {
			ErrorLogger.Printf("cannot launch pilot, cannot load user key: %s\n", err)
			os.Exit(1)
		}
		// fetch remote key
		fetched, err := requestAKey(*tenant, info)
		// if failed retry
		for !fetched {
			// calculates wait interval with exponential backoff and jitter
			interval = nextInterval(failures)
			ErrorLogger.Printf("cannot retrieve activation key, retrying in %.2f minutes: %s\n", interval.Seconds()/60, err)
			// wait interval
			time.Sleep(interval)
			failures++
			// try again
			fetched, err = requestAKey(*tenant, info)
			// if successful
			if err == nil {
				// break the loop
				break
			}
		}
		InfoLogger.Printf("activation key deployed, pilot is ready to launch\n")
	}
	// before doing anything, verify activation key
	ak, err := loadAKey(AkFile())
	if err != nil {
		// if it cannot load activation key exit
		ErrorLogger.Printf("cannot launch pilot: cannot load activation key, %s\n", err)
		os.Exit(1)
	}
	akInfo, err := readAKey(*ak)
	if err != nil {
		// if it cannot load activation key exit
		ErrorLogger.Printf("cannot launch pilot: cannot read activation key, %s\n", err)
		os.Exit(1)
	}
	// set the activation
	A = akInfo
	// check expiration date
	if A.Expiry.Before(time.Now()) {
		// if activation expired the exit
		ErrorLogger.Printf("cannot launch pilot: activation key expired\n")
		os.Exit(1)
	}
	// check if the mac-address is valid
	validMac := false
	for _, address := range info.MacAddress {
		if address == A.MacAddress {
			validMac = true
			break
		}
	}
	if !validMac {
		// if activation expired the exit
		ErrorLogger.Printf("cannot launch pilot: invalid mac address\n")
		os.Exit(1)
	}
	// set host UUID
	info.HostUUID = A.HostUUID
}

func requestAKey(clientKey userKeyInfo, info *ctl.HostInfo) (bool, error) {
	bearerToken := NewAKToken(clientKey, info)
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// TODO: set to false if in production!!!
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/activation-key", clientKey.URI), nil)
	if err != nil {
		return false, fmt.Errorf("cannot create activation request: %s\n", err)
	}
	req.Header.Add("Authorization", bearerToken.String())
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return false, fmt.Errorf("cannot send http request: %s\n", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("activation http request failed with code %d: %s\n", resp.StatusCode, resp.Status)
	}
	ak, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("cannot read activation key from http response: %s\n", err)
	}
	err = os.WriteFile(AkFile(), ak, 0600)
	if err != nil {
		return false, fmt.Errorf("cannot write activation file: %s\n", err)
	}
	return true, nil
}
