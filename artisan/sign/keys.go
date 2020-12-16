/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package sign

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

// generate a new RSA key pair
func NewKeyPair(size int) (*rsa.PrivateKey, error) {
	reader := rand.Reader
	bitSize := size
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key pair: %s", err)
	}
	return key, nil
}

// save the specified RSA private key in OpenPGP armor format to a file
func SavePGPPrivateKey(filename string, key *rsa.PrivateKey, headers map[string]string) error {
	var buf = &bytes.Buffer{}
	w, err := armor.Encode(buf, openpgp.PrivateKeyType, headers)
	if err != nil {
		return fmt.Errorf("error creating OpenPGP Armor: %s", err)
	}
	pgpKey := packet.NewRSAPrivateKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		return fmt.Errorf("error serializing private key: %s", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("error serializing private key: %s", err)
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write key: %s", err)
	}
	return err
}

// save the specified RSA public key in OpenPGP armor format to a file
func SavePGPPublicKey(filename string, key *rsa.PublicKey, headers map[string]string) error {
	var buf = &bytes.Buffer{}
	w, err := armor.Encode(buf, openpgp.PublicKeyType, headers)
	if err != nil {
		fmt.Errorf("error creating OpenPGP Armor: %s", err)
	}
	pgpKey := packet.NewRSAPublicKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		fmt.Errorf("error serializing public key: %s", err)
	}
	err = w.Close()
	if err != nil {
		fmt.Errorf("error serializing public key: %s", err)
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write key: %s", err)
	}
	return err
}

// generates a private and public RSA keys for signing and verifying artefacts
func GeneratePGPKeys(path, name string, size int) {
	if size > 4500 {
		core.RaiseErr("maximum bit size 4500 exceeded")
	}
	if size == 0 {
		size = 2048
	}
	keyFilename, pubFilename := KeyNames(path, name, "pgp")
	key, err := NewKeyPair(size)
	core.CheckErr(err, "cannot create key pair")
	core.CheckErr(SavePGPPrivateKey(keyFilename, key, nil), "cannot save PGP private key")
	core.CheckErr(SavePGPPublicKey(pubFilename, &key.PublicKey, nil), "cannot save PGP public key")
}

func LoadPGPPrivateKey(group, name string) (*rsa.PrivateKey, error) {
	// first attempt to load the key from the registry/keys/group/name path
	private, _ := KeyNames(path.Join(core.RegistryPath(), "keys", group, name), fmt.Sprintf("%s_%s", group, name), "pgp")
	key, err := ReadPGPPrivateKey(private)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		private, _ = KeyNames(path.Join(core.RegistryPath(), "keys", group), group, "pgp")
		key, err = ReadPGPPrivateKey(private)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			private, _ = KeyNames(path.Join(core.RegistryPath(), "keys"), "root", "pgp")
			key, err = ReadPGPPrivateKey(private)
			if err != nil {
				return nil, fmt.Errorf("cannot read private pgp key: %s", err)
			}
		}
	}
	return key, nil
}

func LoadPGPPublicKey(group, name string) (*rsa.PublicKey, error) {
	// first attempt to load the key from the registry/keys/group/name path
	_, public := KeyNames(path.Join(core.RegistryPath(), "keys", group, name), fmt.Sprintf("%s_%s", group, name), "pgp")
	key, err := ReadPGPPublicKey(public)
	if err != nil {
		// if no luck, attempt to load the key from the registry/keys/group path
		_, public = KeyNames(path.Join(core.RegistryPath(), "keys", group), group, "pgp")
		key, err = ReadPGPPublicKey(public)
		if err != nil {
			// final attempt to load the key from the registry/keys/ path
			_, public = KeyNames(path.Join(core.RegistryPath(), "keys"), "root", "pgp")
			key, err = ReadPGPPublicKey(public)
			if err != nil {
				return nil, fmt.Errorf("cannot read public pgp key: %s", err)
			}
		}
	}
	return key, nil
}

func ReadPGPPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// open ascii armored private key
	in, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening private pgp key: %s", err)
	}
	defer in.Close()
	block, err := armor.Decode(in)
	if err != nil {
		return nil, fmt.Errorf("error decoding OpenPGP Armor: %s", err)
	}
	if block.Type != openpgp.PrivateKeyType {
		return nil, fmt.Errorf("invalid private pgp key file: error decoding private key")
	}
	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, fmt.Errorf("error reading private pgp key")
	}
	key, ok := pkt.(*packet.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private pgp key: error parsing private key")
	}
	if rsa, ok := key.PrivateKey.(*rsa.PrivateKey); ok {
		return rsa, nil
	}
	return nil, fmt.Errorf("invalid private pgp key, no RSA key found")
}

func ReadPGPPublicKey(filename string) (*rsa.PublicKey, error) {
	// open ascii armored public key
	in, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening public pgp key: %s", err)
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		return nil, fmt.Errorf("error decoding OpenPGP Armor: %s", err)
	}
	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("invalid public pgp key file: error decoding private key")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	kingpin.FatalIfError(err, "error reading private key")

	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public pgp key: error parsing public key")
	}
	if rsa, ok := key.PublicKey.(*rsa.PublicKey); ok {
		return rsa, nil
	}
	return nil, fmt.Errorf("invalid public pgp key, no RSA key found")
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
