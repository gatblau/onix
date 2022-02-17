/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "fmt"

// InRoute the definition of an inbound route
type InRoute struct {
	// Name the name of the route
	Name string `bson:"_id" json:"name" yaml:"name" example:"SUPPLIER_A_IN_ROUTE"`
	// Description a description indicating the purpose of the route
	Description string `bson:"description "json:"description" yaml:"description" example:"the inbound route for supplier A"`
	// URI the remote URI from where inbound files should be downloaded
	URI string `bson:"uri" json:"uri" yaml:"uri" example:"s3.supplier-a.com"`
	// User the username to authenticate against the remote URI
	User string `bson:"user "json:"user" yaml:"user"`
	// Pwd the password to authenticate against the remote URI
	Pwd string `bson:"pwd" json:"pwd" yaml:"pwd"`
	// PublicKey the PGP public key used to verify the author of the downloaded files
	PublicKey string `bson:"public_key" json:"public_key" yaml:"public_key"`
	// Verify a flag indicating whether author verification should be enabled
	Verify bool `bson:"verify" json:"verify" yaml:"verify"`
}

func (r InRoute) GetName() string {
	return r.Name
}

func (r InRoute) Valid() error {
	if len(r.URI) == 0 {
		return fmt.Errorf("inbound route %s URI is mandatory", r.Name)
	}
	if (len(r.User) > 0 && len(r.Pwd) == 0) || (len(r.User) == 0 && len(r.Pwd) > 0) {
		return fmt.Errorf("inbound route %s requires both username and password to be provided, or none of them", r.Name)
	}
	if r.Verify && len(r.PublicKey) == 0 {
		return fmt.Errorf("inbound route %s requires author verification so, it must specify the author's public key", r.Name)
	}
	return nil
}
