/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gatblau/oxc"
	"net/http"
	"os/exec"
	"syscall"
)

func execute(name string, args []string, w http.ResponseWriter) error {
	command := exec.Command(name, args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Printf("failed creating command stdoutpipe: %s", err)
		return err
	}
	defer func() {
		_ = stdout.Close()
	}()
	stdoutReader := bufio.NewReader(stdout)

	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Printf("failed creating command stderrpipe: %s", err)
		return err
	}
	defer func() {
		_ = stderr.Close()
	}()
	stderrReader := bufio.NewReader(stderr)

	if err = command.Start(); err != nil {
		return err
	}

	go handleReader(stdoutReader, w)
	go handleReader(stderrReader, w)

	if err = command.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if _, ok = exitErr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("run command failed: '%s' - '%s'", name, err)
			}
		}
		return err
	}
	return nil
}

func handleReader(reader *bufio.Reader, w http.ResponseWriter) {
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		w.Write([]byte(str))
		fmt.Printf("! %s\n", []byte(str))
	}
}

func readFlow(key string) (flow []byte, err error) {
	cfg := NewConf()
	oxUri, err := cfg.getOnixWAPIURI()
	if err != nil {
		return nil, err
	}
	oxUser, err := cfg.getOnixWAPIUser()
	if err != nil {
		return nil, err
	}
	oxPwd, err := cfg.getOnixWAPIPwd()
	if err != nil {
		return nil, err
	}
	oxcfg := &oxc.ClientConf{
		BaseURI:            oxUri,
		Username:           oxUser,
		Password:           oxPwd,
		InsecureSkipVerify: true,
	}
	oxcfg.SetAuthMode("basic")
	ox, err := oxc.NewClient(oxcfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create Onix http client: %s", err)
	}
	item, err := ox.GetItem(&oxc.Item{Key: key})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve flow specification from Onix: %s", err)
	}
	flow, err = json.Marshal(item.Meta)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal flow specification: %s", err)
	}
	return flow, nil
}
