//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"strings"
)

// dbman root configuration management struct
type RootCfg struct {
	cfg *viper.Viper
}

// create a new RootCfg instance
func NewRootCfg() *RootCfg {
	rootCfg := &RootCfg{
		cfg: viper.New(),
	}
	rootCfg.load()
	return rootCfg
}

// load the RootCfg
func (c *RootCfg) load() {
	home := homeDir()
	c.cfg.AddConfigPath(home)
	c.cfg.SetConfigName(".dbman")
	c.cfg.SetConfigType("toml")
	err := c.cfg.ReadInConfig()
	if err != nil {
		c.setPath(home)
		c.setName("default")
		c.save()
		err := c.cfg.ReadInConfig()
		if err != nil {
			fmt.Printf("cannot save root configuration: %v", err)
		}
	}
}

// save the config path & name to file
func (c *RootCfg) save() {
	err := c.cfg.WriteConfig()
	if err != nil {
		fmt.Printf("oops! cannot save cache: %v", err)
	}
}

// return the config file path to use
func (c *RootCfg) path() string {
	return c.cfg.GetString("path")
}

// return the config name to use
func (c *RootCfg) name() string {
	return c.cfg.GetString("name")
}

// return the config file name to use
func (c *RootCfg) filename() string {
	return fmt.Sprintf(".dbman_%v", c.name())
}

func (c *RootCfg) setName(name string) {
	// check file name should not contain a file extension
	if strings.Contains(name, ".toml") ||
		strings.Contains(name, ".json") ||
		strings.Contains(name, ".yaml") ||
		strings.Contains(name, ".yml") ||
		strings.Contains(name, ".yaml") ||
		strings.Contains(name, ".txt") {
		fmt.Printf("file extension not allowed in config filename %v", name)
		return
	}

	// invalid if the file name is prepended by '.'
	if strings.Index(name, ".") == 1 {
		fmt.Printf("invalid name '%v' should not start with '.'", name)
		return
	}
	c.cfg.Set("name", name)
}

func (c *RootCfg) setPath(value string) {
	c.cfg.Set("path", value)
}

// get the home directory
func homeDir() string {
	// find home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("cant find home directory: %v\n", err)
		return ""
	}
	return home
}
