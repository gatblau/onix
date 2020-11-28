/*
  Onix Config Manager - Artie
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
	"github.com/gatblau/onix/artie/core"
	"os"
	"path/filepath"
)

// generate a new RSA key pair
func NewKeyPair(size int) *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := size
	key, err := rsa.GenerateKey(reader, bitSize)
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

// works out the fully qualified names of the private and public RSA keys
func KeyNames(path, name string) (key string, pub string) {
	if len(path) == 0 {
		path = "."
	}
	if len(name) == 0 {
		name = "id"
	}
	// if the path is relative then make it absolute
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		core.CheckErr(err, "cannot return an absolute representation of path: '%s'", path)
		path = p
	}
	// works out the private key name
	keyName := filepath.Join(path, fmt.Sprintf("%s_rsa_key.pem", name))
	// works out the public key name
	pubName := filepath.Join(path, fmt.Sprintf("%s_rsa_pub.pem", name))
	return keyName, pubName
}

// generates a private and public RSA keys for signing and verifying artefacts
func GenerateKeys(path, name string, size int) {
	if size > 4500 {
		core.RaiseErr("maximum bit size 4500 exceeded")
	}
	if size == 0 {
		size = 2048
	}
	keyFilename, pubFilename := KeyNames(path, name)
	key := NewKeyPair(size)
	SavePrivateKey(keyFilename, key)
	SavePublicKey(pubFilename, key.PublicKey)
}
