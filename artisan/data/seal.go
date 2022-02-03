/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"io"
	"io/ioutil"
)

// Seal the digital Seal for a package
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
	//     Thus, the receiver can present the package and the digital signature to a third party as evidence if any dispute arises in the future.
	Signature string `json:"signature,omitempty"`
}

// Checksum takes the combined checksum of the Seal information and the compressed file
func (seal *Seal) Checksum(path string) []byte {
	// precondition: the manifest is required
	if seal.Manifest == nil {
		core.RaiseErr("seal has no manifest, cannot create checksum")
	}
	// read the compressed file
	file, err := ioutil.ReadFile(path)
	core.CheckErr(err, "cannot open seal file")
	// serialise the seal info to json
	info := core.ToJsonBytes(seal.Manifest)
	core.Debug("manifest before checksum:\n>> start on next line\n%s\n>> ended on previous line", string(info))
	hash := sha256.New()
	written, err := hash.Write(file)
	core.CheckErr(err, "cannot write package file to hash")
	core.Debug("%d bytes from package written to hash", written)
	written, err = hash.Write(info)
	core.CheckErr(err, "cannot write manifest to hash")
	core.Debug("%d bytes from manifest written to hash", written)
	sum := hash.Sum(nil)
	core.Debug("seal calculated base64 encoded checksum:\n>> start on next line\n%s\n>> ended on previous line", base64.StdEncoding.EncodeToString(sum))
	return sum
}

// PackageId the package id calculated as the hex encoded SHA-256 digest of the artefact Seal
func (seal *Seal) PackageId() (string, error) {
	// serialise the seal info to json
	info := core.ToJsonBytes(seal)
	hash := sha256.New()
	// copy the seal content into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		return "", fmt.Errorf("cannot create hash from package seal: %s", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
