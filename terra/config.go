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
	LogLevel string
	Id       string
	Ox       *OnixConf  // configuration for the Onix client
	Terra    *TerraConf // configuration for Terra, the http backend service for Onix
}

// configuration for Terra
type TerraConf struct {
	Path               string
	Port               string
	AuthMode           string
	Username           string
	Password           string
	Metrics            bool
	InsecureSkipVerify bool
}

// config for Onix client
type OnixConf struct {
	URL          string
	AuthMode     string
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	TokenURI     string
}

func NewConfig() (*Config, error) {
	log.Info().Msg("Loading configuration.")

	// use viper to load configuration data
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	// reads the configuration file
	err := v.ReadInConfig()
	if err != nil { // handle errors reading the config file
		log.Error().Msgf("Fatal error config file: %s \n", err)
		return nil, err
	}

	// binds all environment variables to make it container friendly
	v.AutomaticEnv()
	v.SetEnvPrefix("OX_TERRA") // prefixes all env vars

	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv("Id")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Onix.URL")
	_ = v.BindEnv("Onix.AuthMode")
	_ = v.BindEnv("Onix.Username")
	_ = v.BindEnv("Onix.Password")
	_ = v.BindEnv("Onix.ClientId")
	_ = v.BindEnv("Onix.AppSecret")
	_ = v.BindEnv("Onix.TokenURI")
	_ = v.BindEnv("TerraService.Path")
	_ = v.BindEnv("TerraService.Port")
	_ = v.BindEnv("TerraService.AuthMode")
	_ = v.BindEnv("TerraService.Username")
	_ = v.BindEnv("TerraService.Password")
	_ = v.BindEnv("TerraService.Metrics")
	_ = v.BindEnv("TerraService.InsecureSkipVerify")

	// creates a config struct and populate it with values
	c := new(Config)
	c.Ox = new(OnixConf)
	c.Terra = new(TerraConf)

	// general configuration
	c.Id = v.GetString("Id")
	c.LogLevel = v.GetString("LogLevel")
	c.Ox.URL = v.GetString("Onix.URL")
	c.Ox.AuthMode = v.GetString("Onix.AuthMode")
	c.Ox.Username = v.GetString("Onix.Username")
	c.Ox.Password = v.GetString("Onix.Password")
	c.Ox.ClientId = v.GetString("Onix.ClientId")
	c.Ox.ClientSecret = v.GetString("Onix.AppSecret")
	c.Ox.TokenURI = v.GetString("Onix.TokenURI")
	c.Terra.AuthMode = v.GetString("Service.AuthMode")
	c.Terra.InsecureSkipVerify = v.GetBool("Service.InsecureSkipVerify")
	c.Terra.Metrics = v.GetBool("Service.Metrics")
	c.Terra.Username = v.GetString("Service.Username")
	c.Terra.Password = v.GetString("Service.Password")
	c.Terra.Port = v.GetString("Service.Port")
	c.Terra.Path = v.GetString("Service.Path")

	return c, nil
}
