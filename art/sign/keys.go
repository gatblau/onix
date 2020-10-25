/*
  Onix Config Manager - Art
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package sign

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// generate a new RSA key pair
func NewKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)
	return key
}

// generate a new RSA Multi-Prime key pair
// more efficient key-generation and decryption/signing operation
// however, it might be easier to factor a multi-prime RSA public key than a dual-prime one
func NewMultiPrimeKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateMultiPrimeKey(reader, 16, bitSize)
	checkError(err)
	return key
}

// save the key to a pem file
func SavePrivateKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()
	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

// save the public key to a pem file
func SavePublicKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	checkError(err)
	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()
	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
