/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// UtilPwdCmd generates passwords
type UtilPwdCmd struct {
	Cmd          *cobra.Command
	len          *int
	specialChars *bool
	bcrypt       *bool
}

func NewUtilPwdCmd() *UtilPwdCmd {
	c := &UtilPwdCmd{
		Cmd: &cobra.Command{
			Use:   "pwd [flags]",
			Short: "generates a random password",
			Long:  `generates a random password`,
		},
	}
	c.len = c.Cmd.Flags().IntP("length", "l", 16, "length of the generated password")
	c.specialChars = c.Cmd.Flags().BoolP("special-chars", "s", false, "use special characters in the generated password")
	c.bcrypt = c.Cmd.Flags().BoolP("bcrypt", "b", false, "hash it using bcrypt algorithm")
	c.Cmd.Run = c.Run
	return c
}

func (c *UtilPwdCmd) Run(_ *cobra.Command, _ []string) {
	pwd := core.RandomPwd(*c.len, *c.specialChars)
	fmt.Printf("%s", pwd)
	if *c.bcrypt {
		hash, err := pwdHash([]byte(pwd))
		core.CheckErr(err, "cannot hash password")
		fmt.Printf("\nbcrypt:%s\n", hash)
	}
}

func pwdHash(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
