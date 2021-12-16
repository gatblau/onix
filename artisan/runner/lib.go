package runner

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bufio"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// launch a container and mount the current directory on the host machine into the container
// the current directory must contain a build.yaml file where fxName is defined
func runBuildFileFx(runtimeName, fxName, dir, containerName string, env *merge.Envar) error {
	// if the OS is linux and the user id is not 100,000,000, it cannot continue
	if isWrong, msg := core.WrongUserId(); isWrong {
		// print error
		core.RaiseErr("%s\n", msg)
	}
	// check the local registry path has not been created by the root user othewise the runtime will error
	registryPath := core.RegistryPath()
	if runtime.GOOS == "linux" && strings.HasPrefix(registryPath, "//") {
		// in linux if the user is not root but the local registry folder is owned by the root user, then
		// the registry path in a runtime will start with two consecutive forward slashes
		core.RaiseErr("cannot continue, the local registry folder is owned by root\n" +
			"ensure it is owned by UID=100000000 for the runtime to work")
	}
	if env == nil {
		env = merge.NewEnVarFromSlice([]string{})
	}
	// determine which container tool is available in the host
	tool, err := containerCmd()
	if err != nil {
		return err
	}
	// add runtime vars
	env.Add("OXART_FX_NAME", fxName)
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
			if _, ok = exitErr.Sys().(syscall.WaitStatus); ok {
				return exitErr
			}
		}
		return err
	}
	return nil
}

// launch a container and execute a package function
func runPackageFx(runtimeName, packageName, fxName, containerName, artRegistryUser, artRegistryPwd string, env *merge.Envar) error {
	// if the OS is linux and the user id is not 100,000,000, it cannot continue
	if isWrong, msg := core.WrongUserId(); isWrong {
		// print warning
		fmt.Println(msg)
		os.Exit(1)
	}
	// determine which container tool is available in the host
	tool, err := containerCmd()
	if err != nil {
		return err
	}
	// add add runtime vars
	env.Add("OXART_PACKAGE_NAME", packageName)
	env.Add("OXART_FX_NAME", fxName)
	env.Add("OXART_REG_USER", artRegistryUser)
	env.Add("OXART_REG_PWD", artRegistryPwd)
	// create a slice with docker run args
	args := toContainerArgs(runtimeName, "", containerName, env)
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
			if _, ok = exitErr.Sys().(syscall.WaitStatus); ok {
				return exitErr
			}
		}
		return err
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
func toContainerArgs(imageName, dir, containerName string, env *merge.Envar) []string {
	var result = []string{"run", "--name", containerName}
	vars := env.Slice()
	for _, v := range vars {
		result = append(result, "-e")
		result = append(result, v)
	}
	// create bind mounts
	// note: in order to allow for art runc command to access host mounted files in linux with selinux enabled, a :Z label
	// is added to the volume see https://docs.docker.com/storage/bind-mounts/#configure-the-selinux-label
	// Z modify the selinux label of the host file or directory being mounted into the container indicating that the
	// bind mount content is private and unshared.
	if len(dir) > 0 {
		// add a bind mount for the current folder to the /workspace/source in the runtime
		result = append(result, "-v")
		result = append(result, fmt.Sprintf("%s:%s", dir, "/workspace/source:Z"))
	}
	// add a bind mount for the artisan registry folder
	result = append(result, "-v")
	// note: mind the location of the mount in the runtime must align with its user home!
	result = append(result, fmt.Sprintf("%s:%s", core.RegistryPath(), "/home/runtime/.artisan:Z"))
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
		core.InfoLogger.Printf("%s\n", string(out))
		core.CheckErr(err, "cannot remove temporary container %s", containerName)
	}
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
