/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  This code has been based on https://github.com/AaronO/go-rsa-sign
*/
package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io/ioutil"
	"path"
)

type Signer struct {
}

func (s *Signer) Sign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func (s *Signer) SignHex(key *rsa.PrivateKey, data []byte) (string, error) {
	sig, err := s.Sign(key, data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig), nil
}

func (s *Signer) SignBase64(key *rsa.PrivateKey, data []byte) (string, error) {
	sig, err := s.Sign(key, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func LoadPrivateKey(group, name string) (*rsa.PrivateKey, error) {
	// first attempt to load the key from the registry/keys/group/name path
	private, _ := KeyNames(path.Join(core.RegistryPath(), "keys"), fmt.Sprintf("%s_%s", group, name))
	pemKey, err := ioutil.ReadFile(private)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		private, _ = KeyNames(path.Join(core.RegistryPath(), "keys"), group)
		pemKey, err = ioutil.ReadFile(private)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			private, _ = KeyNames(path.Join(core.RegistryPath(), "keys"), "root")
			pemKey, err = ioutil.ReadFile(private)
		}
	}
	return ParsePrivateKey(pemKey)
}

func LoadPublicKey(group, name string) (*rsa.PublicKey, error) {
	// first attempt to load the key from the registry/keys/group/name path
	_, public := KeyNames(path.Join(core.RegistryPath(), "keys"), fmt.Sprintf("%s_%s", group, name))
	pemKey, err := ioutil.ReadFile(public)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		_, public = KeyNames(path.Join(core.RegistryPath(), "keys"), group)
		pemKey, err = ioutil.ReadFile(public)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			_, public = KeyNames(path.Join(core.RegistryPath(), "keys"), "root")
			pemKey, err = ioutil.ReadFile(public)
		}
	}
	return ParsePublicKey(pemKey)
}
