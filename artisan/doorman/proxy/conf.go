/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"fmt"
	"os"
	"strconv"
)

func getEmailFrom() (string, error) {
	value := os.Getenv("DPROXY_EMAIL_FROM")
	if len(value) == 0 {
		return "", fmt.Errorf("variable DPROXY_EMAIL_FROM is required and not defined")
	}
	return value, nil
}

func getEmailFromPwd() (string, error) {
	value := os.Getenv("DPROXY_EMAIL_FROM_PWD")
	if len(value) == 0 {
		return "", fmt.Errorf("variable DPROXY_EMAIL_FROM_PWD is required and not defined")
	}
	return value, nil
}

func getSmtpHost() (string, error) {
	value := os.Getenv("DPROXY_SMTP_HOST")
	if len(value) == 0 {
		return "", fmt.Errorf("variable DPROXY_SMTP_HOST is required and not defined")
	}
	return value, nil
}

func getSmtpPort() (int, error) {
	value := os.Getenv("DPROXY_SMTP_PORT")
	if len(value) == 0 {
		return -1, fmt.Errorf("variable DPROXY_SMTP_PORT is required and not defined")
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		return -1, fmt.Errorf("invalid DPROXY_SMTP_PORT value: %s; %s", value, err)
	}
	return port, nil
}
