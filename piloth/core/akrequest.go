/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
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
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AKRequestBearerToken struct {
	info     userKeyInfo
	Username string `json:"username"`
	// DeviceId is either the host primary interface MAC address or the device hardware uuid
	DeviceId  string    `json:"device_id"`
	IpAddress string    `json:"ip_address"`
	Hostname  string    `json:"hostname"`
	Time      time.Time `json:"time"`
}

func NewAKRequestBearerToken(clientInfo userKeyInfo, options PilotOptions) AKRequestBearerToken {
	defer TRA(CE())
	var deviceId string
	// if hardware id should be used to identify the device
	if options.UseHwId {
		// then set the device identifier to the hardware id
		deviceId = options.Info.HardwareId
	} else {
		// otherwise, set it to the primary mac address
		deviceId = options.Info.PrimaryMAC
	}
	return AKRequestBearerToken{
		info:      clientInfo,
		Username:  clientInfo.Username,
		DeviceId:  deviceId,
		IpAddress: options.Info.HostIP,
		Hostname:  options.Info.HostName,
		Time:      time.Now(),
	}
}

func (t AKRequestBearerToken) String() string {
	defer TRA(CE())
	b, err := json.Marshal(t)
	if err != nil {
		ErrorLogger.Printf("cannot create activation key request bearer token: %s\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("Bearer %s %s", t.Username, encrypt(t.info.SK, hex.EncodeToString(b), t.info.IV))
}

func activate(options PilotOptions) {
	defer TRA(CE())
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
		fetched, err := requestAKey(*tenant, options)
		// if failed retry
		for !fetched {
			// calculates wait interval with exponential backoff and jitter
			interval = nextInterval(failures)
			ErrorLogger.Printf("cannot retrieve activation key, retrying in %.2f minutes: %s\n", interval.Seconds()/60, err)
			// wait interval
			time.Sleep(interval)
			failures++
			// try again
			fetched, err = requestAKey(*tenant, options)
			// if successful
			if err == nil {
				// break the loop
				break
			}
		}
		InfoLogger.Printf("activation key deployed, pilot is ready to launch\n")
	}
	// before doing anything, verify activation key
	akInfo, err := LoadActivationKey()
	if err != nil {
		fmt.Errorf("cannot start pilot: %s", err)
		os.Exit(1)
	}
	// set the activation
	A = akInfo
	// validate the activation key
	A.Validate()
	// check expiration date
	if A.Expiry.Before(time.Now()) {
		// if activation expired the exit
		ErrorLogger.Printf("cannot launch pilot: activation key expired\n")
		os.Exit(1)
	}
	// if set to use hardware id for device identification
	if options.UseHwId {
		if A.DeviceId != options.Info.HardwareId {
			// if the device Id is not the hardware id; then exit
			ErrorLogger.Printf("cannot launch pilot: invalid host hardware id: %s\n", options.Info.HardwareId)
			os.Exit(1)
		}
	} else { // use mac address for device identification
		// check if the mac-address matches the device id in the activation key
		matchedMac := false
		for _, macAddress := range options.Info.MacAddress {
			if A.DeviceId == macAddress {
				matchedMac = true
				break
			}
		}
		// if the mac address does not match
		if !matchedMac {
			// if the device Id is not the hardware id; then exit
			ErrorLogger.Printf("cannot launch pilot: invalid host mac address: %s\n", options.Info.PrimaryMAC)
			os.Exit(1)
		}
	}
	// set host UUID
	options.Info.HostUUID = A.HostUUID
}

func LoadActivationKey() (*AKInfo, error) {
	ak, err := loadAKey(AkFile())
	if err != nil {
		// if it cannot load activation key exit
		return nil, fmt.Errorf("cannot load activation key, %s\n", err)
	}
	akInfo, err := readAKey(*ak)
	if err != nil {
		// if it cannot load activation key exit
		return nil, fmt.Errorf("cannot read activation key, %s\n", err)
	}
	return akInfo, nil
}

func requestAKey(clientKey userKeyInfo, options PilotOptions) (bool, error) {
	bearerToken := NewAKRequestBearerToken(clientKey, options)
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.InsecureSkipVerify,
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
