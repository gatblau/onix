package core

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const (
	DoormanProxyUser = "DOORMAN_PROXY_USER"
	DoormanProxyPwd  = "DOORMAN_PROXY_PWD"
)

func init() {
	// load env vars from file if present
	godotenv.Load("doorman.env")
}

func getString(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("variable %s is required and not defined", key)
	}
	return value, nil
}

func GetProxyUser() (string, error) {
	return getString(DoormanProxyUser)
}

func GetProxyPwd() (string, error) {
	return getString(DoormanProxyPwd)
}
