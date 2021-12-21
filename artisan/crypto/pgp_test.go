/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
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
	err := p.SavePublicKey(pubName, "xxx", "")
	if err != nil {
		t.FailNow()
	}
	err = p.SavePrivateKey(pkName, "xxx", "")
	if err != nil {
		t.FailNow()
	}
}

func TestLoadKeySignAndVerify(t *testing.T) {
	// load the private key for signing
	priv, err := LoadPGP("id_rsa_key.pgp", "")
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
	pub, err := LoadPGP("id_rsa_pub.pgp", "")
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
	pub, err := LoadPGP("pub.key", "")
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
	priv, err := LoadPGP("priv.key", "")
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

func TestGenerateEncryptedPrivateKey(t *testing.T) {
	// creates a new PGP object
	p := NewPGP("gatblau/boot", "an artisan pgp key for digital signatures", "onix@gatblau.org", 2048)

	// defines a 16-digit passphrase (AES-128)
	// for AES-192 use a 24-digit passphrase
	// for AES-256 use a 32-digit passphrase
	aes_128_passphrase := "0123456789012345"

	// save encrypted key as *.asc file using passphrase
	err := p.SavePrivateKey("myprivatekey.asc", "0.0.1", aes_128_passphrase)
	if err != nil {
		t.Fatal(err)
	}

	// now try and load file and decrypt it using passphrase
	p2, err := LoadPGP("myprivatekey.asc", aes_128_passphrase)
	if err != nil {
		t.Fatal(err)
	}

	// now save pgp key unencrypted, hence no passphrase
	err = p2.SavePrivateKey("myprivatekey.pgp", "0.0.1", "")
	if err != nil {
		t.Fatal(err)
	}
}
