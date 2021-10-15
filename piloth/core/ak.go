/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"crypto/aes"
	"crypto/cipher"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/helper"
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
	_, err := os.Stat(ConfFile())
	return err == nil
}

func LoadAK() (*AK, error) {
	akBytes, err := os.ReadFile(ConfFile())
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
