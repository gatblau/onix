/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  originally taken from gist https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09
*/

package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// TlsSignatureAlgorithm the algorithm to use to generate the TLS certificate
type TlsSignatureAlgorithm int

const (
	// RSA Rivest–Shamir–Adleman Algorithm
	RSA TlsSignatureAlgorithm = iota
	// ECDSA Elliptic Curve Digital Signature Algorithm
	ECDSA
)

// TLSCertificate generates a self-signed TLS x509 certificate and a private key using Elliptic Curve Digital
// Signature Algorithm (ECDSA)
func TLSCertificate(algor TlsSignatureAlgorithm, organisation string, hosts []string, isCA bool) (cert []byte, key []byte, err error) {
	var priv interface{}
	switch algor {
	case RSA:
		priv, err = rsa.GenerateKey(rand.Reader, 2048)
	case ECDSA:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if err != nil {
		return
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{
				organisation,
			},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}
	// add IP addresses and or DNS names to the certificate
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}
	// if the certificates are for a certificate authority
	if isCA {
		// set it as CA cert
		template.IsCA = true
		// add use for signing certs
		template.KeyUsage |= x509.KeyUsageCertSign
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pubKey(priv), priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %s", err)
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	cert = out.Bytes()
	out.Reset()
	pem.Encode(out, pemBlock(priv))
	key = out.Bytes()
	return
}

// pubKey returns the public key for the specified private key
func pubKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

// pemBlock returns a PEM block for the passed-in RSA or ECDSA private key
func pemBlock(privateKey interface{}) *pem.Block {
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
