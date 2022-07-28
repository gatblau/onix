/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	ctl "github.com/gatblau/onix/pilotctl/types"
	pilotCore "github.com/gatblau/onix/piloth/core"
	"github.com/spf13/cobra"
)

// LaunchCmd launches host pilot
type LaunchCmd struct {
	cmd                *cobra.Command
	useHwId            bool   // use hardware uuid to identify device (instead of primary mac address)
	tracing            bool   // enables tracing
	logCollector       bool   // enables log collector
	cpu                bool   // enables cpu profiling
	mem                bool   // enables memory profiling
	insecureSkipVerify bool   // if true, crypto/tls accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to machine-in-the-middle attacks unless custom verification is used.
	cvePath            string // the  path used to collect CVE reports to export
}

func NewLaunchCmd() *LaunchCmd {
	c := &LaunchCmd{
		cmd: &cobra.Command{
			Use:   "launch [flags]",
			Short: "launches host pilot",
			Long:  `launches host pilot`,
		},
	}
	c.cmd.Flags().BoolVarP(&c.useHwId, "hw-id", "w", false, "use hardware uuid to identify device(instead of primary mac address)")
	c.cmd.Flags().BoolVarP(&c.tracing, "trace", "t", false, "enables tracing")
	c.cmd.Flags().BoolVarP(&c.logCollector, "syslog-collector", "s", false, "enables the syslog collector")
	c.cpu = *c.cmd.Flags().Bool("cpu", false, "enables cpu profiling only; cannot profile memory")
	c.mem = *c.cmd.Flags().Bool("mem", false, "enables memory profiling only; cannot profile cpu")
	c.insecureSkipVerify = *c.cmd.Flags().Bool("insecureSkipVerify", false, "disables verification of certificates presented by the server and host name in that certificate; in this mode, TLS is susceptible to machine-in-the-middle attacks unless custom verification is used.")
	c.cmd.Flags().StringVar(&c.cvePath, "cve-path", "", "if set, enables export of CVE reports from specified path")
	c.cmd.Run = c.Run
	return c
}

func (c *LaunchCmd) Run(cmd *cobra.Command, args []string) {
	// collects device/host information
	hostInfo, err := ctl.NewHostInfo()
	if err != nil {
		core.RaiseErr("cannot collect host information")
	}
	// creates pilot instance
	p, err := pilotCore.NewPilot(pilotCore.PilotOptions{
		UseHwId:            c.useHwId,
		Logs:               c.logCollector,
		Tracing:            c.tracing,
		Info:               hostInfo,
		CPU:                c.cpu,
		MEM:                c.mem,
		InsecureSkipVerify: c.insecureSkipVerify,
		CVEPath:            c.cvePath,
	})
	core.CheckErr(err, "cannot start pilot")
	// start the pilot
	p.Start()
}
