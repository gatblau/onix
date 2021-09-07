package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ConfKey string

const (
	ConfDbName                   ConfKey = "OX_PILOTCTL_DB_NAME"
	ConfDbHost                   ConfKey = "OX_PILOTCTL_DB_HOST"
	ConfDbPort                   ConfKey = "OX_PILOTCTL_DB_PORT"
	ConfDbUser                   ConfKey = "OX_PILOTCTL_DB_USER"
	ConfDbPwd                    ConfKey = "OX_PILOTCTL_DB_PWD"
	ConfPingIntervalSecs         ConfKey = "OX_PILOCTL_PING_INTERVAL_SECS"
	ConfOxWapiUri                ConfKey = "OX_WAPI_URI"
	ConfOxWapiUser               ConfKey = "OX_WAPI_USER"
	ConfOxWapiPwd                ConfKey = "OX_WAPI_PWD"
	ConfOxWapiInsecureSkipVerify ConfKey = "OX_WAPI_INSECURE_SKIP_VERIFY"
	ConfArtRegURI                ConfKey = "OX_ART_REG_URI"
	ConfArtRegUser               ConfKey = "OX_ART_REG_USER"
	ConfArtRegPwd                ConfKey = "OX_ART_REG_PWD"
)

type Conf struct {
}

func NewConf() *Conf {
	return &Conf{}
}

func (c *Conf) get(key ConfKey) string {
	return os.Getenv(string(key))
}

func (c *Conf) getDbName() string {
	value := os.Getenv(string(ConfDbName))
	if len(value) == 0 {
		return "pilotctl"
	}
	return value
}

func (c *Conf) getDbHost() string {
	return c.getValue(ConfDbHost)
}

func (c *Conf) getDbPort() string {
	value := os.Getenv(string(ConfDbPort))
	if len(value) == 0 {
		return "5432"
	}
	return value
}

func (c *Conf) getDbUser() string {
	return c.getValue(ConfDbUser)
}

func (c *Conf) getDbPwd() string {
	return c.getValue(ConfDbPwd)
}

// PingIntervalSecs the pilot ping interval
func (c *Conf) PingIntervalSecs() time.Duration {
	defaultValue, _ := time.ParseDuration("15s")
	value := os.Getenv(string(ConfPingIntervalSecs))
	if len(value) == 0 {
		return defaultValue
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("WARNING: %s is invalid, defaulting to %d\n", ConfPingIntervalSecs, defaultValue)
		return defaultValue
	}
	interval, err := time.ParseDuration(fmt.Sprintf("%ds", v))
	if err != nil {
		fmt.Printf("WARNING: %s is invalid, defaulting to %d\n", ConfPingIntervalSecs, defaultValue)
		return defaultValue
	}
	return interval
}

// DisconnectedAfterSecs the period after which pilotctl considers a host disconnected if the last seen time is not updated
func (c *Conf) DisconnectedAfterSecs() int {
	in := int(c.PingIntervalSecs())
	after := in + 5
	return after
}

func (c *Conf) getOxWapiUrl() string {
	return c.getValue(ConfOxWapiUri)
}

func (c *Conf) getOxWapiUsername() string {
	return c.getValue(ConfOxWapiUser)
}

func (c *Conf) getOxWapiPassword() string {
	return c.getValue(ConfOxWapiPwd)
}

func (c *Conf) getArtRegUri() string {
	return c.getValue(ConfArtRegURI)
}

func (c *Conf) getArtRegUser() string {
	return c.getValue(ConfArtRegUser)
}

func (c *Conf) getArtRegPwd() string {
	return c.getValue(ConfArtRegPwd)
}

func (c *Conf) getValue(key ConfKey) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		fmt.Printf("ERROR: variable %s not defined", key)
		os.Exit(1)
	}
	return value
}

func (c *Conf) getOxWapiInsecureSkipVerify() bool {
	b, err := strconv.ParseBool(c.getValue(ConfOxWapiInsecureSkipVerify))
	if err != nil {
		fmt.Printf("ERROR: invalid value for variable %s", ConfOxWapiInsecureSkipVerify)
		os.Exit(1)
	}
	return b
}
