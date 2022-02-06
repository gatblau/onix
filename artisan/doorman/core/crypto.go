/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

var (
	sk = "819da5fa6489428f9b95780d5f5d740d651b50e21c99b33101eceeb37a5c8850"
	iv = "f9e02e9fce1a498a34b08e67"
)

func skBytes() []byte {
	d, _ := hex.DecodeString(sk)
	return d
}

func ivBytes() []byte {
	d, _ := hex.DecodeString(iv)
	return d
}

func decrypt(cypherText string) (string, error) {
	keyBytes, _ := hex.DecodeString(sk)
	ciphertext, _ := hex.DecodeString(cypherText)
	iVecBytes, _ := hex.DecodeString(iv)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := aesgcm.Open(nil, iVecBytes, ciphertext, nil)
	if err != nil {
		return "", err
	}
	s := string(plaintext[:])
	return s, nil
}

func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(skBytes())
	if err != nil {
		return "", fmt.Errorf("cannot create cipher block: %s", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cannot wrap block cipher: %s", err)
	}
	ciphertext := aesgcm.Seal(nil, ivBytes(), []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}
