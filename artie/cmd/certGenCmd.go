/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"github.com/gatblau/onix/artie/sign"
	"github.com/spf13/cobra"
)

// list local artefacts
type CertGenCmd struct {
	cmd     *cobra.Command
	bitSize *int   // the key bit size
	path    string // the file path to the keys
	name    string // the name of the keys
}

func NewCertGenCmd() *CertGenCmd {
	c := &CertGenCmd{
		cmd: &cobra.Command{
			Use:   "gen",
			Short: "generate a new RSA key pair",
			Long:  `RSA keys are used to sign and verify signatures of artefacts`,
		},
	}
	c.bitSize = c.cmd.Flags().IntP("size", "s", 2048, "The bit size of the generated RSA key pair, defaults to s=2048 \nOther common sizes are 1024, 3072 and 4096. \nAny size is possible.")
	c.cmd.Flags().StringVarP(&c.path, "path", "p", ".", "The path of the generated RSA key pair, defaults to the current directory \".\"")
	c.cmd.Flags().StringVarP(&c.name, "name", "n", "id", "The name given to the generated RSA key pair, defaults to id_rsa_key.pem (private key) and id_rsa_pub.pem (public key).\nIf specified, the naming is [name]_rsa_key.pem (private key) and [name]_rsa_pub.pem (public key)")
	c.cmd.Run = c.Run
	return c
}

func (b *CertGenCmd) Run(cmd *cobra.Command, args []string) {
	sign.GenerateKeys(b.path, b.name, *b.bitSize)
}
