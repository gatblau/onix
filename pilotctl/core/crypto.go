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

// Sign create a cryptographic signature for the passed-in object
func sign(obj interface{}) (string, error) {
	// only sign if we have an object
	if obj != nil {
		// load the signing key
		keyFile, err := signingKeyFile()
		if err != nil {
			return "", err
		}
		// retrieve the verification key from the specified location
		pgp, err := crypto.LoadPGP(keyFile, "")
		if err != nil {
			return "", fmt.Errorf("sign => cannot load signing key: %s", err)
		}
		// obtain the object checksum
		cs, err := checksum(obj)
		if err != nil {
			return "", fmt.Errorf("sign => cannot create checksum: %s", err)
		}
		signature, err := pgp.Sign(cs)
		if err != nil {
			return "", fmt.Errorf("sign => cannot create signature: %s", err)
		}
		// return a base64 encoded string with the digital signature
		return base64.StdEncoding.EncodeToString(signature), nil
	}
	return "", nil
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

func signingKeyFile() (string, error) {
	path := filepath.Join(executablePath(), ".pilot_sign.pgp")
	_, err := os.Stat(path)
	if err != nil {
		path = filepath.Join(homePath(), ".pilot_sign.pgp")
		_, err = os.Stat(path)
		if err != nil {
			path = "/keys/.pilot_sign.pgp"
			_, err = os.Stat(path)
			if err != nil {
				return "", fmt.Errorf("cannot find signing key")
			}
		}
		return path, nil
	}
	return path, nil
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
