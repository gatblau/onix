package registry

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func Test(t *testing.T) {
	n, _ := core.ParseName("artisan-registry-amosonline-aws-01-sapgatewaycd.apps.amosds.amosonline.io/gatblau/artisan")
	l := NewLocalRegistry()
	l.Push(n, "admin:nxrpsap", false)
}
