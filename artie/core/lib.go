package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash.Sum(nil)))
}

func CheckErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Printf("error: %s - %s\n", fmt.Sprintf(msg, a), err)
		os.Exit(1)
	}
}

func RaiseErr(msg string, a ...interface{}) {
	fmt.Printf("error: %s\n", fmt.Sprintf(msg, a))
	os.Exit(1)
}

func Msg(msg string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Printf("artie: %s\n", fmt.Sprintf(msg, a...))
	} else {
		fmt.Printf("artie: %s\n", msg)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 80))
}
