//   Onix Config DatabaseProvider - Dbman
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

type RootCmd struct {
	*cobra.Command
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		&cobra.Command{
			Use:   "dbman",
			Short: "database manager",
			Long: `dbman is a CLI tool to manage database schema released versions, upgrade data and perform database backups and restores.
	dbman is part of (and used by) Onix Configuration DatabaseProvider (see https://onix.gatblau.org) to manage its configuration database.
	dbman can also be run from a container (when in http mode) to manage the data / schema life cycle of databases from a container platform.`,
		},
	}
	cobra.OnInitialize(c.initConfig)
	return c
}

// initConfig reads in config file and ENV variables if set.
func (c *RootCmd) initConfig() {
	dm, err := util.NewDbMan()
	if err != nil {
		log.Err(err).Msg("cannot create DbMan instance")
	}
	util.DM = dm
}
