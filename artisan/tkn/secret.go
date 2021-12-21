/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

type Secret struct {
	APIVersion string             `yaml:"apiVersion,omitempty"`
	Kind       string             `yaml:"kind,omitempty"`
	Metadata   *Metadata          `yaml:"metadata,omitempty"`
	Type       string             `yaml:"type,omitempty"`
	StringData *map[string]string `yaml:"stringData,omitempty"`
	SecretName string             `yaml:"secretName,omitempty"`
}

// type StringData struct {
// 	Pwd   string `yaml:"pwd,omitempty"`
// 	User  string `yaml:"user,omitempty"`
// 	Token string `yaml:"token,omitempty"`
// }
