/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package runner

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"path/filepath"
	"time"
)

// runs functions defined in packages or sources containing build.yaml within a runtime
type Runner struct {
	buildFile *data.BuildFile
	path      string
	artHome   string
}

func NewFromPath(path, artHome string) (*Runner, error) {
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
		artHome:   artHome,
	}, nil
}

func New() (*Runner, error) {
	return new(Runner), nil
}

func (r *Runner) RunC(fxName string, interactive bool, env *merge.Envar, network string) error {
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
	runtime = core.QualifyRuntime(runtime)
	// generate a unique name for the running container
	containerName := fmt.Sprintf("art-runc-%s-%s", core.Encode(fxName), core.RandomString(8))
	// if insputs are defined for the function then survey for data
	i := data.SurveyInputFromBuildFile(fxName, r.buildFile, true, false, env, r.artHome)
	// merge the collected input with the current environment
	env.Merge(i.Env())
	core.Debug(fmt.Sprintf("env vars passed to container:\n%s\n", env.String()))
	// launch a container with a bind mount to the path where the build.yaml is located
	err := runBuildFileFx(runtime, fxName, r.path, containerName, network, env, r.artHome)
	if err != nil {
		removeContainer(containerName)
		return err
	}
	// wait for the container to complete its task
	for isRunning(containerName) {
		time.Sleep(500 * time.Millisecond)
	}
	removeContainer(containerName)
	return nil
}

func (r *Runner) ExeC(packageName, fxName, credentials, network string, interactive bool, env *merge.Envar) error {
	var runtime string
	name, _ := core.ParseName(packageName)
	// get a local registry handle
	local := registry.NewLocalRegistry(r.artHome)
	// ensure the package is in the local registry
	local.Pull(name, credentials)
	// get the package manifest
	m := local.GetManifest(name)
	// if the manifest exports the function
	if isExported(m, fxName) {
		// get the runtime to use from the manifest function
		fx := m.Fx(fxName)
		// if the runtime is defined at the function level
		if len(fx.Runtime) > 0 {
			// use the function level runtime
			runtime = fx.Runtime
		} else if len(m.Runtime) > 0 {
			// if not use the manifest level runtime
			runtime = m.Runtime
		} else {
			return fmt.Errorf("runtime attribute is required in manifest for package '%s'", name)
		}
		runtime = core.QualifyRuntime(runtime)
		// interactively survey for required input via CLI
		input := data.SurveyInputFromManifest(name.Group, name.Name, "", name.Domain, fxName, m, interactive, false, env, r.artHome)
		// merge the collected input with the current environment without adding the PGP keys (they must be present locally)
		env.Merge(input.Env())
		// get registry credentials
		uname, pwd := core.RegUserPwd(credentials)
		// create a random container name
		containerName := fmt.Sprintf("art-exec-%s", core.RandomString(8))
		// launch a container with a bind mount to the artisan registry only
		err := runPackageFx(runtime, packageName, fxName, containerName, uname, pwd, network, env, r.artHome)
		if err != nil {
			removeContainer(containerName)
			return err
		}
		// wait for the container to complete its task
		for isRunning(containerName) {
			time.Sleep(500 * time.Millisecond)
		}
		removeContainer(containerName)
		return nil
	} else {
		core.RaiseErr("the function '%s' is not defined in the package manifest, check that it has been exported in the build profile", fxName)
	}
	return nil
}
