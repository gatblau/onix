/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "fmt"

// OutRoute represents an outbound route to distribute packages and images
type OutRoute struct {
	// Name the name uniquely identifying the outbound route
	Name string `bson:"_id" json:"name" example:"ACME_OUT_LOGISTICS"`
	// Description describes the purpose of the route
	Description     string           `bson:"description" json:"description" example:"outbound route for ACME company logistics department"`
	PackageRegistry *PackageRegistry `bson:"package_registry" json:"package_registry"`
	ImageRegistry   *ImageRegistry   `bson:"image_registry" json:"image_registry"`
}

func (r OutRoute) GetName() string {
	return r.Name
}

func (r OutRoute) Valid() error {
	if r.PackageRegistry == nil && r.ImageRegistry == nil {
		return fmt.Errorf("outbound route %s must specify at least one registry", r.Name)
	}
	if r.PackageRegistry != nil {
		if len(r.PackageRegistry.URI) == 0 {
			return fmt.Errorf("inbound route %s requires package registry URI", r.Name)
		}
		if (len(r.PackageRegistry.User) > 0 && len(r.PackageRegistry.Pwd) == 0) || (len(r.PackageRegistry.User) == 0 && len(r.PackageRegistry.Pwd) > 0) {
			return fmt.Errorf("outbound route %s: package registry requires both username and password to be provided, or none of them", r.Name)
		}
		if r.PackageRegistry.Sign && len(r.PackageRegistry.PrivateKey) == 0 {
			return fmt.Errorf("outbound route %s requires signature so, it must specify the signer's private key", r.Name)
		}
	}
	if r.ImageRegistry != nil {
		if len(r.ImageRegistry.URI) == 0 {
			return fmt.Errorf("outbound route %s requires image registry URI", r.Name)
		}
		if (len(r.ImageRegistry.User) > 0 && len(r.ImageRegistry.Pwd) == 0) || (len(r.ImageRegistry.User) == 0 && len(r.ImageRegistry.Pwd) > 0) {
			return fmt.Errorf("outbound route %s: image registry requires both username and password to be provided, or none of them", r.Name)
		}
	}
	return nil
}

// PackageRegistry the details of the target package registry within an outbound route
type PackageRegistry struct {
	// URI the location of the package registry
	URI string `bson:"uri" json:"uri" example:"packages.acme.com:8082"`
	// User the username to authenticate with the package registry
	User string `bson:"user" json:"user" example:"test_user"`
	// Pwd the password to authenticate with the package registry
	Pwd string `bson:"pwd" json:"pwd" example:"d8y2b9fc97y23!$^"`
	// Sign a flag indicating whether packages pushed to the registry should be resigned
	Sign bool `bson:"sign" json:"sign" example:"true"`
	// PrivateKey the name of the private PGP key used to re-sign the packages
	PrivateKey string `bson:"private_key" json:"private_key" example:"SIGNING_KEY_01"`
}

// ImageRegistry the details of the target registry within an outbound route
type ImageRegistry struct {
	// URI the location of the container image registry
	URI string `bson:"uri" json:"uri" example:"images.acme.com:5000"`
	// User the username to authenticate with the container image registry
	User string `bson:"user" json:"user"`
	// Pwd the password to authenticate with the container image registry
	Pwd string `bson:"pwd" json:"pwd"`
}
