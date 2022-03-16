/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package crypto

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type KeyLoader struct {
	root string
	m    []string
}

func newKeyLoader() *KeyLoader {
	loader := &KeyLoader{
		root: core.KeysPath(),
		m:    []string{},
	}
	fSys := os.DirFS(loader.root)
	var path string
	fs.WalkDir(fSys, ".", func(p string, d fs.DirEntry, err error) error {
		// if a key file was found
		path = filepath.Dir(p)
		if !d.IsDir() && filepath.Ext(d.Name()) == ".pgp" && !contains(loader.m, path) {
			loader.m = append(loader.m, path)
		}
		return nil
	})
	sort.Slice(loader.m, func(i, j int) bool {
		return len(loader.m[i]) > len(loader.m[j])
	})
	return loader
}

func (k *KeyLoader) Key(name *core.PackageName, private bool) (*PGP, error) {
	keyPath := k.resolve(name, private)
	return LoadPGP(keyPath, "")
}

// resolve the path of the key to use
func (k *KeyLoader) resolve(name *core.PackageName, private bool) string {
	var path, keyPath, keyName string
	if private {
		keyName = "_key"
	} else {
		keyName = "_pub"
	}
	for _, s := range k.m {
		if strings.HasPrefix(name.Repository(), s) {
			path = s
			break
		}
	}
	if len(path) == 0 {
		keyPath = filepath.Join(k.root, fmt.Sprintf("root_rsa%s.pgp", keyName))
	} else {
		keyPath = filepath.Join(k.root, path, fmt.Sprintf("%s_rsa%s.pgp", strings.ReplaceAll(path, "/", "_"), keyName))
	}
	return keyPath
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
