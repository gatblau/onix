package test

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/flow"
	"testing"
)

func Test(t *testing.T) {
	env, _ := core.NewEnVarFromFile(".env")
	// loads a bare flow from the path
	flow, err := flow.NewWithEnv("setup_flow_bare.yaml", ".", env)
	core.CheckErr(err, "cannot load bare flow")
	// merges input, surveying for required data if in interactive mode
	err = flow.Merge(false)
	core.CheckErr(err, "cannot merge bare flow")
	// marshals the flow to YAML
	json, err := flow.JsonString()
	core.CheckErr(err, "cannot marshal bare flow")
	// print to stdout
	fmt.Println(json)
	flow.SaveYAML()
}
