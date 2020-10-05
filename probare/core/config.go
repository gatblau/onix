/*
*    Onix Probare - Demo Application for reactive config management
*    Copyright (c) 2020 by www.gatblau.org
*
*    Licensed under the Apache License, Version 2.0 (the "License");
*    you may not use this file except in compliance with the License.
*    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*    Unless required by applicable law or agreed to in writing, software distributed under
*    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
*    either express or implied.
*    See the License for the specific language governing permissions and limitations under the License.
*
*    Contributors to this project, hereby assign copyright in this code to the project,
*    to be licensed under the same terms as the rest of the code.
 */
package core

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

var (
	AppBinds     = []string{"Log.Level", "Banner.Type", "Banner.Message"}
	SecretsBinds = []string{"User", "Pwd"}
)

type config struct {
	// the configuration string content
	content string
	// the viper instance
	store *viper.Viper
	// the name of the configuration file
	filename string
	// the viper environment variable bindings
	binds []string
}

func NewConfig(filename string, binds []string) (*config, error) {
	// use viper to load configuration data from the specified file
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	c := &config{
		store:    v,
		filename: filename,
		binds:    binds,
	}

	// load the configuration
	err := c.LoadFromFile()

	return c, err
}

func (c *config) GetString(key string) string {
	return c.store.GetString(key)
}

// get the configuration file content
// if no reader is passed then it reads from the file system otherwise uses the reader
func (c *config) GetFileContent() ([]byte, error) {
	return ioutil.ReadFile(c.store.ConfigFileUsed())
}

func (c *config) LoadFromFile() error {
	sendMsg(Terminal, []string{fmt.Sprintf("reloading '%s' configuration from file", c.filename)})
	return c.Load("")
}

func (c *config) Load(cfg string) error {
	if len(cfg) > 0 {
		c.content = cfg
	} else {
		fileContent, err := c.GetFileContent()
		if err != nil {
			sendMsg(Terminal, []string{fmt.Sprintf("failed to load '%s' configuration from file: %v", c.filename, err)})
		}
		c.content = string(fileContent)
	}
	if len(cfg) == 0 {
		// reads the configuration file
		err := c.store.ReadInConfig()
		if err != nil { // handle errors reading the config file
			log.Error().Msgf("cannot read config from file %s: %s.toml \n", c.filename, err)
			return err
		}
	} else {
		// reads from the reader
		err := c.store.ReadConfig(strings.NewReader(cfg))
		if err != nil { // handle errors reading the config file
			log.Error().Msgf("cannot read config from reader %s: %s.toml \n", c.filename, err)
			return err
		}
	}

	// binds all environment variables to make it container friendly
	c.store.AutomaticEnv()
	c.store.SetEnvPrefix("PROBE_") // prefixes all env vars

	// replace character to support environment variable format
	c.store.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// binds the specified Viper keys to ENV variables
	for _, bind := range c.binds {
		_ = c.store.BindEnv(bind)
	}

	// set the log level
	logLevel := c.GetString("Log.Level")
	if len(logLevel) > 0 {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			log.Warn().Msg(err.Error())
			log.Info().Msg("defaulting log level to INFO")
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		zerolog.SetGlobalLevel(level)
	}
	return nil
}
