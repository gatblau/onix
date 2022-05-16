/*
  Onix Config Manager - Onix Library
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package resx

import "strings"

const (
	File UriType = iota
	S3
	S3S
	Http
	Https
	Ftp
	Ftps
	Unknown
)

type UriType int64

func (t UriType) String() string {
	switch t {
	case File:
		return "file"
	case S3:
		return "s3"
	case S3S:
		return "s3s"
	case Http:
		return "http"
	case Https:
		return "https"
	case Ftp:
		return "ftp"
	case Ftps:
		return "ftps"
	case Unknown:
		return "unknown"
	}
	return "unknown"
}

func ParseUriType(uri string) UriType {
	if strings.HasPrefix(uri, "http://") {
		return Http
	}
	if strings.HasPrefix(uri, "https://") {
		return Https
	}
	if strings.HasPrefix(uri, "s3://") {
		return S3
	}
	if strings.HasPrefix(uri, "s3s://") {
		return S3S
	}
	if strings.HasPrefix(uri, "ftp://") {
		return Ftp
	}
	if strings.HasPrefix(uri, "ftps://") {
		return Ftps
	}
	if !strings.Contains(uri, "://") {
		return File
	}
	return Unknown
}

func IsFile(uri string) bool {
	return ParseUriType(uri) == File
}
