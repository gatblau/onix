/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import (
	"fmt"
	"strings"
)

// InRoute the definition of an inbound route
type InRoute struct {
	// Name the name of the route
	Name string `bson:"_id" json:"name" yaml:"name" example:"SUPPLIER_A_IN_ROUTE"`
	// Description a description indicating the purpose of the route
	Description string `bson:"description "json:"description" yaml:"description" example:"the inbound route for supplier A"`
	// ServiceHost the remote host from where inbound files should be downloaded
	ServiceHost string `bson:"service_host" json:"service_host" yaml:"service_host" example:"s3.supplier-a.com"`
	// ServiceId a unique identifier for the S3 service where inbound files should be downloaded
	ServiceId string `bson:"service_id" json:"service_id" yaml:"service_id"`
	// BucketName the name of the S3 bucket containing files to download
	BucketName string `bson:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
	// User the username to authenticate against the remote ServiceHost
	User string `bson:"user "json:"user" yaml:"user"`
	// Pwd the password to authenticate against the remote ServiceHost
	Pwd string `bson:"pwd" json:"pwd" yaml:"pwd"`
	// PublicKey the PGP public key used to verify the author of the downloaded files
	PublicKey string `bson:"public_key" json:"public_key" yaml:"public_key"`
	// Verify a flag indicating whether author verification should be enabled
	Verify bool `bson:"verify" json:"verify" yaml:"verify"`
	// WebhookToken an authentication token to be passed by an event sender to be authenticated by the doorman's proxy webhook
	// its value can be anything, but it is typically a base64 encoded global unique identifier
	WebhookToken string `bson:"webhook_token" json:"webhook_token" yaml:"webhook_token" example:"JFkxnsn++02UilVkYFFC9w=="`
	// WebhookWhitelist the list of IP addresses accepted by the webhook (whitelist)
	WebhookWhitelist []string `bson:"webhook_whitelist" json:"webhook_whitelist" yaml:"webhook_whitelist"`
	// Filter a regular expression to filter publication events and prevent doorman from being invoked
	// if not defined, no filter is applied
	Filter string `bson:"filter" json:"filter" yaml:"filter"`
}

func (r InRoute) GetName() string {
	return r.Name
}

func (r InRoute) Valid() error {
	if len(r.ServiceHost) == 0 {
		return fmt.Errorf("inbound route %s service_host is mandatory", r.Name)
	}
	if !strings.HasPrefix(r.ServiceHost, "s3") {
		return fmt.Errorf("inbound route %s invalid service_host, should start with s3:// or s3s://", r.Name)
	}
	if (len(r.User) > 0 && len(r.Pwd) == 0) || (len(r.User) == 0 && len(r.Pwd) > 0) {
		return fmt.Errorf("inbound route %s requires both username and password to be provided, or none of them", r.Name)
	}
	if r.Verify && len(r.PublicKey) == 0 {
		return fmt.Errorf("inbound route %s requires author verification so, it must specify the author's public key", r.Name)
	}
	return nil
}

func (r InRoute) Creds() string {
	return fmt.Sprintf("%s:%s", r.User, r.Pwd)
}
