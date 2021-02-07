package backend

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
