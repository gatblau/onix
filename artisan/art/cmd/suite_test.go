package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/artisan/runner"
	"testing"
)

func Test(t *testing.T) {
	packageName := "artisan-registry-amosonline-aws-01-sapgatewaycd.apps.amosds.amosonline.io/recipe/java-quarkus"
	fxName := "setup"
	// create an instance of the runner
	run, err := runner.New()
	core.CheckErr(err, "cannot initialise runner")
	env, err := core.NewEnVarFromFile(".env")
	if err != nil {
		fmt.Printf("cannot load env file: %s\n", err.Error())
		t.FailNow()
	}
	// launch a runtime to execute the function
	err = run.ExeC(packageName, fxName, "admin:nxrpsap", false, env)
	i18n.Err(err, i18n.ERR_CANT_EXEC_FUNC_IN_PACKAGE, fxName, packageName)
}

func TestRunC(t *testing.T) {
	run, err := runner.NewFromPath(".")
	core.CheckErr(err, "cannot initialise runner")
	err = run.RunC("deploy", false, core.NewEnVarFromSlice([]string{}))
}

func TestMerge(t *testing.T) {
	args := []string{"artisan-registry-amosonline-aws-01-sapgatewaycd.apps.amosds.amosonline.io/recipe/java-quarkus", "setup"}
	var input *data.Input
	if len(args) > 0 && len(args) < 3 {
		name, err := core.ParseName(args[0])
		core.CheckErr(err, "invalid package name: %s", name)
		local := registry.NewLocalRegistry()
		manifest := local.GetManifest(name)
		if len(args) == 2 {
			fxName := args[1]
			fx := manifest.Fx(fxName)
			input = fx.Input
		} else {
			for i, function := range manifest.Functions {
				if i == 0 {
					input = function.Input
				} else {
					input.Merge(function.Input)
				}
			}
		}
		// add the credentials to download the package
		input.SurveyRegistryCreds(name.Group, name.Name, "", name.Domain, false, true, core.NewEnVarFromSlice([]string{}))
	}
}

func TestPush(t *testing.T) {
	reg := registry.NewLocalRegistry()
	name, err := core.ParseName("localhost:8082/artisan")
	if err != nil {
		t.FailNow()
	}
	reg.Push(name, "admin:admin")
}

func TestPull(t *testing.T) {
	reg := registry.NewLocalRegistry()
	name, err := core.ParseName("localhost:8082/gatblau/tools/artisan")
	if err != nil {
		t.FailNow()
	}
	reg.Pull(name, "admin:admin")
}
