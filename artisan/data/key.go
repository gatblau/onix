package data

import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"strings"
)

// describes PGP keys required by functions
type Key struct {
	// the unique reference for the PGP key
	Name string `yaml:"name"`
	// a description of the intended use of this key
	Description string `yaml:"description"`
	// indicates if the referred key is private or public
	Private bool `yaml:"private"`
	// the artisan package group used to select the key
	PackageGroup string `yaml:"package_group,omitempty" json:"package_group,omitempty"`
	// the artisan package name used to select the key
	PackageName string `yaml:"package_name,omitempty" json:"package_name,omitempty"`
	// indicates if this key should be aggregated with other keys
	Aggregate bool `yaml:"aggregate,omitempty" json:"aggregate,omitempty"`
	// the key content
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	// the path to the key in the Artisan registry
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
}

func (k *Key) Encrypt(pubKey *crypto.PGP) error {
	encValue, err := pubKey.Encrypt([]byte(k.Value))
	if err != nil {
		return fmt.Errorf("cannot encrypt PGP key %s: %s", k.Name, err)
	}
	k.Value = base64.StdEncoding.EncodeToString(encValue)
	return nil
}

func (k *Key) Decrypt(privateKey *crypto.PGP) error {
	decoded, err := base64.StdEncoding.DecodeString(k.Value)
	core.CheckErr(err, "cannot base64 decode key '%s'", k.Name)
	decValue, err := privateKey.Decrypt(decoded)
	if err != nil {
		return fmt.Errorf("cannot decrypt PGP key %s: %s", k.Name, err)
	}
	k.Value = string(decValue)
	return nil
}

type Keys []*Key

func (list Keys) Len() int { return len(list) }

func (list Keys) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

func (list Keys) Less(i, j int) bool {
	var si string = list[i].Name
	var sj string = list[j].Name
	var si_lower = strings.ToLower(si)
	var sj_lower = strings.ToLower(sj)
	if si_lower == sj_lower {
		return si < sj
	}
	return si_lower < sj_lower
}
