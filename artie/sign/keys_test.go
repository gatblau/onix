package sign

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	size := 2048
	keyFilename, pubKey := KeyNames(".", "pgp", "pgp")
	key, _ := NewKeyPair(size)
	headers := make(map[string]string)
	SavePGPPrivateKey(keyFilename, key, headers)
	SavePGPPublicKey(pubKey, &key.PublicKey, headers)
	k, err := ReadPGPPrivateKey(keyFilename)
	if err != nil {
		fmt.Print(err)
		t.FailNow()
	}
	fmt.Print(k)
}
