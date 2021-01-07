/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gatblau/onix/artisan/core"
	"io"
	"log"
	"os"
)

// the digital Seal for a package
// the Seal contains information to determine if the package or its metadata has been compromised
// and therefore the Seal is broken
type Seal struct {
	// the package metadata
	Manifest *Manifest `json:"manifest"`
	// the combined checksum of the package and its metadata
	Digest string `json:"digest"`
	// the cryptographic signature for:
	// - package & seal authentication:
	//     When the verifier validates the digital signature using public key of the package author, he is assured that signature has been created only by author who possess the corresponding secret private key and no one else.
	// - package & seal data integrity:
	//     In case an attacker has access to the package & seal and modifies them, the digital signature verification at receiver end fails.
	//     The hash of modified package and seal and the output provided by the verification algorithm will not match. Hence, receiver can safely deny the package & seal content assuming that data integrity has been breached.
	//     The seal is broken.
	// - non-repudiation:
	//     Since it is assumed that only the author (signer) has the knowledge of the signature key, they can only create unique signature on a given package.
	//     Thus the receiver can present the package and the digital signature to a third party as evidence if any dispute arises in the future.
	Signature string `json:"signature,omitempty"`
}

// takes the combined checksum of the Seal information and the compressed file
func (seal *Seal) Checksum(path string) []byte {
	// precondition: the manifest is required
	if seal.Manifest == nil {
		core.RaiseErr("seal has no manifest, cannot create checksum")
	}
	// read the compressed file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// serialise the seal info to json
	info := core.ToJsonBytes(seal.Manifest)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	// copy the compressed file into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return hash.Sum(nil)
}

// the artefact id calculated as the hex encoded SHA-256 digest of the artefact Seal
func (seal *Seal) ArtefactId() string {
	// serialise the seal info to json
	info := core.ToJsonBytes(seal)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
