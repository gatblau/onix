/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package runner

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"path/filepath"
	"strings"
	"time"
)

// runs functions defined in packages or sources containing build.yaml within a runtime
type Runner struct {
	buildFile *data.BuildFile
	path      string
}

func NewFromPath(path string) (*Runner, error) {
	if len(path) == 0 {
		path = "."
	}
	path = core.ToAbs(path)
	bf, err := data.LoadBuildFile(filepath.Join(path, "build.yaml"))
	if err != nil {
		return nil, fmt.Errorf("cannot load build file: %s", err)
	}
	return &Runner{
		buildFile: bf,
		path:      path,
	}, nil
}

func (r *Runner) RunC(fxName string) error {
	var runtime string
	fx := r.buildFile.Fx(fxName)
	// if the runtime is defined at the function level
	if len(fx.Runtime) > 0 {
		// use the function level runtime
		runtime = fx.Runtime
	} else if len(r.buildFile.Runtime) > 0 {
		// if not use the build file level runtime
		runtime = r.buildFile.Runtime
	} else {
		return fmt.Errorf("runtime attribute is required in build.yaml within %s", r.path)
	}
	// completes name if the short form is used
	runtime = format(runtime)
	// generate a unique name for the running container
	containerName := fmt.Sprintf("art-%s-%s", core.Encode(fxName), core.RandomString(8))
	// collect any input required to run the function
	env := core.NewEnVarFromSlice([]string{})
	// interactively survey for required input via CLI
	input := data.SurveyInputFromBuildFile(fxName, r.buildFile, true)
	// if there are input data
	if input != nil {
		// add the variables to the environment
		for _, variable := range input.Var {
			env.Add(variable.Name, variable.Value)
		}
		// add the secrets to the environment
		for _, secret := range input.Secret {
			env.Add(secret.Name, secret.Value)
		}
	}
	// launch a container with a bind mount to the path where the build.yaml is located
	err := launchContainerWithBindMount(runtime, fxName, r.path, containerName, env)
	if err != nil {
		return err
	}
	// wait for the container to complete its task
	for isRunning(containerName) {
		time.Sleep(500 * time.Millisecond)
	}
	removeContainer(containerName)
	return nil
}

func format(runtime string) string {
	// container images must be in lower case
	runtime = strings.ToLower(runtime)
	// if no repository is specified then assume artisan library at quay.io/artisan
	if !strings.ContainsAny(runtime, "/") {
		return fmt.Sprintf("quay.io/artisan/%s", runtime)
	}
	return runtime
}
