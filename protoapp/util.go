/*
*    Onix Boot - Copyright (c) 2020 by www.gatblau.org
*
*    Licensed under the Apache License, Version 2.0 (the "License");
*    you may not use this file except in compliance with the License.
*    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*    Unless required by applicable law or agreed to in writing, software distributed under
*    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
*    either express or implied.
*    See the License for the specific language governing permissions and limitations under the License.
*
*    Contributors to this project, hereby assign copyright in this code to the project,
*    to be licensed under the same terms as the rest of the code.
 */
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// executes a command
func execCmd(cmd []string) error {
	var c *exec.Cmd
	if len(cmd) == 1 {
		//nolint:gosec
		c = exec.Command(cmd[0])
	} else {
		//nolint:gosec
		c = exec.Command(cmd[0], cmd[1:]...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stdin
	return startCmd(c, 6)
}

// launch a command with a specified timeout period
func startCmd(cmd *exec.Cmd, timeoutSecs time.Duration) error {
	// execute the command asynchronously
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("executing %s failed: %s", cmd, err)
	}
	done := make(chan error)
	// launch a go routine to wait for the command to execute
	go func() {
		// send a message to the done channel if completed or error
		done <- cmd.Wait()
	}()
	// wait for the done channel
	select {
	case <-done:
		// command completed
	case <-time.After(timeoutSecs * time.Second):
		// command timed out after 6 secs
		return fmt.Errorf("executing '%s' timed out", cmd)
	}
	return nil
}

// cmdToRun := "/path/to/someCommand"
// args := []string{"someCommand", "arg1"}
// starts a process
func startProc(cmd string, args []string) (*os.Process, error) {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	return os.StartProcess(cmd, args, procAttr)
}
