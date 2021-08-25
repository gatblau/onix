package build

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/mattn/go-shellwords"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// Exe executes a command and sends output and error streams to stdout and stderr
func Exe(cmd string, dir string, env *merge.Envar, interactive bool) (string, error) {
	if cmd == "" {
		return "", errors.New("no command provided")
	}
	// create a command parser
	p := shellwords.NewParser()
	// parse the command line
	cmdArr, err := p.Parse(cmd)
	if err != nil {
		return "", err
	}
	// if we are in windows
	if runtime.GOOS == "windows" {
		// prepend "cmd /C" to the command line
		cmdArr = append([]string{"cmd", "/C"}, cmdArr...)
		core.Debug("windows cmd => %s", cmdArr)
	}
	name := cmdArr[0]

	var args []string
	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}

	args, _ = core.MergeEnvironmentVars(args, env.Vars, interactive)

	// create the command to execute
	command := exec.Command(name, args...)
	// set the command working directory
	command.Dir = dir
	// set the command environment
	command.Env = env.Slice()
	// capture the command output and error streams in a buffer
	var outbuf, errbuf strings.Builder // or bytes.Buffer
	command.Stdout = &outbuf
	command.Stderr = &errbuf

	// start the execution of the command
	if err := command.Start(); err != nil {
		return "", err
	}

	// wait for the command to complete
	if err := command.Wait(); err != nil {
		// only happens if the command exits with code > 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			if _, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				core.RaiseErr("run command failed: '%s' - '%s'\n", cmd, exitErr.Error())
			}
		}
		return "", err
	}

	// if we have characters in the error buffer
	if len(errbuf.String()) > 0 {
		// create an error object
		err = fmt.Errorf(errbuf.String())
	} else {
		// otherwise, make the error nil
		err = nil
	}
	return outbuf.String(), err
}
