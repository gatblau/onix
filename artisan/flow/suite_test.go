package flow

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func Test(t *testing.T) {
	e, _ := core.NewEnVarFromFile(".env")
	f, _ := NewWithEnv("ci_flow_bare.yaml", ".", e)
	f.SaveOnixJSON()
}
