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
	l.Remove([]*core.ArtieName{core.ParseName("874484")})
}

func TestTag(t *testing.T) {
	l := NewFileRegistry()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Tag(core.ParseName("localhost:8081/gatblau/boot"), core.ParseName("boot:11"))
}
