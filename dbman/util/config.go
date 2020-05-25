//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// dbman config file name
const CfgFileName = ".dbman"

// global only because it is runtime configuration
var Conf *Config

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

// creates a new configuration object from a file in the specified path
// if no path is specified, then uses the location where dbman is running from
func NewConfig(configPath string) (*Config, error) {
	log.Info().Msg("loading configuration")

	// defines the config file name (always the same)
	viper.SetConfigName(CfgFileName)
	viper.SetConfigType("toml")

	// if no path is specified then use default path ($HOME)
	if len(configPath) == 0 {
		// find home directory
		home, err := homedir.Dir()
		if err != nil {
			log.Error().Msg("cant find home directory whilst trying to set up dbman's configuration file.")
			return nil, err
		}
		configPath = home
	}
	viper.AddConfigPath(configPath)

	// reads the configuration file
	err := viper.ReadInConfig()
	if err != nil { // handle errors reading the config file
		log.Warn().Msgf("cant read from configuration file: %v \n", err)
		err = createDefaultCfgFile(configPath)
		if err != nil {
			return nil, err
		}
	}

	// binds all environment variables to make it container friendly
	viper.AutomaticEnv()
	viper.SetEnvPrefix("OX_DBM") // prefixes all env vars

	// replace character to support environment variable format
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.BindEnv("Id")
	_ = viper.BindEnv("LogLevel")
	_ = viper.BindEnv("Port")
	_ = viper.BindEnv("AuthMode")
	_ = viper.BindEnv("Username")
	_ = viper.BindEnv("Password")
	_ = viper.BindEnv("Metrics")
	_ = viper.BindEnv("Db.Name")
	_ = viper.BindEnv("Db.ConnString")
	_ = viper.BindEnv("Db.Username")
	_ = viper.BindEnv("Db.Password")
	_ = viper.BindEnv("Schema.URI")
	_ = viper.BindEnv("Schema.Username")
	_ = viper.BindEnv("Schema.Token")

	// creates a config struct and populate it with values
	c := new(Config)

	// general configuration
	c.Id = viper.GetString("Id")
	c.LogLevel = viper.GetString("LogLevel")
	c.Port = viper.GetString("Port")
	c.AuthMode = viper.GetString("AuthMode")
	c.Username = viper.GetString("Username")
	c.Password = viper.GetString("Password")
	c.Metrics = viper.GetBool("Metrics")
	c.DbName = viper.GetString("Db.Name")
	c.DbConnString = viper.GetString("Db.ConnString")
	c.DbUsername = viper.GetString("Db.Username")
	c.DbPassword = viper.GetString("Db.Password")
	c.SchemaURI = viper.GetString("Schema.URI")
	c.SchemaUsername = viper.GetString("Schema.Username")
	c.SchemaToken = viper.GetString("Schema.Token")

	return c, nil
}

func (cfg *Config) save() {
	viper.WriteConfig()
}

// creates a default configuration file
func createDefaultCfgFile(filePath string) error {
	log.Info().Msg("creating default configuration")
	f, err := os.Create(fmt.Sprintf("%v/%v.toml", filePath, CfgFileName))
	if err != nil {
		log.Error().Msgf("failed to create configuration file: %s", err)
		return err
	}
	l, err := f.WriteString(cfgFile)
	if err != nil {
		log.Error().Msgf("failed to write content into configuration file: %s", err)
		f.Close()
		return err
	}
	log.Info().Msgf("%v bytes written successfully to %v/%v.toml", l, filePath, CfgFileName)
	err = f.Close()
	if err != nil {
		log.Error().Msgf("failed to close configuration file: %s", err)
		return err
	}
	return err
}

// default config file content
const cfgFile = `
# the dbman service instance Id for logging identification purposes
Id = "DbMan"

# verbosity of logging (Trace, Debug, Warning, Info, Error, Fatal, Panic)
LogLevel = "Trace"

# enables metrics
Metrics = "true"

# configuration for DbMan
AuthMode    = "basic"
Port        = "8085"
Username    = "admin"
Password    = "0n1x"

# configuration for the Onix Web API integration
[Db]
    Name        = "onix"
    ConnString  = "http://localhost:8080"
    Username    = "admin"
    Password    = "0n1x"
[Schema]
    URI         = "https://raw.githubusercontent.com/gatblau/ox-db/master"
    Username    = ""
    Token       = ""
`
