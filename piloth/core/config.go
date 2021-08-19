package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	// path to the current executable
	path string
}

type ConfigKey int

func (k ConfigKey) String() string {
	switch k {
	case PilotLogLevel:
		return "PILOT_LOG_LEVEL"
	case PilotCtlUri:
		return "PILOTCTL_URI"
	case PilotArtRegUser:
		return "PILOT_ART_REG_USER"
	case PilotArtRegPwd:
		return "PILOT_ART_REG_PWD"
	}
	return ""
}

const (
	PilotLogLevel ConfigKey = iota
	PilotCtlUri
	PilotArtRegUser
	PilotArtRegPwd
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

	cFile := confFile()
	if _, err := os.Stat(cFile); err == nil {
		err := godotenv.Load(cFile)
		if err != nil {
			log.Fatal().Msg(err.Error())
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

	return nil
}

func currentPath() string {
	path := os.Getenv("PILOT_CFG_PATH")
	if len(path) > 0 {
		path, err := core.AbsPath(path)
		if err != nil {
			panic(err)
		}
		return path
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func confFile() string {
	return fmt.Sprintf("%s/.pilot", currentPath())
}

func CachePath() string {
	return filepath.Join(currentPath(), "cache")
}

func CheckCachePath() {
	_, err := os.Stat(CachePath())
	if err != nil {
		err = os.MkdirAll(CachePath(), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
