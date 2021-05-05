package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import "strings"

type BackendType int

const (
	FileSystem BackendType = iota
	S3
	Nexus3
	Artifactory
	NotRecognized
)

func (b BackendType) String() string {
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

func ParseBackend(value string) BackendType {
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
