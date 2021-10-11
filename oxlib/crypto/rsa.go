package crypto

/*
  Onix Config Manager - Cryptographic Utilities
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

// GenerateRSAKeys generate an RSA key pair and returns their base64 encoded string for PKCS #1, ASN.1 DER form
func GenerateRSAKeys() (pk string, pub string, err error) {
	key, err := rsa.GenerateMultiPrimeKey(rand.Reader, 2, 2048)
	if err != nil {
		return "", "", err
	}
	pk = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
	pub = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&key.PublicKey))
	return
}

func Decrypt(cipherText string, privateKey string) (string, error) {
	pk, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	privKey, err := x509.ParsePKCS1PrivateKey(pk)
	if err != nil {
		return "", err
	}
	ct, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	label := []byte("OAEP Encrypted")
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ct, label)
	if err != nil {
		return "", err
	}
	return string(plaintext[:]), nil
}

func Encrypt(plainText string, publicKey string) (string, error) {
	pub, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pubKey, err := x509.ParsePKCS1PublicKey(pub)
	if err != nil {
		return "", err
	}
	label := []byte("OAEP Encrypted")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(plainText), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
