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
	"crypto"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/docs"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PGP struct {
	entity  *openpgp.Entity
	conf    *packet.Config
	name    string
	comment string
	email   string
}

func Load(keypath string) *PGP {
	return nil
}

// creates a new PGP entity
func NewPGP(name, comment, email string, bits int) *PGP {
	var p = &PGP{
		name:    name,
		comment: comment,
		email:   email,
		conf: &packet.Config{
			DefaultCipher: packet.CipherAES128,
			DefaultHash:   crypto.SHA256,
			RSABits:       bits,
			Time: func() time.Time {
				return time.Now()
			},
		},
	}
	// creates a pgp entity that contains a fresh RSA/RSA keypair with a
	// single identity composed of the given full name, comment and email
	entity, err := openpgp.NewEntity(p.name, p.comment, p.email, p.conf)
	if err != nil {
		panic(err)
	}
	p.entity = entity
	return p
}

// load a PGP entity from file
func LoadPGP(filename string) (*PGP, error) {
	if !filepath.IsAbs(filename) {
		abs, err := filepath.Abs(filename)
		if err != nil {
			return nil, fmt.Errorf("cannot convert path %s to absolute path: %s", filename, err)
		}
		filename = abs
	}
	// read the key file
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open key file %s: %s", filename, err)
	}
	entityList, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read PGP entity: %s", err)
	}
	if len(entityList) == 0 {
		return nil, fmt.Errorf("no PGP entities found in %s", filename)
	}
	return &PGP{
		entity: entityList[0],
	}, nil
}

func (p *PGP) Sign(message []byte) ([]byte, error) {
	writer := new(bytes.Buffer)
	reader := bytes.NewReader(message)
	err := openpgp.ArmoredDetachSign(writer, p.entity, reader, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot sign message: %s", err)
	}
	return writer.Bytes(), nil
}

func (p *PGP) Verify(message []byte, signature []byte) error {
	sig, err := parseSignature(signature)
	if err != nil {
		return err
	}
	hash := sig.Hash.New()
	messageReader := bytes.NewReader(message)
	io.Copy(hash, messageReader)
	err = p.entity.PrimaryKey.VerifySignature(hash, sig)
	if err != nil {
		return err
	}
	return nil
}

// save the PGP private and public keys to a file
func (p *PGP) SaveKeyPair(path, publicKeyFilename, privateKeyFilename string) error {
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("cannot convert to absolute path: %s", err)
		}
		path = absPath
	}
	priv, pub, err := p.toKeyPair()
	if err != nil {
		return fmt.Errorf("cannot save key pair: %s", err)
	}
	// write the private key to a file
	err = ioutil.WriteFile(filepath.Join(path, privateKeyFilename), priv, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save private key: %s", err)
	}
	// write the public key to a file
	err = ioutil.WriteFile(filepath.Join(path, publicKeyFilename), pub, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save public key: %s", err)
	}
	return nil
}

// return the armor ascii encoded private and public keys
func (p *PGP) toKeyPair() (privateKey, publicKey []byte, err error) {
	// the buffer to contain the serialised key pair
	privateBuf := new(bytes.Buffer)
	// serialises the private key into the buffer
	p.entity.SerializePrivate(privateBuf, p.conf)
	// encode the buffer containing the serialised private key into Armor ASCII format
	privateKey, err = armorEncode(privateBuf, openpgp.PrivateKeyType, p.conf)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot armor encode private key: %s", err)
	}
	// the buffer to contain the serialised public key
	publicBuf := new(bytes.Buffer)
	// serialises the public key into the buffer
	p.entity.Serialize(publicBuf)
	// encode the buffer containing the serialised public key into Armor ASCII format
	publicKey, err = armorEncode(publicBuf, openpgp.PublicKeyType, p.conf)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot armor encode public key: %s", err)
	}
	return privateKey, publicKey, nil
}

// armor ascii encode the passed in buffer
func armorEncode(key *bytes.Buffer, keyType string, conf *packet.Config) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	// adds PEM block headers with some information about the keys
	headers := map[string]string{
		"Version": fmt.Sprintf("golang.org/x/crypto/openpgp - artisan-%s", docs.Version),
		"Cipher":  cipherToString(conf.DefaultCipher),
		"Hash":    conf.DefaultHash.String(),
		"RSABits": strconv.Itoa(conf.RSABits),
		"Time":    conf.Time().String(),
	}
	w, err := armor.Encode(buf, keyType, headers)
	if err != nil {
		return nil, fmt.Errorf("cannot encode keys in armor format: %s", err)
	}
	_, err = w.Write(key.Bytes())
	if err != nil {
		return nil, fmt.Errorf("\"error armoring serializedEntity: %s", err)
	}
	w.Close()
	return buf.Bytes(), nil
}

// returns the string representation of the passed i cipher function
func cipherToString(cipher packet.CipherFunction) string {
	switch cipher {
	case 2:
		return "3DES"
	case 3:
		return "CAST5"
	case 7:
		return "AES128"
	case 8:
		return "AES192"
	case 9:
		return "AES192"
	default:
		return "NotKnown"
	}
}

// parses a string of bytes containing a PGP signature
func parseSignature(signature []byte) (*packet.Signature, error) {
	signatureReader := bytes.NewReader(signature)
	block, err := armor.Decode(signatureReader)
	if err != nil {
		return nil, fmt.Errorf("cannot decode OpenPGP Armor: %s", err)
	}
	if block.Type != openpgp.SignatureType {
		return nil, errors.New("invalid signature file")
	}
	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}
	sig, ok := pkt.(*packet.Signature)
	if !ok {
		return nil, errors.New("cannot parse PGP signature")
	}
	return sig, nil
}
