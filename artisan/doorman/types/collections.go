/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

type Collection string

const (
	KeysCollection                  Collection = "keys"
	CommandsCollection              Collection = "commands"
	InRouteCollection               Collection = "inbound-routes"
	OutRouteCollection              Collection = "outbound-routes"
	PipelineCollection              Collection = "pipelines"
	NotificationTemplatesCollection Collection = "templates"
	NotificationsCollection         Collection = "notifications"
)
