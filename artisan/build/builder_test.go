package build

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func TestBuildContentOnly(t *testing.T) {
	builder := NewBuilder("")
	name, _ := core.ParseName("localhost:8080/lib/test1:1")
	builder.Build("", "", "", name, "", false, false, "test/test")
}
