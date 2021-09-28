// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "gatblau",
            "url": "http://onix.gatblau.org/",
            "email": "onix@gatblau.org"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Checks that the HTTP server is listening on the required port.\nUse a liveliness probe.\nIt does not guarantee the server is ready to accept calls.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "General"
                ],
                "summary": "Check that the HTTP API is live",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/events": {
            "get": {
                "description": "Returns a list of syslog entries following the specified filter",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Query"
                ],
                "summary": "Get filtered events",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the organisation of the device where the syslog entry was created",
                        "name": "og",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the organisation of the device where the syslog entry was created",
                        "name": "or",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the area of the device where the syslog entry was created",
                        "name": "ar",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the location of the device where the syslog entry was created",
                        "name": "lo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "syslog entry tag",
                        "name": "tag",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the syslog entry priority",
                        "name": "pri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the syslog entry severity",
                        "name": "sev",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "the syslog entry time following the format ddMMyyyyHHmmSS",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "submits syslog events to be persisted for further use",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Receiver"
                ],
                "summary": "Submit Syslog Events",
                "parameters": [
                    {
                        "description": "the events to submit",
                        "name": "command",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.Events"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.Event": {
            "type": "object",
            "properties": {
                "area": {
                    "type": "string"
                },
                "boot_time": {
                    "type": "string"
                },
                "client": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "event_id": {
                    "type": "string"
                },
                "facility": {
                    "type": "integer"
                },
                "host_address": {
                    "type": "string"
                },
                "hostname": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "machine_id": {
                    "type": "string"
                },
                "org": {
                    "type": "string"
                },
                "org_group": {
                    "type": "string"
                },
                "priority": {
                    "type": "integer"
                },
                "severity": {
                    "type": "integer"
                },
                "tag": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "tls_peer": {
                    "type": "string"
                }
            }
        },
        "types.Events": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Event"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0.4",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "MongoDB Event Receiver for Pilot Control",
	Description: "Onix Config Manager Event Receiver for Pilot Control using MongoDb",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
