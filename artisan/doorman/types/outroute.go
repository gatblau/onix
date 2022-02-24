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

// OutRoute represents an outbound route to distribute packages and images
type OutRoute struct {
	// Name the name uniquely identifying the outbound route
	Name string `bson:"_id" json:"name" example:"ACME_OUT_LOGISTICS"`
	// Description describes the purpose of the route
	Description string `bson:"description" json:"description" yaml:"description" example:"outbound route for ACME company logistics department"`
	// PackageRegistry the information about the artisan registry that is the destination for the spec packages
	PackageRegistry *PackageRegistry `bson:"package_registry" json:"package_registry" yaml:"package_registry"`
	// ImageRegistry the information about the image registry that is the destination for the spec images
	ImageRegistry *ImageRegistry `bson:"image_registry" json:"image_registry" yaml:"image_registry"`
	// S3Store the information about the S3 service that is the destination for the spec tarball files
	S3Store *S3Store `bson:"s3_store" json:"s3_store" yaml:"s3_store"`
}

func (r OutRoute) GetName() string {
	return r.Name
}

func (r OutRoute) Valid() error {
	if r.PackageRegistry == nil && r.ImageRegistry == nil {
		return fmt.Errorf("outbound route %s must specify at least one registry", r.Name)
	}
	if r.PackageRegistry != nil {
		if len(r.PackageRegistry.Domain) == 0 {
			return fmt.Errorf("inbound route %s requires package registry Domain", r.Name)
		}
		if strings.Contains(r.PackageRegistry.Domain, "/") {
			return fmt.Errorf("package registry Domain for outbound route %s must not have / only root domain (and potentially port)", r.Name)
		}
		if (len(r.PackageRegistry.User) > 0 && len(r.PackageRegistry.Pwd) == 0) || (len(r.PackageRegistry.User) == 0 && len(r.PackageRegistry.Pwd) > 0) {
			return fmt.Errorf("outbound route %s: package registry requires both username and password to be provided, or none of them", r.Name)
		}
		if r.PackageRegistry.Sign && len(r.PackageRegistry.PrivateKey) == 0 {
			return fmt.Errorf("outbound route %s requires signature so, it must specify the signer's private key", r.Name)
		}
	}
	if r.ImageRegistry != nil {
		if len(r.ImageRegistry.Domain) == 0 {
			return fmt.Errorf("outbound route %s requires image registry Domain", r.Name)
		}
		if strings.Contains(r.ImageRegistry.Domain, "/") {
			return fmt.Errorf("image registry Domain for outbound route %s must not have / only root domain (and potentially port)", r.Name)
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
	Domain string `bson:"domain" json:"domain" yaml:"domain" example:"packages.acme.com:8082"`
	// Group the group (location withing the registry) where the packages should be placed
	// if not specified, the group from the package to push is used
	Group string `bson:"group" json:"group" yaml:"group" example:"test/groupA"`
	// User the username to authenticate with the package registry
	User string `bson:"user" json:"user" yaml:"user" example:"test_user"`
	// Pwd the password to authenticate with the package registry
	Pwd string `bson:"pwd" json:"pwd" yaml:"pwd" example:"d8y2b9fc97y23!$^"`
	// Sign a flag indicating whether packages pushed to the registry should be resigned
	Sign bool `bson:"sign" json:"sign" yaml:"sign" example:"true"`
	// PrivateKey the name of the private PGP key used to re-sign the packages
	PrivateKey string `bson:"private_key" json:"private_key" yaml:"private_key" example:"SIGNING_KEY_01"`
}

// ImageRegistry the details of the target registry within an outbound route
type ImageRegistry struct {
	// URI the location of the container image registry
	Domain string `bson:"domain" json:"domain" yaml:"domain" example:"images.acme.com:5000"`
	// Group the group (location withing the registry) where the packages should be placed
	// if not specified, the group from the package to push is used
	Group string `bson:"group" json:"group" yaml:"group" example:"test/groupA"`
	// User the username to authenticate with the container image registry
	User string `bson:"user" json:"user" yaml:"user"`
	// Pwd the password to authenticate with the container image registry
	Pwd string `bson:"pwd" json:"pwd" yaml:"pwd"`
}

// S3Store the details of the target S3 store within an outbound route
type S3Store struct {
	// BucketURI the URI of the folder where to upload the spec tar files
	BucketURI string `bson:"bucket_uri" json:"bucket_uri" yaml:"bucket_uri"`
	// User the username of the outbound S3 bucket
	User string `bson:"user" json:"user" yaml:"user"`
	// Pwd the password of the outbound S3 bucket
	Pwd string `bson:"pwd" json:"pwd" yaml:"pwd"`
	// Sign a flag indicating whether packages pushed to the S3 service should be resigned
	Sign bool `bson:"sign" json:"sign" yaml:"sign" example:"true"`
	// PrivateKey the name of the private PGP key used to re-sign the packages in the tarball files
	PrivateKey string `bson:"private_key" json:"private_key" yaml:"private_key" example:"SIGNING_KEY_01"`
}

func (s S3Store) Creds() string {
	return fmt.Sprintf("%s:%s", s.User, s.Pwd)
}
