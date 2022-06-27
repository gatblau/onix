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
	. "github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
	"strings"
)

// SpecPushCmd downloads the contents of a spec from a remote source
type SpecPushCmd struct {
	cmd    *cobra.Command
	images bool
	tag    string
	clean  bool
	user   string
	creds  string
	logout bool
}

func NewSpecPushCmd() *SpecPushCmd {
	c := &SpecPushCmd{
		cmd: &cobra.Command{
			Use:   "push [FLAGS] SPEC-FILE",
			Short: "tag and pushes packages or images defined in the spec file to the tagged registry",
			Long: `tag and pushes packages or images defined in the spec file to the tagged registry
Usage: art spec push [FLAGS] SPEC-FILE

Example:
   # tag and push packages to tagged registry
   art spec push ./my-release -t package.registry.io/apps -u reg_user:reg_pwd

   # tag and push images to tagged registry (assumed already logged in container registry)
   art spec push ./my-release -i -t container.registry.io/apps
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.user, "user", "u", "", "the credentials used to push to the artisan registry (or docker registry if -i is used)")
	c.cmd.Flags().StringVarP(&c.tag, "tag", "t", "", "the target registry host and optionally user/group (e.g. <host>/<group>)")
	c.cmd.Flags().BoolVarP(&c.images, "images", "i", false, "if defined, the command applies to images instead of packages")
	c.cmd.Flags().BoolVar(&c.logout, "logout", false, "if defined, logs out the docker login session")
	c.cmd.Flags().BoolVar(&c.clean, "clean", false, "if defined, remove packages / images from local registries")
	c.cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials to retrieve the spec file from a remote destination")
	c.cmd.MarkFlagRequired("tag")
	return c
}

func (c *SpecPushCmd) Run(cmd *cobra.Command, args []string) {
	if args == nil {
		args = []string{"."}
	}
	if args != nil && len(args) < 1 {
		core.RaiseErr("the URI of the specification is required")
	}
	if len(c.tag) == 0 {
		core.RaiseErr("a tag value must be provided")
	}
	tagParts := strings.Split(c.tag, "/")
	host := tagParts[0]
	group := ""
	if len(tagParts) > 0 {
		group = strings.Join(tagParts[1:], "/")
	}
	core.CheckErr(PushSpec(
		PushOptions{
			SpecPath: args[0],
			Host:     host,
			Group:    group,
			User:     c.user,
			Creds:    c.creds,
			Image:    c.images,
			Clean:    c.clean,
			Logout:   c.logout,
		}), "cannot push spec")
}
