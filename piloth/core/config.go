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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
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
	case PilotCtlUri:
		return "PILOTCTL_URI"
	case PilotArtRegUser:
		return "PILOT_ART_REG_USER"
	case PilotArtRegPwd:
		return "PILOT_ART_REG_PWD"
	case PilotSyslogPort:
		return "PILOT_SYSLOG_PORT"
	}
	return ""
}

const (
	PilotLogLevel ConfigKey = iota
	PilotCtlUri
	PilotArtRegUser
	PilotArtRegPwd
	PilotSyslogPort
)

func (c *Config) getPilotCtlURI() string {
	uri := c.Get(PilotCtlUri)
	if len(uri) == 0 {
		ErrorLogger.Printf("PILOTCTL_URI not defined\n")
		os.Exit(1)
	}
	uri = strings.ToLower(uri)
	if !strings.HasPrefix(uri, "http") {
		ErrorLogger.Printf("PILOTCTL_URI does not define a protocol (preferably use https:// - http links are not secure)\n")
		os.Exit(1)
	}
	return uri
}

func (c *Config) getSyslogPort() string {
	port := c.Get(PilotSyslogPort)
	if len(port) == 0 {
		// set default
		port = "1514"
	}
	return port
}

func (c *Config) Get(key ConfigKey) string {
	return os.Getenv(key.String())
}

func (c *Config) GetBool(key ConfigKey) bool {
	b, _ := strconv.ParseBool(c.Get(key))
	return b
}

func (c *Config) Load() error {
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

func ConfFile() string {
	return fmt.Sprintf("%s/.pilot", CurrentPath())
}

// DataPath returns the path of the root local folder where files are cached
func DataPath() string {
	return filepath.Join(CurrentPath(), "data")
}

// SubmitPath returns the path of the local folder used to cache information to be submitted to pilotctl
func SubmitPath() string {
	return filepath.Join(DataPath(), "submit")
}

// ProcessPath returns the path of the local folder used to cache jobs to be processed by pilot
func ProcessPath() string {
	return filepath.Join(DataPath(), "process")
}

// CheckPaths check all required local folders used by the pilot to cache data exist and if not creates them
func CheckPaths() {
	_, err := os.Stat(DataPath())
	if err != nil {
		err = os.MkdirAll(DataPath(), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	_, err = os.Stat(path.Join(DataPath(), "submit"))
	if err != nil {
		err = os.MkdirAll(path.Join(DataPath(), "submit"), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	_, err = os.Stat(path.Join(DataPath(), "process"))
	if err != nil {
		err = os.MkdirAll(path.Join(DataPath(), "process"), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
