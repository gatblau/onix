/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

// WebhookAuthInfo the information returned to doorman's proxy so that it can authenticate the Webhook caller
type WebhookAuthInfo struct {
	// ReferrerURL the URL the referrer of the token should come from to be considered valid
	// this is the URL of the S3 bucket where the spec files were uploaded
	ReferrerURL string
	// Whitelist one or more IP addresses that are considered valid for authentication of the Webhook caller
	// this is the real IP of the webhook caller
	Whitelist []string `json:"whitelist" yaml:"whitelist"`
}
