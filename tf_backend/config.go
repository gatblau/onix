/*
   Terraform Http Backend - Onix - Copyright (c) 2018 by www.gatblau.org

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
	. "gatblau.org/onix/wapic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

// the configuration for the http backend service
type Config struct {
	LogLevel string
	Id       string
	Onix     *Onix    // configuration for Onix integration
	Service  *SvcConf // configuration for the http backend service
}

// the configuration for the http backend endpoint
type SvcConf struct {
	Path               string
	Port               string
	AuthMode           string
	Username           string
	Password           string
	Metrics            bool
	InsecureSkipVerify bool
}

func NewConfig() (Config, error) {
	log.Infof("Loading configuration.")
	v := viper.New()
	// loads the configuration file
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	err := v.ReadInConfig() // find and read the config file
	if err != nil {         // handle errors reading the config file
		log.Errorf("Fatal error config file: %s \n", err)
		return Config{}, err
	}

	// binds all environment variables to make it container friendly
	v.AutomaticEnv()
	v.SetEnvPrefix("OXTFB") // short for Onix Terraform Backend

	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv("Id")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Onix.URL")
	_ = v.BindEnv("Onix.AuthMode")
	_ = v.BindEnv("Onix.Username")
	_ = v.BindEnv("Onix.Password")
	_ = v.BindEnv("Onix.ClientId")
	_ = v.BindEnv("Onix.ClientSecret")
	_ = v.BindEnv("Onix.TokenURI")
	_ = v.BindEnv("Service.Path")
	_ = v.BindEnv("Service.Port")
	_ = v.BindEnv("Service.AuthMode")
	_ = v.BindEnv("Service.Username")
	_ = v.BindEnv("Service.Password")
	_ = v.BindEnv("Service.Metrics")
	_ = v.BindEnv("Service.InsecureSkipVerify")

	// creates a config struct and populate it with values
	c := new(Config)
	c.Onix = new(Onix)
	c.Service = new(SvcConf)

	// general configuration
	c.Id = v.GetString("Id")
	c.LogLevel = v.GetString("LogLevel")
	c.Onix.URL = v.GetString("Onix.URL")
	c.Onix.AuthMode = v.GetString("Onix.AuthMode")
	c.Onix.Username = v.GetString("Onix.Username")
	c.Onix.Password = v.GetString("Onix.Password")
	c.Onix.ClientId = v.GetString("Onix.ClientId")
	c.Onix.ClientSecret = v.GetString("Onix.ClientSecret")
	c.Onix.TokeURI = v.GetString("Onix.TokenURI")
	c.Service.AuthMode = v.GetString("Service.AuthMode")
	c.Service.InsecureSkipVerify = v.GetBool("Service.InsecureSkipVerify")
	c.Service.Metrics = v.GetBool("Service.Metrics")
	c.Service.Username = v.GetString("Service.Username")
	c.Service.Password = v.GetString("Service.Password")
	c.Service.Port = v.GetString("Service.Port")
	c.Service.Path = v.GetString("Service.Path")

	return *c, nil
}
