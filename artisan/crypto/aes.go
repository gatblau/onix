package crypto

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  NOTICE: original code borrowed from https://github.com/kashifsoofi/crypto-sandbox

*/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type CipherMode int

const (
	CBC CipherMode = iota
	GCM
)

type Padding int

const (
	NoPadding Padding = iota
	PKCS7
)

type AesCrypto struct {
	CipherMode CipherMode
	Padding    Padding
}

const AesIvSize = 16

func (c AesCrypto) Encrypt(plainText string, key []byte) (string, error) {
	// create a new aes cipher using key
	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if c.CipherMode == GCM {
		return c.EncryptGcm(aes, plainText)
	} else {
		return c.EncryptCbc(aes, plainText)
	}
}

func (c AesCrypto) EncryptGcm(aes cipher.Block, plainText string) (string, error) {
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	cipherText := gcm.Seal(nil, nonce, plainTextBytes, nil)

	return c.PackCipherData(cipherText, nonce, gcm.Overhead()), nil
}

func (c AesCrypto) EncryptCbc(aes cipher.Block, plainText string) (string, error) {
	iv := make([]byte, AesIvSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(aes, iv)

	plainTextBytes := []byte(plainText)
	plainTextBytes, err := pkcs7Pad(plainTextBytes, encrypter.BlockSize())
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plainTextBytes))
	encrypter.CryptBlocks(cipherText, plainTextBytes)

	return c.PackCipherData(cipherText, iv, 0), nil
}

func (c AesCrypto) Decrypt(cipherText string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	encryptedBytes, iv, tagSize := c.UnpackCipherData(data)

	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if c.CipherMode == GCM {
		return DecryptGcm(aes, encryptedBytes, iv, tagSize)
	} else {
		return DecryptCbc(aes, encryptedBytes, iv)
	}
}

func DecryptGcm(aes cipher.Block, encrypted []byte, nonce []byte, tagSize int) (string, error) {
	var aesgcm cipher.AEAD
	var err error
	aesgcm, err = cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes[:len(decryptedBytes)]), nil
}

func DecryptCbc(aes cipher.Block, encrypted []byte, iv []byte) (string, error) {
	decryptor := cipher.NewCBCDecrypter(aes, iv)

	decryptedBytes := make([]byte, len(encrypted))
	decryptor.CryptBlocks(decryptedBytes, encrypted)

	decryptedBytes, err := pkcs7Unpad(decryptedBytes, decryptor.BlockSize())
	if err != nil {
		return "", err
	}

	return string(decryptedBytes[:len(decryptedBytes)]), nil
}

func (c AesCrypto) PackCipherData(cipherText []byte, iv []byte, tagSize int) string {
	ivLength := len(iv)
	dataLength := len(cipherText) + ivLength + 1
	if c.CipherMode == GCM {
		dataLength += 1
	}

	data := make([]byte, dataLength)

	// set first 2 bytes as nonceSize, to make cipher data compatible with crypto methods in other languages in this repo

	data[0] = byte(ivLength)
	index := 1
	if c.CipherMode == GCM {
		data[1] = byte(tagSize)
		index += 1
	}
	copy(data[index:], iv[0:ivLength])
	index += ivLength
	copy(data[index:], cipherText)

	return base64.StdEncoding.EncodeToString(data)
}

func (c AesCrypto) UnpackCipherData(data []byte) ([]byte, []byte, int) {
	ivSize := int(data[0])
	index := 1
	tagSize := 0
	if c.CipherMode == GCM {
		tagSize = int(data[index])
		index += 1
	}
	iv, encryptedBytes := data[index:index+ivSize], data[index+ivSize:]

	return encryptedBytes, iv, tagSize
}

// ref: https://golang-examples.tumblr.com/post/98350728789/pkcs7-padding
// Appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("Invalid block length %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...), nil
}

// Returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("Invalid block length %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("Invalid data length %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("Invalid padding")
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("Invalid padding")
		}
	}

	return data[:len(data)-padlen], nil
}
