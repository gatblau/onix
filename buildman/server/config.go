package server

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
)

const (
	VarMetricsEnabled = "OXB_METRICS_ENABLED"
	VarSwaggerEnabled = "OXB_SWAGGER_ENABLED"
	VarHTTPPort       = "OXB_HTTP_PORT"
	VarHTTPUname      = "OXB_HTTP_UNAME"
	VarHTTPPwd        = "OXB_HTTP_PWD"
)

type Config struct {
}

func (c *Config) SwaggerEnabled() bool {
	return c.getBoolean(VarSwaggerEnabled, true)
}

func (c *Config) MetricsEnabled() bool {
	return c.getBoolean(VarMetricsEnabled, true)
}

func (c *Config) HttpPort() string {
	return c.getString(VarHTTPPort, "8082")
}

func (c *Config) HttpUser() string {
	return c.getString(VarHTTPUname, "admin")
}

func (c *Config) HttpPwd() string {
	return c.getString(VarHTTPPwd, "admin")
}

func (c *Config) getBoolean(varName string, defaultValue bool) bool {
	value := os.Getenv(varName)
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		// set as default value
		enabled = defaultValue
	}
	return enabled
}

func (c *Config) getString(varName string, defaultValue string) string {
	value := os.Getenv(varName)
	if len(value) == 0 {
		// set as default value
		value = defaultValue
	}
	return value
}

func (c *Config) BasicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.HttpUser(), c.HttpPwd()))))
}
