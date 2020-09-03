/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"crypto/tls"
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

const (
	LogLevel = "LogLevel"
)

// pilot configuration
type Config struct {
	LogLevel string
	// configuration for the Onix client
	OxConf *oxc.ClientConf
	// message broker config
	EmConf *oxc.EventConfig
	// app conf file path
	CfgFile string
	// a command to reload the app configuration
	ReloadCmd string
	// a URI of the endpoint to call to reload the configuration
	ReloadURI string
	// a URI to check the application is ready after reloading the configuration
	ReadyURI string
	// the viper instance
	store *viper.Viper
}

func (c *Config) GetString(key string) string {
	return c.store.GetString(key)
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
	v.SetEnvPrefix("OXP") // prefixes all env vars

	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv(LogLevel)

	_ = v.BindEnv("Onix.URL")
	_ = v.BindEnv("Onix.AuthMode")
	_ = v.BindEnv("Onix.Username")
	_ = v.BindEnv("Onix.Password")
	_ = v.BindEnv("Onix.ClientId")
	_ = v.BindEnv("Onix.AppSecret")
	_ = v.BindEnv("Onix.TokenURI")
	_ = v.BindEnv("Onix.InsecureSkipVerify")

	_ = v.BindEnv("Broker.Server")
	_ = v.BindEnv("Broker.Username")
	_ = v.BindEnv("Broker.Password")
	_ = v.BindEnv("Broker.InsecureSkipVerify")

	_ = v.BindEnv("App.Key")
	_ = v.BindEnv("App.CfgFile")
	_ = v.BindEnv("App.ReloadCmd")
	_ = v.BindEnv("App.ReloadURI")
	_ = v.BindEnv("App.ReadyURI")

	// creates a config struct and populate it with values
	c := new(Config)
	c.store = v

	// log level
	c.LogLevel = v.GetString("LogLevel")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		log.Warn().Msg(err.Error())
		log.Info().Msg("defaulting log level to INFO")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(logLevel)

	// Ox client config
	oxcfg := &oxc.ClientConf{
		BaseURI:            v.GetString("Onix.URL"),
		Username:           v.GetString("Onix.Username"),
		Password:           v.GetString("Onix.Password"),
		ClientId:           v.GetString("Onix.ClientId"),
		AppSecret:          v.GetString("Onix.AppSecret"),
		TokenURI:           v.GetString("Onix.TokenURI"),
		InsecureSkipVerify: v.GetBool("Onix.InsecureSkipVerify"),
	}
	oxcfg.SetAuthMode(v.GetString("Onix.AuthMode"))

	// event manager config
	emcfg := &oxc.EventConfig{
		Server:             v.GetString("Broker.Server"),
		ItemInstance:       v.GetString("App.Key"),
		Qos:                2,
		Username:           v.GetString("Broker.Username"),
		Password:           v.GetString("Broker.Password"),
		InsecureSkipVerify: v.GetBool("Broker.InsecureSkipVerify"),
		ClientAuthType:     tls.NoClientCert,
		OnMsgReceived:      nil,
	}
	c.OxConf = oxcfg
	c.EmConf = emcfg
	c.CfgFile = v.GetString("App.CfgFile")
	c.ReloadCmd = v.GetString("App.ReloadCmd")
	c.ReloadURI = v.GetString("App.ReloadURI")
	c.ReadyURI = v.GetString("App.ReadyURI")
	return c, nil
}
