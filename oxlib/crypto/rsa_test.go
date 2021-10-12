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
	"os"
	"testing"
)

// test encryption / decryption
func TestEncryptDecrypt(t *testing.T) {
	pk, pub, err := GenerateRSAKeysPKCS1()
	if err != nil {
		t.Error(err)
	}
	message := "Hello World!"
	cipherText, err := EncryptPKCS1(message, pub)
	if err != nil {
		t.Error(err)
	}
	plainText, err := DecryptPKCS1(cipherText, pk)
	if err != nil {
		t.Error(err)
	}
	if message != plainText {
		t.Errorf("mismatch %s vs %s", message, plainText)
	}
}

// prints 2 RSA key-pairs
func TestPrintKeys(t *testing.T) {
	pk, pub, err := GenerateRSAKeysPKCS1()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("PK-1 => '%s'\n", pk)
	fmt.Printf("PUB-1 => '%s'\n\n", pub)
	pk, pub, err = GenerateRSAKeysPKCS1()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("PK-2 => '%s'\n", pk)
	fmt.Printf("PUB-2 => '%s'\n", pub)
}

func TestSaveKeys(t *testing.T) {
	pk, pub, err := GenerateRSAKeysPKCS1()
	if err != nil {
		t.Error(err)
	}
	// user:   read/write
	// group:  none
	// others: none
	var perm os.FileMode = 0600
	os.WriteFile("id_rsa", []byte(pk), perm)
	os.WriteFile("id_rsa.pub", []byte(pub), perm)
}
