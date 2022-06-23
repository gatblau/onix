/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package mqtt

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gatblau/onix/artisan/core"
)

type ConfKey string

const (
	ConfOxMsgBrokerUri                 ConfKey = "OX_MSGBROKER_URI"
	ConfOxMsgBrokerUser                ConfKey = "OX_MSGBROKER_USER"
	ConfOxMsgBrokerPwd                 ConfKey = "OX_MSGBROKER_PWD"
	ConfOxMsgBrokerInsecureSkipVerify  ConfKey = "OX_MSGBROKER_INSECURE_SKIP_VERIFY"
	ConfOxMsgBrokerClientId            ConfKey = "OX_MSGBROKER_CLIENT_ID"
	ConfOxMsgBrokerQoS                 ConfKey = "OX_MSGBROKER_QoS"
	ConfOxMsgBrokerTopic               ConfKey = "OX_MSGBROKER_TOPIC"
	ConfOxMsgBrokerShutdownGracePeriod ConfKey = "OX_MSGBROKER_SHUTDOWN_GRACE_PERIOD"
)

type Conf struct {
}

func NewConf() *Conf {
	return &Conf{}
}

func (c *Conf) get(key ConfKey) string {
	return os.Getenv(string(key))
}

func (c *Conf) getConfOxMsgBrokerUri() string {
	return c.getValue(ConfOxMsgBrokerUri)
}

func (c *Conf) getConfOxMsgBrokerClientId() string {
	return c.getValue(ConfOxMsgBrokerClientId)
}

func (c *Conf) getConfOxMsgBrokerQoS() int {
	if len(c.getValue(ConfOxMsgBrokerQoS)) > 0 {
		i, err := strconv.Atoi(c.getValue(ConfOxMsgBrokerQoS))
		if err != nil {
			fmt.Printf("ERROR: invalid value for variable %s, value [%s] \n", ConfOxMsgBrokerQoS, c.getValue(ConfOxMsgBrokerQoS))
			os.Exit(1)
		}
		return i
	} else {
		return 0
	}
}

func (c *Conf) getConfOxMsgBrokerTopic() string {
	return c.getValue(ConfOxMsgBrokerTopic)
}

func (c *Conf) getConfOxMsgBrokerUser() string {
	return c.getValue(ConfOxMsgBrokerUser)
}

func (c *Conf) getConfOxMsgBrokerPwd() string {
	return c.getValue(ConfOxMsgBrokerPwd)
}

func (c *Conf) getConfOxMsgBrokerInsecureSkipVerify() bool {
	b, err := strconv.ParseBool(c.getValue(ConfOxMsgBrokerInsecureSkipVerify))
	if err != nil {
		fmt.Printf("ERROR: invalid value for variable %s \n", ConfOxMsgBrokerInsecureSkipVerify)
		os.Exit(1)
	}
	return b
}

func (c *Conf) getConfOxMsgBrokerShutdownGracePeriod() uint {
	if len(c.getValue(ConfOxMsgBrokerShutdownGracePeriod)) > 0 {
		i, err := strconv.ParseUint(c.getValue(ConfOxMsgBrokerShutdownGracePeriod), 0, 32)
		if err != nil {
			fmt.Printf("ERROR: invalid value for variable %s, value [%s] \n", ConfOxMsgBrokerShutdownGracePeriod, c.getValue(ConfOxMsgBrokerShutdownGracePeriod))
			return 250
		}
		return uint(i)
	} else {
		return 250
	}
}

func (c *Conf) getValue(key ConfKey) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		core.Debug("\n Warning: variable %s not defined\n", key)
	}
	return value
}

func IsMqttConfigAvailable() bool {
	if len(os.Getenv(string(ConfOxMsgBrokerUri))) > 0 {
		return true
	}
	return false
}
