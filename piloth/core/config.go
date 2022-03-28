/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config pilot configuration
type Config struct {
	LogLevel string
	// path to the current executable
	path string
}

type ConfigKey int

func (k ConfigKey) String() string {
	switch k {
	case PilotLogLevel:
		return "PILOT_LOG_LEVEL"
	case PilotSyslogPort:
		return "PILOT_SYSLOG_PORT"
	case PilotActivationURI:
		return "PILOT_ACTIVATION_URI"
	case PilotUserKey:
		return "PILOT_USER_KEY"
	}
	return ""
}

const (
	PilotLogLevel ConfigKey = iota
	PilotSyslogPort
	PilotActivationURI
	PilotUserKey
)

func (c *Config) getSyslogPort() string {
	defer TRA(CE())
	port := c.Get(PilotSyslogPort)
	if len(port) == 0 {
		// set default
		port = "1514"
	}
	return port
}

func (c *Config) Get(key ConfigKey) string {
	defer TRA(CE())
	return os.Getenv(key.String())
}

func (c *Config) GetBool(key ConfigKey) bool {
	defer TRA(CE())
	b, _ := strconv.ParseBool(c.Get(key))
	return b
}

func (c *Config) Load() error {
	defer TRA(CE())
	// set the file path to where pilot is running
	// c.path = currentPath()
	c.path = CurrentPath()

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

	return nil
}

func CurrentPath() string {
	// check if the current path is overridden
	path := os.Getenv("PILOT_CFG_PATH")
	// if so
	if len(path) > 0 {
		// works out the absolute path and return
		path, err := core.AbsPath(path)
		if err != nil {
			ErrorLogger.Printf(err.Error())
			os.Exit(1)
		}
		return path
	}
	// otherwise, get the path where pilot is located
	path, err := core.AbsPath(".")
	if err != nil {
		ErrorLogger.Printf(err.Error())
		os.Exit(1)
	}
	return path
}

func AkFile() string {
	defer TRA(CE())
	return fmt.Sprintf("%s/.pilot", CurrentPath())
}

func UserKeyFile() string {
	defer TRA(CE())
	return fmt.Sprintf("%s/.userkey", CurrentPath())
}

// DataPath returns the path of the root local folder where files are cached
func DataPath() string {
	return filepath.Join(CurrentPath(), "data")
}

// SubmitPath returns the path of the local folder used to cache information to be submitted to pilotctl
func SubmitPath() string {
	defer TRA(CE())
	return filepath.Join(DataPath(), "submit")
}

// ProcessPath returns the path of the local folder used to cache jobs to be processed by pilot
func ProcessPath() string {
	return filepath.Join(DataPath(), "process")
}
