package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	c "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/gatblau/onix/artisan/crypto"
	"os"
	"path/filepath"
)

func decrypt(key string, cypherText string, iv string) (string, error) {
	keyBytes, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(cypherText)
	ivBytes, _ := hex.DecodeString(iv)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := aesgcm.Open(nil, ivBytes, ciphertext, nil)
	if err != nil {
		return "", err
	}
	s := string(plaintext[:])
	return s, nil
}

func encrypt(key []byte, plaintext string, iv []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext)
}

func verify(text, signature string) (bool, error) {
	msg := c.NewPlainMessageFromString(text)
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("cannot decode PGP signature: %s\n", err)
	}
	pgpSig := c.NewPGPSignature(sigBytes)
	if err != nil {
		return false, fmt.Errorf("cannot read PGP signature: %s\n", err)
	}
	d, err := decrypt(sk, pub, iv)
	if err != nil {
		return false, fmt.Errorf("cannot decrypt public PGP key: %s\n", err)
	}
	pub, err := c.NewKeyFromArmored(d)
	if err != nil {
		return false, fmt.Errorf("cannot read public PGP key: %s\n", err)
	}
	signKR, err := c.NewKeyRing(pub)
	err = signKR.VerifyDetached(msg, pgpSig, c.GetUnixTime())
	return err == nil, err
}

// sign create a cryptographic signature for the passed-in object
func verify2(obj interface{}, signature string) error {
	// decode the  signature
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("verify => cannot decode signature string '%s': %s\n", signature, err)
	}
	// obtain the object checksum
	sum, err := checksum(obj)
	if err != nil {
		return fmt.Errorf("verify => cannot calculate checksum: %s\n", err)
	}
	// load verification key from activation key
	pgp, err := crypto.LoadPGPBytes([]byte(A.VerifyKey))
	if err != nil {
		return fmt.Errorf("verify => cannot load host verification key: %s", err)
	}
	// check loaded key is not private
	if pgp.HasPrivate() {
		return fmt.Errorf("verify => verification key should be public, private key found\n")
	}
	// verify digital signature
	return pgp.Verify(sum, sig)
}

// checksum create a checksum of the passed-in object
func checksum(obj interface{}) ([]byte, error) {
	source, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("checksum => cannot convert object to JSON to produce checksum: %s", err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	err = json.Indent(dest, source, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("checksum => cannot indent JSON to produce checksum: %s", err)
	}
	// create a new hash
	hash := sha256.New()
	// write object bytes into hash
	_, err = hash.Write(dest.Bytes())
	if err != nil {
		return nil, fmt.Errorf("checksum => cannot write JSON bytes to hash: %s", err)
	}
	// obtain checksum
	sum := hash.Sum(nil)
	return sum, nil
}

func executablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func homePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}
