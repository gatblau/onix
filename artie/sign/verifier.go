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
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
)

type Verifier struct {
	Key *rsa.PublicKey
}

func NewVerifier(pemKey []byte) (*Verifier, error) {
	key, err := parsePublicKey(pemKey)
	if err != nil {
		return nil, err
	}
	return &Verifier{key}, nil
}

func (v *Verifier) Verify(data, sig []byte) error {
	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(v.Key, hash, hashed, sig)
}

func (v *Verifier) VerifyHex(data []byte, sigHex string) error {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return err
	}
	return v.Verify(data, sig)
}

func (v *Verifier) VerifyBase64(data []byte, sig64 string) error {
	sig, err := base64.StdEncoding.DecodeString(sig64)
	if err != nil {
		return err
	}
	return v.Verify(data, sig)
}
