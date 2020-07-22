//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	AppVersion       = "AppVersion"
	ThemeName        = "Theme"
	Plugins          = "Plugins"
	HttpMetrics      = "Http.Metrics"
	HttpAuthMode     = "Http.AuthMode"
	HttpUsername     = "Http.Username"
	HttpPassword     = "Http.Password"
	HttpPort         = "Http.Port"
	RepoURI          = "Repo.URI"
	RepoUsername     = "Repo.Username"
	RepoPassword     = "Repo.Password"
	DbProvider       = "Db.Provider"
	DbHost           = "Db.Host"
	DbPort           = "Db.Port"
	DbName           = "Db.Name"
	DbUsername       = "Db.Username"
	DbPassword       = "Db.Password"
	DbAdminUser      = "Db.AdminUsername"
	DbAdminPwd       = "Db.AdminPassword"
	DbObjectsPattern = "Db.ObjectsPattern"
)

// dbman configuration management struct
type Config struct {
	Cache *Cache
	cfg   *viper.Viper
}

// create a new instance
func NewConfig(path string, name string) *Config {
	conf := &Config{
		Cache: NewCache(),
		cfg:   viper.New(),
	}
	conf.Load(path, name)
	return conf
}

// load a dbman configuration
// path: the configuration file path - if empty is passed-in then home directory is used
// name: the configuration name used to create a filename as follows: .dbman_[name].toml
func (c *Config) Load(path string, name string) error {
	// if no name is specified then use the cached name
	if len(name) == 0 {
		// get it from the root configuration
		name = c.Cache.filename()
	} else if name != c.Cache.name() {
		// if the name is different from the one cached then update the cache
		c.Cache.setName(name)
		c.Cache.save()
	}

	// if no path is used, then used the cached path
	if len(path) == 0 {
		path = c.Cache.Path()
	} else if path != c.Cache.Path() {
		// if the path is different from the one cached then update the cache
		c.Cache.setPath(path)
		c.Cache.save()
	}

	// ensures the config file name is prepended with a dot to make it hidden
	c.cfg.SetConfigName(c.Cache.filename())
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
		err = c.createCfgFile(path, c.Cache.filename())
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

	_ = c.cfg.BindEnv("AppVersion")
	_ = c.cfg.BindEnv("Theme")
	_ = c.cfg.BindEnv("Http.Port")
	_ = c.cfg.BindEnv("Http.AuthMode")
	_ = c.cfg.BindEnv("Http.Username")
	_ = c.cfg.BindEnv("Http.Password")
	_ = c.cfg.BindEnv("Http.Metrics")
	_ = c.cfg.BindEnv("Db.Name")
	_ = c.cfg.BindEnv("Db.Host")
	_ = c.cfg.BindEnv("Db.Port")
	_ = c.cfg.BindEnv("Db.Provider")
	_ = c.cfg.BindEnv("Db.Username")
	_ = c.cfg.BindEnv("Db.Password")
	_ = c.cfg.BindEnv("Db.AdminUsername")
	_ = c.cfg.BindEnv("Db.AdminPassword")
	_ = c.cfg.BindEnv("Repo.URI")
	_ = c.cfg.BindEnv("Repo.Username")
	_ = c.cfg.BindEnv("Repo.Password")

	return nil
}

// return the configuration file used
func (c *Config) ConfigFileUsed() string {
	return c.cfg.ConfigFileUsed()
}

// creates a default configuration file
func (c *Config) createCfgFile(filePath string, filename string) error {
	fmt.Println("? I am writing configuration file to disk")
	f, err := os.Create(fmt.Sprintf("%v/%v.toml", filePath, filename))
	if err != nil {
		fmt.Printf("!!! I failed to create a new configuration file: %s\n", err)
		return err
	}
	l, err := f.WriteString(cfgFile)
	if err != nil {
		fmt.Printf("!!! I failed to create a new configuration file: %s\n", err)
		f.Close()
		return err
	}
	fmt.Printf("? I have written %v bytes to %v/%v.toml\n", l, filePath, filename)
	err = f.Close()
	if err != nil {
		fmt.Printf("!!! I cannot close the configuration file: %s\n", err)
		return err
	}
	return err
}

// save the configuration to file
func (c *Config) Save() {
	err := c.cfg.WriteConfig()
	if err != nil {
		fmt.Printf("!!! I could not save configuration: %v", err)
	}
	err = c.cfg.ReadInConfig()
	if err != nil {
		fmt.Printf("!!! I could not read configuration: %v", err)
	}
}

// check if a key is contained in the internal viper registry
func (c *Config) contains(key string) bool {
	keys := c.cfg.AllKeys()
	for _, a := range keys {
		if a == key {
			return true
		}
	}
	return false
}

// Get returns the value for the specified key
func (c *Config) Get(ctx context.Context, key string) interface{} {
	return c.cfg.GetString(key)
}

// Set sets the value for the specified key
func (c *Config) Set(key string, value interface{}) {
	key = strings.ToLower(key)
	// if key passed in is not standard (i.e. not part of the default set of config keys)
	if !c.contains(key) {
		// warn the user in case they misspelled a standard key
		fmt.Printf("! The key '%v' you provided is not standard, I am adding it to the configuration set.\n", key)
	}
	// updates the key
	c.cfg.Set(key, value)
}

func (c *Config) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

// toString the current configuration file
func (c *Config) ToString() string {
	var (
		buffer bytes.Buffer
		line   string
	)
	for _, key := range c.cfg.AllKeys() {
		if !strings.Contains(strings.ToLower(key), "password") {
			line = fmt.Sprintf("%s = %v\n", key, c.cfg.Get(key))
		} else {
			line = fmt.Sprintf("%s = ???????\n", key)
		}
		buffer.WriteString(line)
	}
	return buffer.String()
}

func (c *Config) All() string {
	m := c.cfg.AllSettings()
	bytes, e := json.Marshal(m)
	if e != nil {
		return ""
	}
	return string(bytes)
}

func (c *Config) GetString(key string) string {
	return c.cfg.GetString(key)
}

// default config file content
const cfgFile = `AppVersion = "0.0.4"
Theme = ""
[Http]
	Metrics  = "true"
	AuthMode = "basic"
	Port     = "8085"
	Username = "admin"
	Password = "0n1x"
[Db]
    Provider      = "_pgsql"
    Name          = "onix"
    Host          = "localhost"
    Port          = "5432"
    Username      = "onix"
    Password      = "onix"
    AdminUsername = "postgres"
    AdminPassword = "onix"
[Repo]
    URI      = "https://raw.githubusercontent.com/gatblau/ox-db/master"
    Username = ""
    Password = ""
`
