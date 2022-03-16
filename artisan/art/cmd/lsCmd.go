/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
)

// ListCmd list packages
type ListCmd struct {
	cmd      *cobra.Command
	quiet    bool
	registry string
	creds    string
}

func NewListCmd() *ListCmd {
	c := &ListCmd{
		cmd: &cobra.Command{
			Use:   "ls [FLAGS]",
			Short: "list packages in the local or a remote registry",
			Long:  `list packages in the local or a remote registry`,
			Example: `
# list packages in local registry
art ls

# list packages from remote registry at localhost:8082
art ls -r localhost:8082 -u <user>:<pwd>
`,
		},
	}
	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "only show numeric IDs")
	c.cmd.Flags().StringVarP(&c.registry, "registry", "r", "", "the domain name or IP of the remote registry (e.g. my-remote-registry); port can also be specified using a colon syntax")
	c.cmd.Flags().StringVarP(&c.creds, "user", "u", "", "the credentials used to retrieve the information from the remote registry")
	c.cmd.Run = c.Run
	return c
}

func (c *ListCmd) Run(cmd *cobra.Command, args []string) {
	if len(c.registry) == 0 {
		local := registry.NewLocalRegistry()
		if c.quiet {
			local.ListQ()
		} else {
			local.List()
		}
	}
	if len(c.registry) > 0 {
		uname, pwd := core.RegUserPwd(c.creds)
		remote, err := registry.NewRemoteRegistry(c.registry, uname, pwd)
		core.CheckErr(err, "invalid registry name")
		remote.List(c.quiet)
	}
}
