/*
   Onix Config Manager - OxTerra - Terraform Http Backend for Onix
   Copyright (c) 2018-2020 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

// the configuration for the http backend service
type Config struct {
	Id             string
	LogLevel       string
	AuthMode       string
	Path           string
	Port           string
	Username       string
	Password       string
	Metrics        bool
	DbName         string
	DbConnString   string
	DbUsername     string
	DbPassword     string
	SchemaURI      string
	SchemaUsername string
	SchemaToken    string
}

func NewConfig() (*Config, error) {
	log.Info().Msg("loading configuration")

	// use viper to load configuration data
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	// reads the configuration file
	err := v.ReadInConfig()
	if err != nil { // handle errors reading the config file
		log.Error().Msgf("fatal error config file: %s \n", err)
		return nil, err
	}

	// binds all environment variables to make it container friendly
	v.AutomaticEnv()
	v.SetEnvPrefix("OX_DBM") // prefixes all env vars

	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv("Id")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Port")
	_ = v.BindEnv("AuthMode")
	_ = v.BindEnv("Username")
	_ = v.BindEnv("Password")
	_ = v.BindEnv("Metrics")
	_ = v.BindEnv("Db.Name")
	_ = v.BindEnv("Db.ConnString")
	_ = v.BindEnv("Db.Username")
	_ = v.BindEnv("Db.Password")
	_ = v.BindEnv("Schema.URI")
	_ = v.BindEnv("Schema.Username")
	_ = v.BindEnv("Schema.Token")

	// creates a config struct and populate it with values
	c := new(Config)

	// general configuration
	c.Id = v.GetString("Id")
	c.LogLevel = v.GetString("LogLevel")
	c.Port = v.GetString("Port")
	c.AuthMode = v.GetString("AuthMode")
	c.Username = v.GetString("Username")
	c.Password = v.GetString("Password")
	c.Metrics = v.GetBool("Metrics")
	c.DbName = v.GetString("Db.Name")
	c.DbConnString = v.GetString("Db.ConnString")
	c.DbUsername = v.GetString("Db.Username")
	c.DbPassword = v.GetString("Db.Password")
	c.SchemaURI = v.GetString("Schema.URI")
	c.SchemaUsername = v.GetString("Schema.Username")
	c.SchemaToken = v.GetString("Schema.Token")

	return c, nil
}
