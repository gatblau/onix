/*
  Onix Config Manager - crypto utils
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func EncryptAES(plaintext string) string {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	c, err := aes.NewCipher(key)
	CheckError(err)
	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))
	key = reverse(key)
	out = bytes.Join([][]byte{key[:15], out, key[15:]}, []byte(""))
	return hex.EncodeToString(out)
}

func DecryptAES(encryptedText string) string {
	ciphertext, _ := hex.DecodeString(encryptedText)
	key := bytes.Join([][]byte{ciphertext[:15], ciphertext[len(ciphertext)-17:]}, []byte(""))
	key = reverse(key)
	ciphertext = ciphertext[15 : len(ciphertext)-17]
	c, err := aes.NewCipher(key)
	CheckError(err)
	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)
	return string(pt[:])
}

func reverse(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return append(reverse(input[1:]), input[0])
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
