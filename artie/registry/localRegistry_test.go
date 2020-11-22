package registry

import (
	"github.com/gatblau/onix/artie/core"
	"testing"
)

func TestUpload(t *testing.T) {
	named := core.ParseName("localhost:8082/gatblau/artie:v10")
	local := NewLocalRegistry()
	local.Push(named, "admin:admin", false)
}

func TestRemove(t *testing.T) {
	l := NewLocalRegistry()
	l.Remove([]*core.ArtieName{core.ParseName("874484")})
}

func TestTag(t *testing.T) {
	l := NewLocalRegistry()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Tag(core.ParseName("localhost:8081/gatblau/boot"), core.ParseName("boot:11"))
}
