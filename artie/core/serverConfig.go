/*
  Onix ServerConfig Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	VarMetricsEnabled = "OXA_METRICS_ENABLED"
	VarSwaggerEnabled = "OXA_SWAGGER_ENABLED"
	VarHTTPPort       = "OXA_HTTP_PORT"
	VarHTTPUname      = "OXA_HTTP_UNAME"
	VarHTTPPwd        = "OXA_HTTP_PWD"
	VarBackendType    = "OXA_HTTP_BACKEND"
	VarBackendDomain  = "OXA_HTTP_BACKEND_DOMAIN"
)

type Backend int

const (
	FileSystem Backend = iota
	S3
	Nexus3
	Artifactory
	NotRecognized
)

func (b Backend) String() string {
	switch b {
	case FileSystem:
		return "FileSystem"
	case S3:
		return "S3"
	case Nexus3:
		return "Nexus3"
	case Artifactory:
		return "Artifactory"
	default:
		return "NotRecognized"
	}
}

func ParseBackend(value string) Backend {
	switch strings.ToLower(value) {
	case "filesystem":
		return FileSystem
	case "s3":
		return S3
	case "nexus3":
		return Nexus3
	case "artifactory":
		return Artifactory
	default:
		return NotRecognized
	}
}

type ServerConfig struct {
}

func (c *ServerConfig) Backend() Backend {
	value := os.Getenv(VarBackendType)
	if len(value) == 0 {
		// defaults to Nexus3
		return Nexus3
	}
	backend := ParseBackend(value)
	if backend == NotRecognized {
		RaiseErr("backend '%s' is not a valid backend", value)
	}
	return backend
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
