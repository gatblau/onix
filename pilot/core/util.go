/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/google/renameio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// executes a command
func execute(cmd string) (string, error) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		//nolint:gosec
		c = exec.Command(strArr[0])
	} else {
		//nolint:gosec
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	// execute the command asynchronously
	if err := c.Start(); err != nil {
		return stderr.String(), fmt.Errorf("executing %s failed: %s", cmd, err)
	}
	done := make(chan error)
	// launch a go routine to wait for the command to execute
	go func() {
		// send a message to the done channel if completed or error
		done <- c.Wait()
	}()
	// wait for the done channel
	select {
	case <-done:
		// command completed
	case <-time.After(6 * time.Second):
		// command timed out after 6 secs
		return stderr.String(), fmt.Errorf("executing '%s' timed out", cmd)
	}
	return stdout.String(), nil
}

// creates a new Basic Authentication Token
func basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// copy a file from a source to a destination
func copyFile(src, dest string) error {
	srcContent, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcContent.Close()

	data, err := ioutil.ReadAll(srcContent)
	if err != nil {
		return err
	}
	return renameio.WriteFile(dest, data, 0644)
}

// compute an MD5 checksum for the specified string
func checksum(txt string) [16]byte {
	return md5.Sum([]byte(txt))
}

// cmdToRun := "/path/to/someCommand"
// args := []string{"someCommand", "arg1"}
// starts a process
func startProc(cmd string, args []string) (*os.Process, error) {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	return os.StartProcess(cmd, args, procAttr)
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
