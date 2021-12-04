/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/app"
	"github.com/gatblau/onix/artisan/core"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// AppCmd creates deployment config files for a target platform
type AppCmd struct {
	cmd    *cobra.Command
	creds  string
	path   string
	format string
}

func NewAppCmd() *AppCmd {
	c := &AppCmd{
		cmd: &cobra.Command{
			Use:   "app [flags] [app manifest URI]",
			Short: "generates application deployment configuration files for a target platform",
			Long: `generates application deployment configuration files for a target platform\n
the URI can be either http(s):// or file://`,
			Args: cobra.ExactArgs(1),
		},
	}
	c.cmd.Flags().StringVarP(&c.creds, "creds", "u", "", "-u user:password; the credentials for git authentication")
	c.cmd.Flags().StringVarP(&c.format, "format", "f", "compose", "the target format for the configuration files (i.e. compose or k8s)")
	c.cmd.Flags().StringVarP(&c.path, "path", "p", ".", "the output path where the configuration files will be written")
	c.cmd.Run = c.Run
	return c
}

func (c *AppCmd) Run(cmd *cobra.Command, args []string) {
	uri := args[0]
	manifest, err := app.NewAppMan(uri)
	if err != nil {
		core.ErrorLogger.Fatalf(err.Error())
	}
	var builderType app.BuilderType
	switch strings.ToLower(c.format) {
	case "compose":
		builderType = app.DockerCompose
	case "k8s":
		builderType = app.Kubernetes
	default:
		core.ErrorLogger.Fatalf("invalid format, valid formats are compose or k8s")
	}
	builder, err := app.NewBuilder(builderType, *manifest)
	if err != nil {
		core.ErrorLogger.Fatalf(err.Error())
	}
	files, err := builder.Build()
	if err != nil {
		core.ErrorLogger.Fatalf(err.Error())
	}
	path, err := filepath.Abs(c.path)
	if err != nil {
		core.ErrorLogger.Fatalf(err.Error())
	}
	// ensure path exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			core.ErrorLogger.Fatalf("cannot create folder '%s': %s\n", path, err)
		}
	}
	for _, file := range files {
		fpath := filepath.Join(path, file.Name)
		err = os.WriteFile(fpath, file.Content, os.ModePerm)
		if err != nil {
			core.ErrorLogger.Fatalf("cannot write file %s: %s\n", fpath, err)
		}
	}
}
