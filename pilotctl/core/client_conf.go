package core

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// ClientConf client configuration information
type ClientConf struct {
	// the base URI for the service
	BaseURI string
	// disables TLS certificate verification
	InsecureSkipVerify bool
	// the user username for Basic and OpenId user authentication
	Username string
	// the user password for Basic and OpenId user authentication
	Password string
	// time out
	Timeout time.Duration
}

// gets the authentication token based on the authentication mode selected
func (cfg *ClientConf) getAuthToken() (string, error) {
	return cfg.basicToken(cfg.Username, cfg.Password), nil
}

// validates the client configuration
func checkConf(cfg *ClientConf) error {
	if len(cfg.BaseURI) == 0 {
		return errors.New("BaseURI is not defined")
	}
	// if the protocol is not specified, the add http as default
	// this is to avoid the server producing empty responses if no protocol is specified in the URI
	if !strings.HasPrefix(strings.ToLower(cfg.BaseURI), "http") {
		log.Warn().Msgf("no protocol defined for Onix URI '%s', 'http://' will be added to it", cfg.BaseURI)
		cfg.BaseURI = fmt.Sprintf("http://%s", cfg.BaseURI)
	}

	if len(cfg.Username) == 0 {
		return errors.New("username is not defined")
	}
	if len(cfg.Password) == 0 {
		return errors.New("password is not defined")
	}

	// if timeout is zero, it never timeout so is not good
	if cfg.Timeout == 0*time.Second {
		// set a default timeout of 30 secs
		cfg.Timeout = 30 * time.Second
	}
	return nil
}

// creates a new Basic Authentication Token
func (cfg *ClientConf) basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}
