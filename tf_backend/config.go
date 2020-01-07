/*
   Onix Terra - Copyright (c) 2020 by www.gatblau.org

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

type Config struct {
	LogLevel string
	Id       string
	Onix     *Onix
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
	v.SetEnvPrefix("OXTE")

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

	// creates a config struct and populate it with values
	c := new(Config)

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

	return *c, nil
}
