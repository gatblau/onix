/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// convert the passed in parameter to a Json Byte Array
func ToJsonBytes(s interface{}) []byte {
	// serialise the seal to json
	source, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	err = json.Indent(dest, source, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return dest.Bytes()
}

func ToJsonFile(obj interface{}) (*os.File, error) {
	// create an UUId
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	// generate an internal random and transient name based on the UUId
	name := path.Join(TmpPath(), fmt.Sprintf("%s.json", uuid))
	// marshals the object into Json bytes
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// create a transient temp file
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	// write the bytes into it
	_, err = file.Write(b)
	if err != nil {
		return nil, err
	}
	// closes the file
	file.Close()
	// open the created file
	file = openFile(name)
	if err != nil {
		return nil, err
	}
	// remove the file from the file system
	err = os.Remove(name)
	// return the File object
	return file, err
}

func openFile(path string) *os.File {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return r
}

// remove an element in a slice
func RemoveElement(a []string, value string) []string {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == value {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = ""   // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

// the artefact id calculated as the hex encoded SHA-256 digest of the artefact Seal
func ArtefactId(seal *Seal) string {
	// serialise the seal info to json
	info := ToJsonBytes(seal)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Printf("error: %s - %s\n", fmt.Sprintf(msg, a...), err)
		os.Exit(1)
	}
}

func RaiseErr(msg string, a ...interface{}) {
	fmt.Printf("error: %s\n", fmt.Sprintf(msg, a...))
	os.Exit(1)
}

func Msg(msg string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Printf("info: %s\n", fmt.Sprintf(msg, a...))
	} else {
		fmt.Printf("info: %s\n", msg)
	}
}

// takes the combined checksum of the Seal information and the compressed file
func SealChecksum(path string, sealData *Manifest) []byte {
	// read the compressed file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// serialise the seal info to json
	info := ToJsonBytes(sealData)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	// copy the compressed file into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return hash.Sum(nil)
}

// gets the checksum of the passed string
func StringCheckSum(value string) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader([]byte(value))); err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func UserPwd(creds string) (user, pwd string) {
	// if credentials not provided then no user / pwd
	if len(creds) == 0 {
		return "", ""
	}
	parts := strings.Split(creds, ":")
	if len(parts) != 2 {
		log.Fatal(errors.New("credentials in incorrect format, they should be USER:PWD"))
	}
	return parts[0], parts[1]
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
