package registry

import (
	"github.com/gatblau/onix/artie/core"
	"testing"
)

func TestUpload(t *testing.T) {
	named := core.ParseName("localhost:8082/gatblau/boot")
	local := NewLocalAPI()
	remote := NewRemoteAPI(false)
	local.Push(named, remote, "admin:admin")
}

func TestRemove(t *testing.T) {
	l := NewLocalAPI()
	l.Remove([]*core.ArtieName{core.ParseName("874484")})
}

func TestTag(t *testing.T) {
	l := NewLocalAPI()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Tag(core.ParseName("localhost:8081/gatblau/boot"), core.ParseName("boot:11"))
}
