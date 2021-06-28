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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/crypto"
	"os"
	"path/filepath"
)

// sign create a cryptographic signature for the passed-in object
func verify(obj interface{}, signature string) error {
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
	// retrieve the verification key from the specified location
	keyFile, err := verifyKeyFile()
	if err != nil {
		return fmt.Errorf("verify => cannot find host verification key: %s", err)
	}
	pgp, err := crypto.LoadPGP(keyFile)
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

func verifyKeyFile() (string, error) {
	verifyKeyFilename := ".pilot_verify.pgp"
	file := filepath.Join(executablePath(), verifyKeyFilename)
	_, err := os.Stat(filepath.Join(executablePath(), verifyKeyFilename))
	if err != nil {
		file = filepath.Join(homePath(), verifyKeyFilename)
		_, err = os.Stat(filepath.Join(homePath(), verifyKeyFilename))
		if err != nil {
			return "", fmt.Errorf("cannot find signing key")
		}
		return file, nil
	}
	return file, nil
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
