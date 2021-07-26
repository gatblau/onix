package cmd

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/spf13/cobra"
)

// FlowCmd provides functions to manage Artisan execution flows
type FlowCmd struct {
	cmd *cobra.Command
}

func NewFlowCmd() *FlowCmd {
	c := &FlowCmd{
		cmd: &cobra.Command{
			Use:   "flow",
			Short: "provides functions to manage Artisan execution flows",
			Long:  `provides functions to manage Artisan execution flows`,
		},
	}
	return c
}
