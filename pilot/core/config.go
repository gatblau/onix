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
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LogLevel = "LogLevel"
)

// create a logger with added contextual info
// to differentiate pilot logging from application logging
var logger = log.With().Str("agent", "pilot").Logger()

// Config pilot configuration
type Config struct {
	LogLevel string
	// configuration for the Onix client
	OxConf *oxc.ClientConf
	// message broker config
	EmConf *oxc.EventConfig
	// the viper instance
	store *viper.Viper
	// path to the current executable
	path string
}

type ConfigKey int

func (k ConfigKey) String() string {
	switch k {
	case PilotKey:
		return "PILOT_APP_KEY"
	case PilotLogLevel:
		return "PILOT_LOG_LEVEL"
	case PilotOxWapiUrl:
		return "PILOT_OX_WAPI_URL"
	case PilotOxWapiAuthMode:
		return "PILOT_OX_WAPI_AUTHMODE"
	case PilotOxWapiUsername:
		return "PILOT_OX_WAPI_USERNAME"
	case PilotOxWapiPassword:
		return "PILOT_OX_WAPI_PASSWORD"
	case PilotOxWapiClientId:
		return "PILOT_OX_WAPI_CLIENT_ID"
	case PilotOxWapiAppSecret:
		return "PILOT_OX_WAPI_APP_SECRET"
	case PilotOxWapiTokenUri:
		return "PILOT_OX_WAPI_TOKEN_URI"
	case PilotOxWapiInsecureSkipVerify:
		return "PILOT_OX_WAPI_INSECURESKIPVERIFY"
	case PilotOxBrokerUrl:
		return "PILOT_OX_BROKER_URL"
	case PilotOxBrokerUsername:
		return "PILOT_OX_BROKER_USERNAME"
	case PilotOxBrokerPassword:
		return "PILOT_OX_BROKER_PASSWORD"
	case PilotOxBrokerInsecureSkipVerify:
		return "PILOT_OX_BROKER_INSECURESKIPVERIFY"
	case PilotRemUri:
		return "PILOT_OX_REM_URI"
	case PilotRemUsername:
		return "PILOT_OX_REM_USERNAME"
	case PilotRemPassword:
		return "PILOT_OX_REM_PASSWORD"
	}
	return ""
}

const (
	PilotKey ConfigKey = iota
	PilotLogLevel
	PilotOxWapiUrl
	PilotOxWapiAuthMode
	PilotOxWapiUsername
	PilotOxWapiPassword
	PilotOxWapiClientId
	PilotOxWapiAppSecret
	PilotOxWapiTokenUri
	PilotOxWapiInsecureSkipVerify
	PilotOxBrokerUrl
	PilotOxBrokerUsername
	PilotOxBrokerPassword
	PilotOxBrokerInsecureSkipVerify
	PilotRemUri
	PilotRemUsername
	PilotRemPassword
)

func (c *Config) Get(key ConfigKey) string {
	return os.Getenv(key.String())
}

func (c *Config) GetBool(key ConfigKey) bool {
	b, _ := strconv.ParseBool(c.Get(key))
	return b
}

func (c *Config) Load() error {
	logger.Info().Msg("Loading configuration.")

	// set the file path to where pilot is running
	c.path = currentPath()

	// load configuration from .env file if exist
	if _, err := os.Stat(fmt.Sprintf("%s/.env", c.path)); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	} else {
		envPath := os.Getenv("PILOT_ENV_PATH")
		if len(envPath) > 0 {
			path := fmt.Sprintf("%s/.env", envPath)
			if _, err := os.Stat(path); err == nil {
				err := godotenv.Load(path)
				if err != nil {
					log.Fatal().Msg(err.Error())
				}
			}
		}
	}

	// log level
	c.LogLevel = c.Get(PilotLogLevel)

	logLevel, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		log.Warn().Msg(err.Error())
		log.Info().Msg("defaulting log level to INFO")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Ox client config
	oxcfg := &oxc.ClientConf{
		BaseURI:            c.Get(PilotOxWapiUrl),
		Username:           c.Get(PilotOxWapiUsername),
		Password:           c.Get(PilotOxWapiPassword),
		ClientId:           c.Get(PilotOxWapiClientId),
		AppSecret:          c.Get(PilotOxWapiAppSecret),
		TokenURI:           c.Get(PilotOxWapiTokenUri),
		InsecureSkipVerify: c.GetBool(PilotOxWapiInsecureSkipVerify),
	}
	oxcfg.SetAuthMode(c.Get(PilotOxWapiAuthMode))

	// event manager config
	emcfg := &oxc.EventConfig{
		Server:             c.Get(PilotOxBrokerUrl),
		ItemInstance:       c.Get(PilotKey),
		Qos:                2,
		Username:           c.Get(PilotOxBrokerUsername),
		Password:           c.Get(PilotOxBrokerPassword),
		InsecureSkipVerify: c.GetBool(PilotOxBrokerInsecureSkipVerify),
		ClientAuthType:     tls.NoClientCert,
		OnMsgReceived:      nil,
	}
	c.OxConf = oxcfg
	c.EmConf = emcfg

	return nil
}

func currentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
