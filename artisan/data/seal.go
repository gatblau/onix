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
    "path"
)

// Seal the digital Seal for a package
// the Seal contains information to determine if the package or its metadata has been compromised
// and therefore the Seal is broken
type Seal struct {
    // the package metadata
    Manifest *Manifest `json:"manifest"`
    // the combined checksum of the package and its metadata
    Digest string `json:"digest"`
    // the cryptographic seal for:
    // - author authentication: verify the author of the package
    // - data integrity: recognise if the package files differ from the original files at the time the package was built
    // - non-repudiation: since only the author (signer) can create the seal, any user of the package can present the package
    //      seal information to a third party as evidence if any dispute arises in the future
    Seal string `json:"seal,omitempty"`
}

// NoAuthority returns true if the seal does not have an authority
func (seal *Seal) NoAuthority() bool {
    return seal.Manifest == nil || (seal.Manifest != nil && len(seal.Manifest.Authority) == 0)
}

// DSha256 calculates the package SHA-256 digest by taking the combined checksum of the Seal information and the compressed file
func (seal *Seal) DSha256(path string) string {
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
    checksum := hash.Sum(nil)
    core.Debug("seal calculated base64 encoded checksum:\n>> start on next line\n%s\n>> ended on previous line", base64.StdEncoding.EncodeToString(checksum))
    return fmt.Sprintf("sha256:%s", base64.StdEncoding.EncodeToString(checksum))
}

func (seal *Seal) ZipFile(registryRoot string) string {
    return path.Join(core.RegistryPath(""), fmt.Sprintf("%s.zip", seal.Manifest.Ref))
}

func (seal *Seal) SealFile(registryRoot string) string {
    return path.Join(core.RegistryPath(""), fmt.Sprintf("%s.json", seal.Manifest.Ref))
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

// Valid checks that the digest stored in the seal is the same as the digest generated using the passed-in zip file path
// and the seal
// path: the path to the package zip file to validate
func (seal *Seal) Valid(path string) (valid bool, err error) {
    // calculates the digest using the zip file
    digest := seal.DSha256(path)
    // compare to the digest stored in the seal
    if seal.Digest == digest {
        return true, nil
    }
    return false, fmt.Errorf("downloaded package digest: %s does not match digest in manifest %s", digest, seal.Digest)
}
