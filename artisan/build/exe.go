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
	// if in debug mode, print environment variables
	env.Debug()
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
		// only happens if the command exits with os.Exit(>0)
		// if this happens then the only error available is the exit error code
		// for this reason artisan exit with code 0 and fills the stderr buffer
		if exitErr, ok := err.(*exec.ExitError); ok {
			if _, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return "", fmt.Errorf("run command failed: '%s'\n%s (%s)", cmd, errbuf.String(), exitMsg(exitErr.ExitCode()))
			}
		}
		return "", err
	}

	// NOTE: I have observed that some programs exit with no error (code 0) but write to stderr instead of stdout
	// probably due to misuse of the print() function in golang or lack mistake
	// this condition can be found if we reach this point and the errbuf contains bytes
	// at this point I have to assume that as the exit code is 0 there is no actual error and whatever is in stderr
	// it should be in stdout, therefore code below
	if len(errbuf.String()) > 0 {
		// append to stdout
		outbuf.WriteString(errbuf.String())
		if core.InDebugMode() {
			// issue a warning to alert people just in case
			core.WarningLogger.Printf("command %s returned successfully but data was found in stderr. it is assumed that it is not an error and therefore, it has been added to stdout\n", cmd)
		}
	}

	return outbuf.String(), err
}
