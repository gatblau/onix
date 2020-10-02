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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"strings"
)

var (
	AppBinds     = []string{"Log.Level", "Banner.Type", "Banner.Message"}
	SecretsBinds = []string{"User", "Pwd"}
)

type config struct {
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
	err := c.Load(nil)

	return c, err
}

func (c *config) GetString(key string) string {
	return c.store.GetString(key)
}

func (c *config) GetContent() ([]byte, error) {
	return ioutil.ReadFile(c.store.ConfigFileUsed())
}

func (c *config) Load(r io.Reader) error {
	if r == nil {
		// reads the configuration file
		err := c.store.ReadInConfig()
		if err != nil { // handle errors reading the config file
			log.Error().Msgf("cannot read config from file %s: %s.toml \n", c.filename, err)
			return err
		}
	} else {
		// reads from the reader
		err := c.store.ReadConfig(r)
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
