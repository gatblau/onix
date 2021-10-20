/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	c "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	ctl "github.com/gatblau/onix/pilotctl/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AK struct {
	HostUUID   string    `json:"host_uuid"`
	MacAddress string    `json:"mac_address"`
	CtlURI     string    `json:"ctl_uri"`
	Expiry     time.Time `json:"expiry"`
	VerifyKey  string    `json:"verify_key"`
}

func AkExist() bool {
	_, err := os.Stat(AkFile())
	return err == nil
}

func LoadAK() (*AK, error) {
	akBytes, err := os.ReadFile(AkFile())
	if err != nil {
		return nil, fmt.Errorf("cannot read activation key file: %s\n", err)
	}
	content, err := helper.DecryptMessageArmored(decrypt(k, a, i), nil, string(akBytes[:]))
	if err != nil {
		return nil, fmt.Errorf("invalid activation key: %s\n", err)
	}
	ak := new(AK)
	err = json.Unmarshal([]byte(content), ak)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal activation key: %s\n", err)
	}
	return ak, nil
}

func decrypt(key string, ct string, iv string) string {
	keyBytes, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(ct)
	ivBytes, _ := hex.DecodeString(iv)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, ivBytes, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	s := string(plaintext[:])
	return s
}

type AKRequest struct {
	Data      AKRequestEnvelope `json:"data"`
	Signature string            `json:"signature"`
}

type AKRequestEnvelope struct {
	MacAddress string    `json:"mac_address"`
	IpAddress  string    `json:"ip_address"`
	Hostname   string    `json:"hostname"`
	Time       time.Time `json:"time"`
}

func NewAKRequest(data AKRequestEnvelope) (*AKRequest, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("cannot serialise request envelope: %s\n", err)
	}
	pk, err := c.NewKeyFromArmored(decrypt(k, a, i))
	if err != nil {
		return nil, fmt.Errorf("cannot load signing key: %s\n", err)
	}
	kr, err := c.NewKeyRing(pk)
	sig, err := kr.SignDetached(c.NewPlainMessageFromString(string(dataBytes[:])))
	if err != nil {
		return nil, fmt.Errorf("cannot sign request envelope: %s\n", err)
	}
	s, err := sig.GetArmored()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve request envelope armored signature: %s\n", err)
	}
	return &AKRequest{
		Data:      data,
		Signature: s,
	}, nil
}

func activate(info *ctl.HostInfo) {
	var (
		failures float64 = 0
		interval time.Duration
	)
	// first check for a valid activation key
	if !AkExist() {
		InfoLogger.Printf("cannot find activation key, initiating activation protocol\n")
		// fetch remote key
		fetched, err := fetchToken(info)
		// if failed retry
		for !fetched {
			// calculates wait interval with exponential backoff and jitter
			interval = nextInterval(failures)
			ErrorLogger.Printf("cannot retrieve activation key, retrying in %.2f minutes: %s\n", interval.Seconds()/60, err)
			// wait interval
			time.Sleep(interval)
			failures++
			// try again
			fetched, err = fetchToken(info)
			// if successful
			if err == nil {
				// break the loop
				break
			}
		}
	}
	// before doing anything, verify activation key
	ak, err := LoadAK()
	if err != nil {
		// if it cannot load activation key exit
		ErrorLogger.Printf("cannot launch pilot: cannot load activation key, %s\n", err)
		os.Exit(1)
	}
	// set the activation
	A = ak
	// check expiration date
	if A.Expiry.Before(time.Now()) {
		// if activation expired the exit
		ErrorLogger.Printf("cannot launch pilot: activation key expired\n")
		os.Exit(1)
	}
	// check if the mac-adress is valid
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

func fetchToken(info *ctl.HostInfo) (bool, error) {
	akreq, err := NewAKRequest(
		AKRequestEnvelope{
			MacAddress: info.MacAddress[0],
			IpAddress:  info.HostIP,
			Hostname:   info.HostName,
			Time:       time.Now(),
		})
	if err != nil {
		return false, fmt.Errorf("cannot sign activation request: %s\n", err)
	}
	body, err := json.Marshal(akreq)
	if err != nil {
		return false, fmt.Errorf("cannot marshal activation request: %s\n", err)
	}
	cf := new(Config)
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// TODO: set to false if in production!!!
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("POST", cf.getActivationURI(), io.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return false, fmt.Errorf("cannot create activation request: %s\n", err)
	}
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
