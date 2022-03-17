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

// LoadKeys load primary and backup keys to use for a given package name
func LoadKeys(name core.PackageName, isPrivate bool) (primaryKey, backupKey *PGP, err error) {
	primaryKeyPath, backupKeyPath := resolveKeyPath(name, isPrivate)
	// if the primary key does not exist
	if !pathExists(primaryKeyPath) {
		// cannot continue
		return nil, nil, fmt.Errorf("primary key not found in path %s", primaryKeyPath)
	}
	// tries and load the primary key
	primaryKey, err = LoadPGP(primaryKeyPath, "")
	// if it failed
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load primary key from '%s': %s", primaryKeyPath, err)
	}
	// only loads the backup key if it exists
	if pathExists(backupKeyPath) {
		// tries and loads it
		backupKey, err = LoadPGP(backupKeyPath, "")
		// if there was an error trying to load the backup key
		if err != nil {
			// return without backup key, but with an error
			return primaryKey, nil, fmt.Errorf("cannot load backup key from '%s': %s", backupKeyPath, err)
		}
	}
	// return both keys
	return primaryKey, backupKey, nil
}

// resolveKeyPath returns the primary and backup key paths to use for a given package name and keys deployed in the local artisan registry
func resolveKeyPath(name core.PackageName, isPrivate bool) (primary, backup string) {
	keyPaths := getAllKeyPaths()
	var path string
	for _, s := range keyPaths {
		if strings.HasPrefix(name.Repository(), s) {
			path = s
			break
		}
	}
	if len(path) == 0 {
		primary = filepath.Join(core.KeysPath(), fmt.Sprintf("root%s", keySuffix(isPrivate)))
		backup = filepath.Join(core.KeysPath(), fmt.Sprintf("root_backup%s", keySuffix(isPrivate)))
	} else {
		primary = filepath.Join(core.KeysPath(), path, fmt.Sprintf("%s%s", strings.ReplaceAll(path, "/", "_"), keySuffix(isPrivate)))
		backup = filepath.Join(core.KeysPath(), path, fmt.Sprintf("%s_backup%s", strings.ReplaceAll(path, "/", "_"), keySuffix(isPrivate)))
	}
	return primary, backup
}

// getAllKeyPaths returns a list of all possible key paths in the priority they should be used
// it is driven from the existence of keys in the keys folder hierarchy
func getAllKeyPaths() (paths []string) {
	// defines a function to check if a slice slice contains a given element
	var contains = func(elems []string, v string) bool {
		for _, s := range elems {
			if v == s {
				return true
			}
		}
		return false
	}
	// create a file system instance with root in the folder where keys are stored in the artisan local registry
	fSys := os.DirFS(core.KeysPath())
	// variable to track walked path without its last element (the folder without filename)
	var folder string
	// walks the keys' root folder tree to collect a list of sub folders that contain keys
	fs.WalkDir(fSys, ".", func(p string, d fs.DirEntry, err error) error {
		// stores the path without filename for the path being walked
		folder = filepath.Dir(p)
		// if a key file was found
		if !d.IsDir() && // the path is not a directory
			filepath.Ext(d.Name()) == ".pgp" && // the file in the path is a pgp key (*.pgp extension)
			!contains(paths, folder) { // the path has not been previously recorded
			// tracks the path where the key is
			paths = append(paths, folder)
		}
		return nil
	})
	// because the paths slice is in top to bottom ordered, reverse it so that is ordered from the bottom to the top
	// this is, the bottom most path overrides the top most path (i.e. specific keys override more generic ones if defined)
	sort.Slice(paths, func(i, j int) bool {
		// use the path length to compare for ordering (i.e. the longest path has priority over the shorter; because
		// longest paths will always be folders at the bottom of the folder tree)
		return len(paths[i]) > len(paths[j])
	})
	return paths
}

// keySuffix returns the suffix to use for a key depending on whether the key is public or private
func keySuffix(private bool) string {
	if private {
		return "_rsa_key.pgp"
	}
	return "_rsa_pub.pgp"
}

// pathExists returns true if the specified path exists; otherwise false
func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
