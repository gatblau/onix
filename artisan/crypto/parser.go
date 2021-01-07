/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  This code has been based on https://github.com/AaronO/go-rsa-sign
*/
package crypto

// func ParseX509PrivateKey(data []byte, location string) (*rsa.PrivateKey, error) {
// 	pemData, err := pemParse(data, "RSA PRIVATE KEY", location)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return x509.ParsePKCS1PrivateKey(pemData)
// }

// func ParseX509PublicKey(data []byte, location string) (*rsa.PublicKey, error) {
// 	pemData, err := pemParse(data, "RSA PUBLIC KEY", location)
// 	if err != nil {
// 		return nil, err
// 	}
// 	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pubKey, ok := keyInterface.(*rsa.PublicKey)
// 	if !ok {
// 		return nil, fmt.Errorf("could not cast parsed key to *rsa.PublickKey")
// 	}
// 	return pubKey, nil
// }

// func pemParse(data []byte, pemType string, location string) ([]byte, error) {
// 	block, _ := pem.Decode(data)
// 	if block == nil {
// 		return nil, fmt.Errorf("cannot load RSA key from '%s', either not there or corrupted", location)
// 	}
// 	if pemType != "" && block.Type != pemType {
// 		return nil, fmt.Errorf("public key's type is '%s', expected '%s'", block.Type, pemType)
// 	}
// 	return block.Bytes, nil
// }
