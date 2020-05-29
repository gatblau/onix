//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SchemaURI      = "Schema.URI"
	SchemaUsername = "Schema.Username"
	SchemaToken    = "Schema.Token"
	DbConnString   = "Db.ConnString"
	DbName         = "Db.Name"
	DbUsername     = "Db.Username"
	DbPassword     = "Db.Password"
)

// dbman configuration management struct
type AppCfg struct {
	root *RootCfg
	cfg  *viper.Viper
	path string
	name string
}

// create a new instance
func NewAppCfg(path string, name string) *AppCfg {
	conf := &AppCfg{
		root: NewRootCfg(),
		cfg:  viper.New(),
	}
	conf.load(path, name)
	return conf
}

// load a dbman configuration
// path: the configuration file path - if empty is passed-in then home directory is used
// name: the configuration name used to create a filename as follows: .dbman_[name].toml
func (c *AppCfg) load(path string, name string) error {
	// if no name is specified then use the cached name
	if len(name) == 0 {
		// get it from the root configuration
		name = c.root.filename()
	} else if name != c.root.name() {
		// if the name is different from the one cached then update the cache
		c.root.setName(name)
		c.root.save()
	}

	// if no path is used, then used the cached path
	if len(path) == 0 {
		path = c.root.path()
	} else if path != c.root.path() {
		// if the path is different from the one cached then update the cache
		c.root.setPath(path)
		c.root.save()
	}

	// ensures the config file name is prepended with a dot to make it hidden
	c.cfg.SetConfigName(c.root.filename())
	// always use toml as format
	c.cfg.SetConfigType("toml")

	// if no path is specified then use default path ($HOME)
	if len(path) == 0 {
		// find home directory
		path = homeDir()
	}
	c.cfg.AddConfigPath(path)

	// reads the configuration file
	err := c.cfg.ReadInConfig()
	if err != nil { // handle errors reading the config file
		fmt.Println(err)
		err = c.createCfgFile(path, name)
		if err != nil {
			return err
		}
		c.cfg.ReadInConfig()
	}

	// binds all environment variables to make it container friendly
	c.cfg.AutomaticEnv()
	c.cfg.SetEnvPrefix("OX_DBM") // prefixes all env vars

	// replace character to support environment variable format
	c.cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = c.cfg.BindEnv("LogLevel")
	_ = c.cfg.BindEnv("Http.Port")
	_ = c.cfg.BindEnv("Http.AuthMode")
	_ = c.cfg.BindEnv("Http.Username")
	_ = c.cfg.BindEnv("Http.Password")
	_ = c.cfg.BindEnv("Http.Metrics")
	_ = c.cfg.BindEnv("Db.Name")
	_ = c.cfg.BindEnv("Db.ConnString")
	_ = c.cfg.BindEnv("Db.Username")
	_ = c.cfg.BindEnv("Db.Password")
	_ = c.cfg.BindEnv("Schema.URI")
	_ = c.cfg.BindEnv("Schema.Username")
	_ = c.cfg.BindEnv("Schema.Token")

	return nil
}

// return the configuration file used
func (c *AppCfg) ConfigFileUsed() string {
	return c.cfg.ConfigFileUsed()
}

// creates a default configuration file
func (c *AppCfg) createCfgFile(filePath string, filename string) error {
	fmt.Println("writing configuration file to disk")
	f, err := os.Create(fmt.Sprintf("%v/%v.toml", filePath, filename))
	if err != nil {
		fmt.Printf("failed to create configuration file: %s\n", err)
		return err
	}
	l, err := f.WriteString(cfgFile)
	if err != nil {
		fmt.Printf("failed to write content into configuration file: %s\n", err)
		f.Close()
		return err
	}
	fmt.Printf("%v bytes written successfully to %v/%v.toml\n", l, filePath, filename)
	err = f.Close()
	if err != nil {
		fmt.Printf("failed to close configuration file: %s\n", err)
		return err
	}
	return err
}

// save the configuration to file
func (c *AppCfg) save() {
	c.cfg.WriteConfig()
	c.cfg.ReadInConfig()
}

// check if a key is contained in the internal viper registry
func (c *AppCfg) contains(key string) bool {
	keys := c.cfg.AllKeys()
	for _, a := range keys {
		if a == key {
			return true
		}
	}
	return false
}

// set the configuration value for the passed-in key
// return: true if the value was set or false otherwise
func (c *AppCfg) set(key string, value string) bool {
	key = strings.ToLower(key)
	validKey := c.contains(key)
	// only updates if a valid key is passed in
	if validKey {
		c.cfg.Set(key, value)
	} else {
		fmt.Printf("oops! key '%v' is not recognised, cannot update configuration\n", key)
	}
	return validKey
}

// get a configuration value
func (c *AppCfg) get(key string) string {
	return c.cfg.GetString(key)
}

// print the current configuration file
func (c *AppCfg) print() {
	dat, err := ioutil.ReadFile(c.cfg.ConfigFileUsed())
	if err != nil {
		fmt.Sprintf("cannot read config file: %v", err)
		return
	}
	fmt.Print(string(dat))
}
