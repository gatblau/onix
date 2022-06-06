package core

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	DoormanProxyUser         = "DOORMAN_PROXY_USER"
	DoormanProxyPwd          = "DOORMAN_PROXY_PWD"
	DoormanNotificationURI   = "DOORMAN_NOTIFICATION_URI"
	DoormanNotificationUser  = "DOORMAN_NOTIFICATION_USER"
	DoormanNotificationPwd   = "DOORMAN_NOTIFICATION_PWD"
	OxWapiUri                = "OX_WAPI_URI"
	OxWapiUser               = "OX_WAPI_USER"
	OxWapiPwd                = "OX_WAPI_PWD"
	OxWapiInsecureSkipVerify = "OX_WAPI_INSECURE_SKIP_VERIFY"
	ArtRegUser               = "ART_REG_USER"
	ArtRegPwd                = "ART_REG_PWD"
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

func getBoolean(key string) (bool, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return false, fmt.Errorf("variable %s is required and not defined", key)
	}
	return strconv.ParseBool(value)
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

func GetOxWapiUri() (string, error) {
	return getString(OxWapiUri)
}

func GetOxWapiUser() (string, error) {
	return getString(OxWapiUser)
}

func GetOxWapiPwd() (string, error) {
	return getString(OxWapiPwd)
}

func GetOxWapiInsecureSkipVerify() (bool, error) {
	return getBoolean(OxWapiInsecureSkipVerify)
}

func GetArRegUser() (string, error) {
	return getString(ArtRegUser)
}

func GetArRegPwd() (string, error) {
	return getString(ArtRegPwd)
}
