/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// list local artefacts
type PGPEncryptCmd struct {
	cmd     *cobra.Command
	keyPath string
	group   string
	name    string
}

func NewPGPEncryptCmd() *PGPEncryptCmd {
	c := &PGPEncryptCmd{
		cmd: &cobra.Command{
			Use:   "encrypt [flags] filename",
			Short: "encrypts a file using a designated PGP public key",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.keyPath, "key", "k", "", "the path to the private key to use")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "the artefact group of the private key to use")
	c.cmd.Flags().StringVarP(&c.name, "name", "n", "", "the artefact name of the private key to use")
	c.cmd.Run = c.Run
	return c
}

func (b *PGPEncryptCmd) Run(cmd *cobra.Command, args []string) {
	var (
		pgp  *crypto.PGP
		file []byte
		err  error
	)
	if len(args) == 0 {
		core.RaiseErr("the name of the file to encrypt is required")
	} else if len(args) > 1 {
		core.RaiseErr("only the name of the file to encrypt is required")
	}
	path := core.ToAbs(args[0])
	if len(b.keyPath) > 0 {
		// load the crypto key
		pgp, err = crypto.LoadPGP(core.ToAbs(b.keyPath))
		core.CheckErr(err, "cannot load public key")
	} else
	// load the key based on the local repository resolution process
	{
		pgp, err = crypto.LoadPGPPublicKey(b.group, b.name)
		core.CheckErr(err, "cannot load public key")
	}
	// read the file to encrypt
	file, err = ioutil.ReadFile(path)
	core.CheckErr(err, "cannot load file to encrypt: %s", path)
	// encrypt the file content
	encrypted, err := pgp.Encrypt(file)
	core.CheckErr(err, "cannot encrypt file: %s", path)
	// save the encrypted file as *.asc (OpenPGP Armor ASCII)
	core.CheckErr(ioutil.WriteFile(fmt.Sprintf("%s.asc", path), encrypted, os.ModePerm), "cannot write encrypted file")
}
