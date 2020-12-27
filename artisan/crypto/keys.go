/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package crypto

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"os"
	"path"
	"path/filepath"
	"time"
)

// generates a private and public RSA keys for signing and verifying artefacts
func GeneratePGPKeys(path, prefix, name, comment, email string, size int) {
	if size > 4500 {
		core.RaiseErr("maximum bit size 4500 exceeded")
	}
	if size == 0 {
		size = 2048
	}
	hostname, err := os.Hostname()
	core.CheckErr(err, "cannot retrieve hostname")
	// if no name is provided, then use the hostname
	if len(name) == 0 {
		name = hostname
	}
	// if no comment is provided, then autogenerate a comment
	if len(comment) == 0 {
		comment = fmt.Sprintf("%d keys, generated on %s", size, time.Now())
	}
	// if no email is provided, then autogenerate a fake one
	if len(email) == 0 {
		email = fmt.Sprintf("artisan@%s.com", hostname)
	}
	// work out the file names for the keys
	keyFilename, pubFilename := KeyNames(path, prefix, "pgp")
	// create a new PGP entity
	pgp := NewPGP(name, comment, email, size)
	// save the public key part
	core.CheckErr(pgp.SavePublicKey(pubFilename), "cannot save public key")
	// save the private key part
	core.CheckErr(pgp.SavePrivateKey(keyFilename), "cannot save private key")
}

func LoadPGPPrivateKey(group, name string) (*PGP, error) {
	// first attempt to load the key from the registry/keys/group/name path
	private, _ := KeyNames(path.Join(core.RegistryPath(), "keys", group, name), fmt.Sprintf("%s_%s", group, name), "pgp")
	key, err := LoadPGP(private)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		private, _ = KeyNames(path.Join(core.RegistryPath(), "keys", group), group, "pgp")
		key, err = LoadPGP(private)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			private, _ = KeyNames(path.Join(core.RegistryPath(), "keys"), "root", "pgp")
			key, err = LoadPGP(private)
			if err != nil {
				return nil, fmt.Errorf("cannot read private pgp key: %s", err)
			}
		}
	}
	return key, nil
}

func LoadPGPPublicKey(group, name string) (*PGP, error) {
	// first attempt to load the key from the registry/keys/group/name path
	_, public := KeyNames(path.Join(core.RegistryPath(), "keys", group, name), fmt.Sprintf("%s_%s", group, name), "pgp")
	key, err := LoadPGP(public)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		_, public = KeyNames(path.Join(core.RegistryPath(), "keys", group), group, "pgp")
		key, err = LoadPGP(public)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			_, public = KeyNames(path.Join(core.RegistryPath(), "keys"), "root", "pgp")
			key, err = LoadPGP(public)
			if err != nil {
				return nil, fmt.Errorf("cannot read public pgp key: %s", err)
			}
		}
	}
	return key, nil
}

func PrivateKeyName(prefix string, extension string) string {
	if len(prefix) == 0 {
		prefix = "id"
	}
	return fmt.Sprintf("%s_rsa_key.%s", prefix, extension)
}

func PublicKeyName(prefix string, extension string) string {
	if len(prefix) == 0 {
		prefix = "id"
	}
	return fmt.Sprintf("%s_rsa_pub.%s", prefix, extension)
}

// works out the fully qualified names of the private and public RSA keys
func KeyNames(path, prefix string, extension string) (key string, pub string) {
	if len(path) == 0 {
		path = "."
	}
	// if the path is relative then make it absolute
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		core.CheckErr(err, "cannot return an absolute representation of path: '%s'", path)
		path = p
	}
	// works out the private key name
	keyName := filepath.Join(path, PrivateKeyName(prefix, extension))
	// works out the public key name
	pubName := filepath.Join(path, PublicKeyName(prefix, extension))
	return keyName, pubName
}
