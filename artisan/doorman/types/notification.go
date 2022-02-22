/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

type Notification struct {
	// unique identifier for the notification
	Name string `yaml:"name" json:"name" bson:"_id"`
	// Recipient of the notification if type is email
	Recipient string `yaml:"recipient" json:"recipient" bson:"recipient"`
	// Type of the notification (e.g. email, snow, etc.)
	Type string `yaml:"type" json:"type" bson:"type"`
	// Template to use for content of the notification
	Template string `yaml:"template" json:"template" bson:"template"`
}

func (n Notification) GetName() string {
	return n.Name
}

func (n Notification) Valid() error {
	return nil
}

type NotificationTemplate struct {
	// Name unique identifier for notification template
	Name string `yaml:"name" json:"name" bson:"_id"`
	// Subject of the notification
	Subject string `yaml:"subject" json:"subject" bson:"subject"`
	// Content of the template
	Content string `yaml:"content" json:"content" bson:"content"`
}

func (t NotificationTemplate) GetName() string {
	return t.Name
}

func (t NotificationTemplate) Valid() error {
	return nil
}
