package tkn

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/flow"
	"testing"
)

func TestBuild(t *testing.T) {
	f, err := flow.LoadFlow("/Users/andresalos/go/src/github.com/gatblau/onix/artisan/flow/ci-flow-merged.yaml")
	core.CheckErr(err, "cannot load flow")
	pk, err := crypto.LoadPGP("/Users/andresalos/go/src/github.com/gatblau/onix/artisan/flow/id_rsa_key.pgp")
	core.CheckErr(err, "cannot load decryption key")
	builder := NewBuilder(f, pk)
	buf := builder.Create()
	fmt.Println(buf.String())
}
