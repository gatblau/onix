/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "fmt"

// PipelineConf a pipeline configuration used to connect inbound and outbound routes
type PipelineConf struct {
	// Name the name uniquely identifying the pipeline
	Name string `bson:"_id" json:"name" yaml:"name" example:"ACME_PIPELINE"`
	// InboundRoute  the name of the inbound route to use in the pipeline
	InboundRoute string `json:"in_route" yaml:"inbound_route" bson:"inbound_route"`
	// OutboundRoute  the name of the outbound route to use in the pipeline
	OutboundRoute string `json:"out_route" yaml:"outbound_route" bson:"outbound_route"`
	// Commands a list of the command names to be executed between inbound and outbound routes
	Commands []string `json:"commands" yaml:"commands" bson:"commands"`
	// SuccessNotification notification to use in case of success
	SuccessNotification string `json:"success_notification" yaml:"success_notification" bson:"success_notification"`
	// ErrorNotification notification to use in case of errors
	ErrorNotification string `json:"error_notification" yaml:"error_notification" bson:"error_notification"`
	// CmdFailedNotification notification to use in case of command failure
	CmdFailedNotification string `json:"cmd_failed_notification" yaml:"cmd_failed_notification" bson:"cmd_failed_notification"`
}

func (r PipelineConf) GetName() string {
	return r.Name
}

func (r PipelineConf) Valid() error {
	if len(r.InboundRoute) == 0 {
		return fmt.Errorf("pipeline %s must define an inbound route", r.Name)
	}
	if len(r.OutboundRoute) == 0 {
		return fmt.Errorf("pipeline %s must define an outbound route", r.Name)
	}
	return nil
}

// Pipeline provides all information for one pipeline
type Pipeline struct {
	// Name the name uniquely identifying the pipeline
	Name string `json:"name"`
	// InboundRoute  the name of the inbound route to use in the pipeline
	InboundRoute InRoute `json:"in_route"`
	// OutboundRoute  the name of the outbound route to use in the pipeline
	OutboundRoute OutRoute `json:"out_route"`
	// Commands a list of the command names to be executed between inbound and outbound routes
	Commands []string `json:"commands"`
}
