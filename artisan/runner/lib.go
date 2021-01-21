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
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// launch a container and mount the current directory on the host machine into the container
// the current directory must contain a build.yaml file where fxName is defined
func launchContainerWithBindMount(runtimeName, fxName, dir, containerName string, env *core.Envar) error {
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
	args := toCmdArgs(runtimeName, dir, containerName, env)
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
func runPackageContainer(runtimeName, packageName, fxName, artRegistryUser, artRegistryPwd string, env *core.Envar) error {
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
	containerName := fmt.Sprintf("art-run-%s", core.RandomString(8))
	// launch the container with an art exec command
	cmd := exec.Command(tool, toCmdArgs(runtimeName, "", containerName, env)...)
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
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// return an array of environment variable arguments to pass to docker
func toCmdArgs(imageName, dir, containerName string, env *core.Envar) []string {
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

// prints the logs of the container
func printLogs(containerName string) {
	tool, err := containerCmd()
	core.CheckErr(err, "")
	cmd := exec.Command(tool, "logs", containerName)
	logs, _ := cmd.Output()
	log.Printf("%s\n", logs)
}

// removes a docker container
func removeContainer(containerName string) {
	tool, err := containerCmd()
	core.CheckErr(err, "")
	_ = exec.Command(tool, "rm", containerName)
}
