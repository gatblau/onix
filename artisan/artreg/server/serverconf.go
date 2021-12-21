/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/artreg/backend"
	"github.com/gatblau/onix/artisan/core"
	"os"
	"strconv"
)

const (
	VarMetricsEnabled  = "OXA_METRICS_ENABLED"
	VarSwaggerEnabled  = "OXA_SWAGGER_ENABLED"
	VarHTTPPort        = "OXA_HTTP_PORT"
	VarHTTPUname       = "OXA_HTTP_UNAME"
	VarHTTPPwd         = "OXA_HTTP_PWD"
	VarBackendType     = "OXA_HTTP_BACKEND"
	VarBackendDomain   = "OXA_HTTP_BACKEND_DOMAIN"
	VarHTTPUploadLimit = "OXA_HTTP_UPLOAD_LIMIT"
)

type ServerConfig struct {
}

func (c *ServerConfig) Backend() backend.BackendType {
	value := os.Getenv(VarBackendType)
	if len(value) == 0 {
		// defaults to Nexus3
		return backend.Nexus3
	}
	be := backend.ParseBackend(value)
	if be == backend.NotRecognized {
		core.RaiseErr("backend '%s' is not a valid backend", value)
	}
	return be
}

func (c *ServerConfig) MetricsEnabled() bool {
	return c.getBoolean(VarMetricsEnabled, true)
}

func (c *ServerConfig) SwaggerEnabled() bool {
	return c.getBoolean(VarSwaggerEnabled, true)
}

func (c *ServerConfig) HttpPort() string {
	return c.getString(VarHTTPPort, "8082")
}

func (c *ServerConfig) HttpUser() string {
	return c.getString(VarHTTPUname, "admin")
}

func (c *ServerConfig) HttpUploadLimit() int64 {
	limit, err := strconv.ParseInt(c.getString(VarHTTPUploadLimit, "30"), 0, 0)
	core.CheckErr(err, "invalid upload limit specified")
	return limit
}

func (c *ServerConfig) BackendDomain() string {
	return c.getString(VarBackendDomain, "http://localhost:8081")
}

func (c *ServerConfig) HttpPwd() string {
	return c.getString(VarHTTPPwd, "admin")
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
