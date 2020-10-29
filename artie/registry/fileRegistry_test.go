package registry

import (
	"github.com/gatblau/onix/artie/core"
	"testing"
)

func TestUpload(t *testing.T) {
	named := core.ParseName("localhost:8081/gatblau/boot")
	l := NewFileRegistry()
	r := NewNexus3Registry()
	l.Push(named, r, "admin:admin123")
}

func TestRemove(t *testing.T) {
	l := NewFileRegistry()
	l.Remove([]*core.ArtieName{core.ParseName("localhost:8081/gatblau/boot:v32")})
}
