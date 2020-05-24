//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package cmd

import (
	"github.com/gatblau/onix/dbman/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// the path to dbman's config file
var cfgPath string

type RootCmd struct {
	*cobra.Command
	cfg *util.Config
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		Command: &cobra.Command{
			Use:   "dbman",
			Short: "database manager",
			Long: `DbMan is a CLI tool to manage database schema released versions, upgrade data and perform database backups and restores.
	DbMan is part of (and used by) Onix Configuration Manager (see https://onix.gatblau.org) to manage its configuration database.
	DbMan can also be run from a container (when in http mode) to manage the data / schema life cycle of databases from a container platform.`,
		}}
	c.init()
	return c
}

func (c *RootCmd) init() {
	cobra.OnInitialize(c.initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// bind the config-path flag
	c.PersistentFlags().StringVar(&cfgPath, "config-path", "", "path to dbman's config file (default is $HOME)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func (c *RootCmd) initConfig() {
	// if cfgPath is empty, then NewConfig will default to $HOME path
	cfg, err := util.NewConfig(cfgPath)
	if err != nil {
		log.Err(err).Msg("cannot initialise configuration file")
	}
	c.cfg = cfg
}
