package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
	"strconv"
)

type ConfKey string

const (
	ConfDbName          ConfKey = "OX_REM_DB_NAME"
	ConfDbHost          ConfKey = "OX_REM_DB_HOST"
	ConfDbPort          ConfKey = "OX_REM_DB_PORT"
	ConfDbUser          ConfKey = "OX_REM_DB_USER"
	ConfDbPwd           ConfKey = "OX_REM_DB_PWD"
	ConfRefreshInterval ConfKey = "OX_REM_REFRESH_INTERVAL"
	ConfPingInterval    ConfKey = "OX_REM_PING_INTERVAL"
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
		return "rem"
	}
	return value
}

func (c *Conf) getDbHost() string {
	value := os.Getenv(string(ConfDbHost))
	if len(value) == 0 {
		return "localhost"
	}
	return value
}

func (c *Conf) getDbPort() string {
	value := os.Getenv(string(ConfDbPort))
	if len(value) == 0 {
		return "5432"
	}
	return value
}

func (c *Conf) getDbUser() string {
	value := os.Getenv(string(ConfDbUser))
	if len(value) == 0 {
		return "rem"
	}
	return value
}

func (c *Conf) getDbPwd() string {
	value := os.Getenv(string(ConfDbPwd))
	if len(value) == 0 {
		return "r3m"
	}
	return value
}

func (c *Conf) getRefreshInterval() int {
	defaultValue := 120
	value := os.Getenv(string(ConfRefreshInterval))
	if len(value) == 0 {
		return defaultValue
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("WARNING: %s is not a number, defaulting to %d\n", ConfRefreshInterval, defaultValue)
		return defaultValue
	}
	return v
}

func (c *Conf) GetPingInterval() int {
	defaultValue := 60
	value := os.Getenv(string(ConfPingInterval))
	if len(value) == 0 {
		return defaultValue
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("WARNING: %s is not a number, defaulting to %d\n", ConfPingInterval, defaultValue)
		return defaultValue
	}
	return v
}
