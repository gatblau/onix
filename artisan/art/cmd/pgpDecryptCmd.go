/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// list local packages
type PGPDecryptCmd struct {
	cmd     *cobra.Command
	keyPath string
	group   string
	name    string
}

func NewPGPDecryptCmd() *PGPDecryptCmd {
	c := &PGPDecryptCmd{
		cmd: &cobra.Command{
			Use:   "decrypt [flags] filename",
			Short: "decrypts a file using a designated PRIVATE PGP key",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.keyPath, "private-key", "k", "", "the path to the PRIVATE PGP key to use")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "the package group of the PRIVATE PGP key to use")
	c.cmd.Flags().StringVarP(&c.keyPath, "name", "n", "", "the package name of PRIVATE PGP key to use")
	c.cmd.Run = c.Run
	return c
}

func (b *PGPDecryptCmd) Run(cmd *cobra.Command, args []string) {
	var (
		pgp  *crypto.PGP
		file []byte
		err  error
	)
	if len(args) == 0 {
		core.RaiseErr("the name of the file to decrypt is required")
	} else if len(args) > 1 {
		core.RaiseErr("only the name of the file to decrypt is required")
	}
	path := core.ToAbs(args[0])
	// the file to decrypt must have the .asc extension
	if !strings.HasSuffix(filepath.Ext(path), "asc") {
		core.RaiseErr("decrypt can only process files with .asc extensions (OpenPGP Armor ASCII message format)")
	}
	if len(b.keyPath) > 0 {
		// load the crypto key
		pgp, err = crypto.LoadPGP(core.ToAbs(b.keyPath), "")
		i18n.Err(err, i18n.ERR_CANT_LOAD_PRIV_KEY)
	} else
	// load the key based on the local repository resolution process
	{
		pgp, err = crypto.LoadPGPPrivateKey(b.group, b.name)
		i18n.Err(err, i18n.ERR_CANT_LOAD_PRIV_KEY)
	}
	// check the key file provided has a private key
	if !pgp.HasPrivate() {
		core.RaiseErr("the provided key file does not contain a private key, cannot decrypt")
	}
	// read the file to encrypt
	file, err = ioutil.ReadFile(path)
	core.CheckErr(err, "cannot load file to decrypt: %s", path)
	// decrypt the file content
	decrypted, err := pgp.Decrypt(file)
	core.CheckErr(err, "cannot decrypt file: %s", path)
	// save the encrypted file as *.asc (OpenPGP Armor ASCII)
	core.CheckErr(ioutil.WriteFile(strings.TrimSuffix(path, filepath.Ext(path)), decrypted, os.ModePerm), "cannot write decrypted file")
}
