/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"github.com/gatblau/onix/artisan/doorman/db"
	"github.com/gatblau/onix/artisan/doorman/types"
)

type Processor interface {
	// Start the processing of an event
	Start()

	// Pipeline executes a pipeline
	Pipeline(pipe *types.Pipeline) error

	// InboundRoute executes an inbound route in a pipeline
	InboundRoute(pipe *types.Pipeline, route types.InRoute) error

	// PreImport runs pre import checks
	PreImport(route types.InRoute, err error) error

	// Command executes the specified command
	Command(command types.Command) error

	// OutboundRoute execute the outbound route
	OutboundRoute(outRoute types.OutRoute) error

	// ExportFiles from a specification to an S3 bucket
	ExportFiles(s3Store *types.S3Store) error

	// PushImages in a specification to a container registry
	PushImages(imageRegistry *types.ImageRegistry) error

	// PushPackages in a specification to an artisan registry
	PushPackages(pkgRegistry *types.PackageRegistry) error

	// ImportFiles from a specification
	ImportFiles() error

	// SendNotification to the notification service
	SendNotification(nType db.NotificationType) error

	// BeforeComplete executes tasks at the end of the pipeline process
	BeforeComplete(pipe *types.Pipeline) error

	// Info logger
	Info(format string, a ...interface{})

	// Error logger
	Error(format string, a ...interface{}) error

	// Warn logger
	Warn(format string, a ...interface{})
}
