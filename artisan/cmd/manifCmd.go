/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"github.com/yalp/jsonpath"
)

// return artefact's manifest
type ManifestCmd struct {
	cmd    *cobra.Command
	filter string
}

func NewManifestCmd() *ManifestCmd {
	c := &ManifestCmd{
		cmd: &cobra.Command{
			Use:   "manifest [flags] name:tag",
			Short: "returns the artefact manifest",
			Long:  ``,
		},
	}
	c.cmd.Flags().StringVarP(&c.filter, "filter", "f", "", "--filter=JSONPath or -f=JSONPath")
	c.cmd.Run = c.Run
	return c
}

func (b *ManifestCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		core.RaiseErr("the artefact name:tag is required")
	} else if len(args) > 1 {
		core.RaiseErr("too many arguments")
	}
	// create a local registry
	local := registry.NewLocalRegistry()
	// get the artefact manifest
	m := local.GetManifest(core.ParseName(args[0]))
	// marshal the manifest
	bytes, err := json.MarshalIndent(m, "", "  ")
	core.CheckErr(err, "cannot marshal manifest")
	// if no filter is set then return the whole manifest
	if len(b.filter) == 0 {
		fmt.Printf("%v\n", string(bytes))
	} else {
		var jason interface{}
		err := json.Unmarshal(bytes, &jason)
		// otherwise apply the jsonpath to extract a value from the manifest
		result, err := jsonpath.Read(jason, b.filter)
		core.CheckErr(err, "cannot apply filter expression '%s'", b.filter)
		fmt.Printf("%v", result)
	}
}
