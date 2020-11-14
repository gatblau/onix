/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

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

type Manifest struct {
	// the package type
	Type string `json:"type,omitempty"`
	// the license associated to the package
	License string `json:"license"`
	// the name of the package file
	Ref string `json:"ref"`
	// the build profile used
	Profile string `json:"profile"`
	// the labels assigned to the package in the Pakfile
	Labels map[string]string `json:"labels,omitempty"`
	// the URI of the package source
	Source string `json:"source"`
	// the path within the source where the project is (for uber repos)
	SourcePath string `json:"source_path,omitempty"`
	// the commit hash
	Commit string `json:"commit"`
	// repo branch
	Branch string `json:"branch,omitempty"`
	// repo tag
	Tag string `json:"tag,omitempty"`
	// the name of the file or folder that has been packaged
	Target string `json:"target,omitempty"`
	// the timestamp
	Time string `json:"time"`
	// the size of the artefact
	Size string `json:"size"`
}
