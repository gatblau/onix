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
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SaveCmd save one or more packages from the local registry to a tar archive to allow copying without using registries
type SaveCmd struct {
	cmd         *cobra.Command
	srcCreds    string
	targetCreds string
	output      string
}

func NewSaveCmd() *SaveCmd {
	c := &SaveCmd{
		cmd: &cobra.Command{
			Use:   "save [OPTIONS] PACKAGE [flags]",
			Short: "Save one or more packages to a tar archive",
			Long: `Usage:  art save [OPTIONS] PACKAGE [PACKAGE...]

Save one or more packages to a tar archive, streamed to STDOUT by default or saved to a URI (file system or S3)

Options:
  -o, --output        string   Write to a URI, instead of STDOUT 
                               URI can be file system (absolute or relative path) or s3 bucket (s3:// or s3s:// using TLS)
  -u, --source-creds  string   The credentials used to pull packages from an authenticated artisan registry, if the packages are not already in the local registry
  -v  --target-creds  string   The credentials to write packages to a destination, if such destination implements authentication (e.g. s3)

Examples:
   # save package1 and package2 to a tar archive by redirecting STDOUT to file (using >)
   art save package1 package2 > archive.tar 
   
   # save package1 and package2 to a tar archive by specifying relative file path via URI (using -o)
   art save package1 package2 -o ./test/archive.tar 

   # save package1 and package2 from remote artisan registry to an authenticated and TLS enabled s3 bucket
   art save package1 package2 -u reg-USER:reg-PWD -o s3s://endpoint/bucket/archive.tar -c s3-ID:s3-SECRET
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.output, "output", "o", "", "-o exported/archive.tar; the output where the archive will be written including filename")
	c.cmd.Flags().StringVarP(&c.srcCreds, "user", "u", "", "-u USER:PASSWORD; artisan registry username and password")
	c.cmd.Flags().StringVarP(&c.targetCreds, "creds", "c", "", "-c USER:PASSWORD; destination URI username and password")
	return c
}

func (c *SaveCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if len(args) < 1 {
		log.Fatal("at least the name of one package to save is required")
	}
	// validate the package names
	names, err := core.ValidateNames(args)
	i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
	// create a local registry
	local := registry.NewLocalRegistry()
	// export packages into tar bytes
	content, err := local.Save(names, c.srcCreds)
	core.CheckErr(err, "cannot export package(s)")
	if len(c.output) == 0 {
		fmt.Print(string(content[:]))
	} else {
		targetPath := c.output
		// if the path does not implement an URI scheme (i.e. is a file path)
		if !strings.Contains(c.output, "://") {
			targetPath, err = filepath.Abs(targetPath)
			core.CheckErr(err, "cannot obtain the absolute output path")
			ext := filepath.Ext(targetPath)
			if len(ext) == 0 || ext != ".tar" {
				core.RaiseErr("output path must contain a filename with .tar extension")
			}
			// creates target directory
			core.CheckErr(os.MkdirAll(filepath.Dir(targetPath), 0755), "cannot create target output folder")
		}
		core.CheckErr(core.WriteFile(content, targetPath, c.targetCreds), "cannot save exported package file")
	}
}
