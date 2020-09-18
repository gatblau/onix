/*
   Onix Config Manager - SeS - Onix Webhook Receiver for AlertManager
   Copyright (c) 2020 by www.gatblau.org

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
package server

import (
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

// the configuration for the http backend service
type Config struct {
	LogLevel string
	Id       string
	AuthMode string
	Path     string
	Port     string
	Username string
	Password string
	Metrics  bool
	Ox       *oxc.ClientConf // configuration for the Onix client
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
	v.SetEnvPrefix("OXSES") // prefixes all env vars

	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv("Id")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Path")
	_ = v.BindEnv("Port")
	_ = v.BindEnv("AuthMode")
	_ = v.BindEnv("Username")
	_ = v.BindEnv("Password")
	_ = v.BindEnv("Metrics")
	_ = v.BindEnv("Onix.URL")
	_ = v.BindEnv("Onix.AuthMode")
	_ = v.BindEnv("Onix.Username")
	_ = v.BindEnv("Onix.Password")
	_ = v.BindEnv("Onix.ClientId")
	_ = v.BindEnv("Onix.AppSecret")
	_ = v.BindEnv("Onix.TokenURI")
	_ = v.BindEnv("Onix.InsecureSkipVerify")

	// creates a config struct and populate it with values
	c := new(Config)
	c.Ox = new(oxc.ClientConf)

	// general configuration
	c.Id = v.GetString("Id")
	c.LogLevel = v.GetString("LogLevel")
	c.AuthMode = v.GetString("AuthMode")
	c.Metrics = v.GetBool("Metrics")
	c.Username = v.GetString("Username")
	c.Password = v.GetString("Password")
	c.Port = v.GetString("Port")
	c.Path = v.GetString("Path")

	c.Ox.BaseURI = v.GetString("Onix.URL")
	c.Ox.Username = v.GetString("Onix.Username")
	c.Ox.Password = v.GetString("Onix.Password")
	c.Ox.ClientId = v.GetString("Onix.ClientId")
	c.Ox.AppSecret = v.GetString("Onix.AppSecret")
	c.Ox.TokenURI = v.GetString("Onix.TokenURI")
	c.Ox.SetAuthMode(v.GetString("Onix.AuthMode"))
	c.Ox.InsecureSkipVerify = v.GetBool("Onix.InsecureSkipVerify")

	// set the log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		log.Warn().Msg(err.Error())
		log.Info().Msg("defaulting log level to INFO")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(logLevel)

	return c, nil
}

func (c *Config) debugLevel() bool {
	return strings.ToLower(c.LogLevel) == "debug"
}
