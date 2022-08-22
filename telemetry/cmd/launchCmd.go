/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/gatblau/onix/telemetry/opentelemetry"
	"github.com/spf13/cobra"
)

// LaunchCmd launches host pilot
type LaunchCmd struct {
	cmd         *cobra.Command
	configPaths string // location of telemetry config yaml file
	version     string // version
}

func NewLaunchCmd() *LaunchCmd {
	c := &LaunchCmd{
		cmd: &cobra.Command{
			Use:   "launch [flags]",
			Short: "launches telemetry",
			Long:  `launches telemetry`,
		},
	}
	c.cmd.Flags().StringVarP(&c.configPaths, "configpaths", "c", "", "location of telemetry config yaml file")
	c.cmd.Flags().StringVarP(&c.version, "version", "v", "", "version")

	c.cmd.Run = c.Run
	return c
}

func (c *LaunchCmd) Run(cmd *cobra.Command, args []string) {
	var cfg []string
	if len(c.configPaths) > 0 {
		cfg = strings.Split(c.configPaths, ",")
	} else {
		defaultpath := string("telem.yaml")
		cfg = append(cfg, defaultpath)
	}

	collector := opentelemetry.NewOpenTelemetry(cfg, c.version, nil)
	err := collector.Run(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	status := <-collector.Status()
	if !status.Running {
		log.Fatal("service should be running")
	}
	if status.Err != nil {
		log.Fatal(status.Err.Error())
	}
	collector.Stop()
	status = <-collector.Status()
	if status.Running {
		log.Fatal("service should be stopped")
	}
}
