package core

/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// HomeDir pilot's home directory
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// reverse the passed-in string
func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func newToken(hostUUID, hostIP, hostName string) string {
	// create an authentication token as follows:
	// 1. takes host uuid (i.e. machine Id + hostname hash), host ip, name and unix time
	// 2. base 64 encode
	// 3. reverse string
	return reverse(
		base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s|%s|%s|%d", hostUUID, hostIP, hostName, time.Now().Unix()))))
}

func readToken(token string) (string, bool, error) {
	// read token by:
	// 1. reverse string
	// 2. base 64 decode
	// 3. break down into parts
	value, err := base64.StdEncoding.DecodeString(reverse(token))
	if err != nil {
		return "", false, err
	}
	str := string(value)
	parts := strings.Split(str, "|")
	tokenTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", false, err
	}
	timeOk := (time.Now().Unix() - tokenTime) < (5 * 60)
	return parts[0], timeOk, nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
