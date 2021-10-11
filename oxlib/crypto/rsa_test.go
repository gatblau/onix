package crypto

/*
  Onix Config Manager - Cryptographic Utilities
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"fmt"
	"testing"
)

// test encryption / decryption
func TestCrypto(t *testing.T) {
	pk, pub, err := GenerateRSAKeys()
	if err != nil {
		t.Error(err)
	}
	message := "Hello World!"
	cipherText, err := Encrypt(message, pub)
	if err != nil {
		t.Error(err)
	}
	plainText, err := Decrypt(cipherText, pk)
	if err != nil {
		t.Error(err)
	}
	if message != plainText {
		t.Errorf("mismatch %s vs %s", message, plainText)
	}
}

// prints 2 RSA key-pairs
func TestPrintKeys(t *testing.T) {
	pk, pub, err := GenerateRSAKeys()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("PK-1 => '%s'\n", pk)
	fmt.Printf("PUB-1 => '%s'\n\n", pub)
	pk, pub, err = GenerateRSAKeys()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("PK-2 => '%s'\n", pk)
	fmt.Printf("PUB-2 => '%s'\n", pub)
}
