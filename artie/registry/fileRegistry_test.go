package registry

import (
	"github.com/gatblau/onix/artie/core"
	"log"
	"testing"
)

func TestUpload(t *testing.T) {
	named, err := core.ParseNormalizedNamed("localhost:8081/gatblau/boot")
	if err != nil {
		log.Fatal(err)
	}
	l := NewFileRegistry()
	r := NewNexus3Registry()
	l.Push(named, r, "admin:admin123")
}
