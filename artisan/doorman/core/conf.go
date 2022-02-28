package core

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const (
	DoormanProxyUser        = "DOORMAN_PROXY_USER"
	DoormanProxyPwd         = "DOORMAN_PROXY_PWD"
	DoormanNotificationURI  = "DOORMAN_NOTIFICATION_URI"
	DoormanNotificationUser = "DOORMAN_NOTIFICATION_USER"
	DoormanNotificationPwd  = "DOORMAN_NOTIFICATION_PWD"
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

func GetNotificationURI() (string, error) {
	return getString(DoormanNotificationURI)
}

func GetNotificationUser() (string, error) {
	return getString(DoormanNotificationUser)
}

func GetNotificationPwd() (string, error) {
	return getString(DoormanNotificationPwd)
}
