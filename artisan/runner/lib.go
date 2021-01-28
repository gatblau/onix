/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package runner

import (
	"bufio"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// launch a container and mount the current directory on the host machine into the container
// the current directory must contain a build.yaml file where fxName is defined
func runBuildFileFx(runtimeName, fxName, dir, containerName string, env *core.Envar) error {
	// if wrong UID
	if ok, msg := wrongUserId(); !ok {
		// print warning
		fmt.Println(msg)
	}
	if env == nil {
		env = core.NewEnVarFromSlice([]string{})
	}
	// determine which container tool is available in the host
	tool, err := containerCmd()
	if err != nil {
		return err
	}
	// add runtime vars
	env.Add("FX_NAME", fxName)
	// get the docker run arguments
	args := toContainerArgs(runtimeName, dir, containerName, env)
	// launch the container with an art exec command
	cmd := exec.Command(tool, args...)
	core.Debug("! launching runtime: %s %s\n", tool, strings.Join(args, " "))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed creating command stdoutpipe: %s", err)
	}
	defer func() {
		_ = stdout.Close()
	}()
	stdoutReader := bufio.NewReader(stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed creating command stderrpipe: %s", err)
	}
	defer func() {
		_ = stderr.Close()
	}()
	stderrReader := bufio.NewReader(stderr)

	if err = cmd.Start(); err != nil {
		return err
	}

	go handleReader(stdoutReader)
	go handleReader(stderrReader)

	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if _, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return exitErr
			}
		}
		return err
	}
	return nil
}

// launch a container and execute a package function
func runPackageFx(runtimeName, packageName, fxName, dir, containerName, artRegistryUser, artRegistryPwd string, env *core.Envar) error {
	// if wrong UID
	if ok, msg := wrongUserId(); !ok {
		// print warning
		fmt.Println(msg)
	}
	// determine which container tool is available in the host
	tool, err := containerCmd()
	if err != nil {
		return err
	}
	// add runtime vars
	env.Add("PACKAGE_NAME", packageName)
	env.Add("FX_NAME", fxName)
	env.Add("ART_REG_USER", artRegistryUser)
	env.Add("ART_REG_PWD", artRegistryPwd)
	// create a slice with docker run args
	args := toContainerArgs(runtimeName, dir, containerName, env)
	// launch the container with an art exec command
	cmd := exec.Command(tool, args...)
	core.Debug("! launching runtime: %s %s\n", tool, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot launch container: %s", err)
	}
	return nil
}

// return the command to run to launch a container
func containerCmd() (string, error) {
	if isCmdAvailable("docker") {
		return "docker", nil
	} else if isCmdAvailable("podman") {
		return "podman", nil
	}
	return "", fmt.Errorf("either podman or docker is required to launch a container")
}

// checks if a command is available
func isCmdAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// return an array of environment variable arguments to pass to docker
func toContainerArgs(imageName, dir, containerName string, env *core.Envar) []string {
	var result = []string{"run", "--name", containerName} // , "-d", "--rm"
	vars := env.Slice()
	for _, v := range vars {
		result = append(result, "-e")
		result = append(result, v)
	}
	if len(dir) > 0 {
		// add a bind mount for the current folder to the /workspace/source in the runtime
		result = append(result, "-v")
		result = append(result, fmt.Sprintf("%s:%s", dir, "/workspace/source"))
		// add a bind mount for the artisan registry folder
		result = append(result, "-v")
		result = append(result, fmt.Sprintf("%s:%s", core.RegistryPath(), "/.artisan"))
	}
	result = append(result, imageName)
	return result
}

func handleReader(reader *bufio.Reader) {
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(str)
	}
}

func isRunning(containerName string) bool {
	tool, err := containerCmd()
	core.CheckErr(err, "")
	cmd := exec.Command(tool, "container", "inspect", "-f", "'{{.State.Running}}'", containerName)
	out, _ := cmd.Output()
	if strings.Contains(strings.ToLower(string(out)), "error") {
		return false
	}
	running, err := strconv.ParseBool(string(out))
	if err != nil {
		return false
	}
	return running
}

// removes a docker container
func removeContainer(containerName string) {
	tool, err := containerCmd()
	core.CheckErr(err, "")
	rm := exec.Command(tool, "rm", containerName)
	out, err := rm.Output()
	if err != nil {
		core.Msg(string(out))
		core.CheckErr(err, "cannot remove temporary container %s", containerName)
	}
}

// check the user id is correct for bind mounts
func wrongUserId() (bool, string) {
	// if running in linux docker does not run in a VM and uid and gid of bind mounts must
	// match the one in the runtime
	if runtime.GOOS == "linux" {
		// if the user id is not the id of the runtime user
		if os.Geteuid() != 100000000 {
			return true, fmt.Sprintf(`
WARNING! The UID of the running user does not match the one in the runtime.
This can render the bind mounts unusable and red/write errors can ocurr if the process tries to read / or wirte to them.
If you intend to use this command in a linux machine ensure it is run by a user with UID = 100000000.
For example, assuming the user is call "runtime"", you can:
	- create a user with UID 100000000 as follows:
      $ useradd -u 100000000 -g 100000000 runtime
    - create a group with GID 100000000 as follows:
      $ groupadd -g 100000000 -o runtime
	- log a the "runtime" user before running the art command
	- if using docker, add the runtime user to the docker group
      $ sudo usermod -aG docker runtime
`)
		}
	}
	return false, ""
}

// check the the specified function is in the manifest
func isExported(m *data.Manifest, fx string) bool {
	for _, function := range m.Functions {
		if function.Name == fx {
			return true
		}
	}
	return false
}
