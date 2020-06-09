// Onix Database Manager Web API.
//
// This service makes some of DbMan's cli commands available via HTTP.
// Its main purpose is to facilitate key database maintenance operations in containerised environments.
//
// NOTE: Any commands which are meant to modify configuration sets are not available from this service because it is
// assumed to run from a container. Therefore, service configuration can only be set via environment variables.
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.4
//     License: Apache License, Version 2.0 http://www.apache.org/licenses/LICENSE-2.0
//     Contact: gatblau<onix@gatblau.org>
//
//     Produces:
//     - text/plain
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package docs
