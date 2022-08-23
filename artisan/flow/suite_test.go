package flow

import (
	"github.com/gatblau/onix/artisan/merge"
	"testing"
)

func Test(t *testing.T) {
	e, _ := merge.NewEnVarFromFile(".env")
	f, _ := NewWithEnv("ci_flow_bare.yaml", ".", e, "")
	f.SaveOnixJSON()
}
