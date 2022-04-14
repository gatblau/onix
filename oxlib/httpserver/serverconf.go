package httpserver

/*
  Onix Config Manager - Http Client
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	VarMetricsEnabled         = "OX_METRICS_ENABLED"
	VarSwaggerEnabled         = "OX_SWAGGER_ENABLED"
	VarHTTPPort               = "OX_HTTP_PORT"
	VarHTTPUname              = "OX_HTTP_UNAME"
	VarHTTPPwd                = "OX_HTTP_PWD"
	VarHTTPRealm              = "OX_HTTP_REALM"
	VarHTTPUploadPayloadLimit = "OX_HTTP_UPLOAD_LIMIT"
	VarHTTPUploadInMemSize    = "OX_HTTP_UPLOAD_IN_MEM_SIZE"
)

type ServerConfig struct{}

func (c *ServerConfig) MetricsEnabled() bool {
	return c.getBoolean(VarMetricsEnabled, true)
}

func (c *ServerConfig) SwaggerEnabled() bool {
	return c.getBoolean(VarSwaggerEnabled, true)
}

func (c *ServerConfig) HttpPort() string {
	return c.getString(VarHTTPPort, "8080")
}

func (c *ServerConfig) HttpUser() string {
	return c.getString(VarHTTPUname, "admin")
}

func (c *ServerConfig) HttpPwd() string {
	return c.getString(VarHTTPPwd, "adm1n")
}

func (c *ServerConfig) HttpRealm() string {
	return c.getString(VarHTTPRealm, "*")
}

func (c *ServerConfig) getBoolean(varName string, defaultValue bool) bool {
	value := os.Getenv(varName)
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		// set as default value
		enabled = defaultValue
	}
	return enabled
}

func (c *ServerConfig) getString(varName string, defaultValue string) string {
	value := os.Getenv(varName)
	if len(value) == 0 {
		// set as default value
		value = defaultValue
	}
	return value
}

func (c *ServerConfig) BasicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.HttpUser(), c.HttpPwd()))))
}

func (c *ServerConfig) HttpUploadLimit() int64 {
	limit, err := strconv.ParseInt(c.getString(VarHTTPUploadPayloadLimit, "250"), 0, 0)
	if err != nil {
		log.Fatalf("invalid upload limit specified")
	}
	return limit
}

func (c *ServerConfig) HttpUploadInMemorySize() int64 {
	limit, err := strconv.ParseInt(c.getString(VarHTTPUploadInMemSize, "150"), 0, 0)
	if err != nil {
		log.Fatalf("invalid upload limit specified")
	}
	return limit
}
