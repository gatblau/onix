/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"github.com/spf13/cobra"
	"log"
)

// list local artefacts
type PushCmd struct {
	cmd  *cobra.Command
	repo *core.LocalRegistry
}

func NewPushCmd() *PushCmd {
	c := &PushCmd{
		cmd: &cobra.Command{
			Use:   "artie push [OPTIONS] NAME[:TAG]",
			Short: "Push an artifact to a registry",
			Long:  ``,
		},
		repo: core.NewRepository(),
	}
	c.cmd.Run = c.Run
	return c
}

func (b *PushCmd) Run(cmd *cobra.Command, args []string) {
	nameTag := args[0]
	named, err := core.ParseNormalizedNamed(nameTag)
	if err != nil {
		log.Fatal(err)
	}
	b.repo.GetArtefactsByRepo(named.String())
	log.Print(fmt.Sprintf("found artefact %s in local registry", named))
}
