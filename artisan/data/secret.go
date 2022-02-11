package data

import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"strings"
)

// Secret describes the secrets required by functions
type Secret struct {
	// the unique reference for the secret
	Name string `yaml:"name" json:"name"`
	// a description of the intended use or meaning of this secret
	Description string `yaml:"description" json:"description"`
	// the value of the secret
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	// the value is required
	Required bool `yaml:"required,omitempty" json:"required,omitempty"`
}

func (s *Secret) Encrypt(pubKey *crypto.PGP) error {
	encValue, err := pubKey.Encrypt([]byte(s.Value))
	if err != nil {
		return fmt.Errorf("cannot encrypt secret %s: %s", s.Name, err)
	}
	s.Value = base64.StdEncoding.EncodeToString(encValue)
	return nil
}

func (s *Secret) Decrypt(pk *crypto.PGP) {
	if !pk.HasPrivate() {
		core.RaiseErr("provided key is not private")
	}
	// decode encrypted value
	decoded, err := base64.StdEncoding.DecodeString(s.Value)
	core.CheckErr(err, "cannot decode encrypted value using base64")
	decValueBytes, err := pk.Decrypt([]byte(decoded))
	core.CheckErr(err, "cannot decrypt secret")
	s.Value = string(decValueBytes)
}

type Secrets []*Secret

func (list Secrets) Len() int { return len(list) }

func (list Secrets) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

func (list Secrets) Less(i, j int) bool {
	var si string = list[i].Name
	var sj string = list[j].Name
	var si_lower = strings.ToLower(si)
	var sj_lower = strings.ToLower(sj)
	if si_lower == sj_lower {
		return si < sj
	}
	return si_lower < sj_lower
}
