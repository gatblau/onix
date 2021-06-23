package registry

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func Test(t *testing.T) {
	r := NewLocalRegistry()
	name, _ := core.ParseName("localhost:8082/aps/anthos")
	r.Push(name, "admin:admin")
}
