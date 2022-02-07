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

// Key a digital key
// @Description a digital key used to sign or verify signatures
type Key struct {
	// a unique identifier for the digital key
	Name string `bson:"_id" json:"name"`
	// the name of the entity owning the key
	Owner string `bson:"owner" json:"owner"`
	// a description of the intended use of the key
	Description string `bson:"description" json:"description"`
	// the actual content of the key
	Value string `bson:"value" json:"value"`
	// indicates if the key is private, otherwise public
	IsPrivate bool `bson:"is_private" json:"is_private"`
}

func (k Key) Valid() error {
	if len(k.Name) == 0 {
		return fmt.Errorf("name attribute must be prodided")
	}
	// ensure name has no blank spaces and is in uppercase
	k.Name = strings.ReplaceAll(strings.ToUpper(k.Name), " ", "_")
	if len(k.Owner) == 0 {
		return fmt.Errorf("owner attribute must be prodided")
	}
	if len(k.Value) == 0 {
		return fmt.Errorf("value attribute must be prodided")
	}
	return nil
}

func (k Key) GetName() string {
	return k.Name
}
