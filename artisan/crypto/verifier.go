/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  This code has been based on https://github.com/AaronO/go-rsa-sign
*/
package crypto

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
)

type Verifier struct {
}

func (v *Verifier) Verify(key *rsa.PublicKey, data, sig []byte) error {
	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}

func (v *Verifier) VerifyHex(key *rsa.PublicKey, data []byte, sigHex string) error {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return err
	}
	return v.Verify(key, data, sig)
}

func (v *Verifier) VerifyBase64(key *rsa.PublicKey, data []byte, sig64 string) error {
	sig, err := base64.StdEncoding.DecodeString(sig64)
	if err != nil {
		return err
	}
	return v.Verify(key, data, sig)
}
