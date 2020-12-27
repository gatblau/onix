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
	"testing"
)

func TestGenerateKeys(t *testing.T) {
	// creates a new PGP object
	p := NewPGP("gatblau/boot", "an artisan pgp key for digital signatures", "onix@gatblau.org", 2048)

	// saves private and public keys
	pkName, pubName := KeyNames(".", "root", "pgp")
	err := p.SaveKeyPair(pubName, pkName)
	if err != nil {
		t.FailNow()
	}
}

func TestLoadKeySignAndVerify(t *testing.T) {
	// load the private key for signing
	priv, err := LoadPGP("priv.key")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// sign the message
	signature, err := priv.Sign([]byte("Hello World"))
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// load the public key for verification of the signature
	pub, err := LoadPGP("pub.key")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// verify the signature
	err = pub.Verify([]byte("Hello World"), signature)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	// anyone can encrypt with the public key
	// load the public key for encryption of the message
	pub, err := LoadPGP("pub.key")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	encrypted, err := pub.Encrypt([]byte("Hello World"))
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// only the holder of the private key can decrypt
	// load the private key for decryption of the message
	priv, err := LoadPGP("priv.key")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	msg, err := priv.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Print(msg)
}
