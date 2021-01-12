/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package crypto

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/docs"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// PGP entity for signing, verification, encryption and decryption
type PGP struct {
	entity  *openpgp.Entity
	conf    *packet.Config
	name    string
	comment string
	email   string
}

const (
	defaultDigest = crypto.SHA256
	defaultCipher = packet.CipherAES128
)

// creates a new PGP entity
func NewPGP(name, comment, email string, bits int) *PGP {
	var p = &PGP{
		name:    name,
		comment: comment,
		email:   email,
		conf: &packet.Config{
			DefaultCipher: defaultCipher,
			DefaultHash:   defaultDigest,
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
	entity := entityList[0]

	// NOTE: if this is a public key, adds the default cipher to the id self-signature so that the public key can be used to
	// encrypt messages without failing with message:
	// "cannot encrypt because no candidate hash functions are compiled in. (Wanted RIPEMD160 in this case.)
	// PGP Encrypt defaults to PreferredSymmetric=Cast5 & PreferredHash=Ripemd160
	// To avoid the error above, it has to change the required values
	// It needs to be in the list of preferred algorithms specified in the self-signature of the primary identity
	// there should only be one, but cycle over all identities for completeness
	for _, id := range entity.Identities {
		preferredHashId, _ := s2k.HashToHashId(defaultDigest)
		id.SelfSignature.PreferredHash = []uint8{preferredHashId}
		id.SelfSignature.PreferredSymmetric = []uint8{uint8(defaultCipher)}
	}

	return &PGP{
		entity: entity,
	}, nil
}

// check if the PGP entity has a private key, if not an error is returned
func (p *PGP) HasPrivate() bool {
	if p.entity == nil {
		core.RaiseErr("PGP object does not contain entity")
	}
	return p.entity.PrivateKey != nil
}

// signs the specified message (requires loading a private key)
func (p *PGP) Sign(message []byte) ([]byte, error) {
	writer := new(bytes.Buffer)
	reader := bytes.NewReader(message)
	err := openpgp.ArmoredDetachSign(writer, p.entity, reader, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot sign message: %s", err)
	}
	return writer.Bytes(), nil
}

// verifies the message using a specified signature (requires loading a public key)
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

// encrypts the specified message
func (p *PGP) Encrypt(message []byte) ([]byte, error) {
	// create buffer to write output to
	buf := new(bytes.Buffer)
	// create armor format encoder
	encoderWriter, err := armor.Encode(buf, "Message", make(map[string]string))
	if err != nil {
		return []byte{}, fmt.Errorf("cannot create PGP armor: %v", err)
	}
	// create the encryptor with the encoder
	encryptorWriter, err := openpgp.Encrypt(encoderWriter, []*openpgp.Entity{p.entity}, nil, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot create encryptor: %v", err)
	}
	// create the compressor with the encryptor
	compressorWriter, err := gzip.NewWriterLevel(encryptorWriter, gzip.BestCompression)
	if err != nil {
		return []byte{}, fmt.Errorf("invalid compression level: %v", err)
	}
	// write the message to the compressor
	messageReader := bytes.NewReader(message)
	_, err = io.Copy(compressorWriter, messageReader)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot write data to the compressor: %v", err)
	}
	compressorWriter.Close()
	encryptorWriter.Close()
	encoderWriter.Close()
	// returns an encoded, encrypted, and compressed message
	return buf.Bytes(), nil
}

// decrypts the specified message
func (p *PGP) Decrypt(encrypted []byte) ([]byte, error) {
	// Decode message
	block, err := armor.Decode(bytes.NewReader(encrypted))
	if err != nil {
		return []byte{}, fmt.Errorf("cannot decode the PGP armor encrypted string: %v", err)
	}
	if block.Type != "Message" {
		return []byte{}, errors.New("invalid message type")
	}
	// decrypt the message
	entityList := openpgp.EntityList{
		p.entity,
	}
	messageReader, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read message: %v", err)
	}
	read, err := ioutil.ReadAll(messageReader.UnverifiedBody)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read unverified body: %v", err)
	}
	// unzip the message
	reader := bytes.NewReader(read)
	uncompressed, err := gzip.NewReader(reader)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot initialise gzip reader: %v", err)
	}
	defer uncompressed.Close()
	out, err := ioutil.ReadAll(uncompressed)
	if err != nil {
		return []byte{}, err
	}
	// return the unencoded, unencrypted, and uncompressed message
	return out, nil
}

func (p *PGP) SavePublicKey(keyFilename string) error {
	keyBytes, err := p.toPublicKey()
	if err != nil {
		return fmt.Errorf("cannot save public key: %s", err)
	}
	// write the public key to a file
	err = ioutil.WriteFile(keyFilename, keyBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save public key: %s", err)
	}
	return nil
}

func (p *PGP) SavePrivateKey(keyFilename string) error {
	keyBytes, err := p.toPrivateKey()
	if err != nil {
		return fmt.Errorf("cannot save private key: %s", err)
	}
	// write the private key to a file
	err = ioutil.WriteFile(keyFilename, keyBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save private key: %s", err)
	}
	return nil
}

func (p *PGP) toPrivateKey() (privateKey []byte, err error) {
	// the buffer to contain the serialised key pair
	privateBuf := new(bytes.Buffer)
	// serialises the private key into the buffer
	p.entity.SerializePrivate(privateBuf, p.conf)
	// create PEM headers for the PGP key
	headers := pemHeaders(cipherToString(p.conf.DefaultCipher), p.conf.DefaultHash.String(), p.conf.RSABits, p.conf.Time())
	// encode the buffer containing the serialised private key into Armor ASCII format
	privateKey, err = armorEncode(privateBuf, openpgp.PrivateKeyType, headers)
	if err != nil {
		return nil, fmt.Errorf("cannot armor encode private key: %s", err)
	}
	return privateKey, nil
}

func (p *PGP) toPublicKey() (publicKey []byte, err error) {
	// the buffer to contain the serialised public key
	publicBuf := new(bytes.Buffer)
	// serialises the public key into the buffer
	p.entity.Serialize(publicBuf)
	// create PEM headers for the PGP key
	headers := pemHeaders(cipherToString(p.conf.DefaultCipher), p.conf.DefaultHash.String(), p.conf.RSABits, p.conf.Time())
	// encode the buffer containing the serialised public key into Armor ASCII format
	publicKey, err = armorEncode(publicBuf, openpgp.PublicKeyType, headers)
	if err != nil {
		return nil, fmt.Errorf("cannot armor encode public key: %s", err)
	}
	return publicKey, nil
}

// create PEM headers for the PGP key
func pemHeaders(cipher, hash string, rsaBits int, time time.Time) map[string]string {
	headers := map[string]string{
		"Version": fmt.Sprintf("golang.org/x/crypto/openpgp - artisan-%s", docs.Version),
		"Comment": fmt.Sprintf("Cipher: %s, Hash: %s, RSA Bits: %s, Created: %s", cipher, hash, strconv.Itoa(rsaBits), time.String()),
		"Hash":    fmt.Sprintf("%s/%s", cipher, hash),
	}
	return headers
}

// armor ascii encode the passed in buffer
func armorEncode(key *bytes.Buffer, keyType string, headers map[string]string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
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
		return "AES256"
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
