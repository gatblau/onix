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
	"os"
	"os/user"
	"path"
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

// IsRegistered is the host registered?
func IsRegistered() bool {
	fi, err := os.Stat(regpath())
	return os.IsExist(err) || fi != nil
}

// SetRegistered set the host as registered
func SetRegistered() error {
	regFile, err := os.Create(regpath())
	if err != nil {
		return err
	}
	regFile.Close()
	return nil
}

func regpath() string {
	return path.Join(HomeDir(), ".pilot_reg")
}

// reverse the passed0in string
func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func newToken(hostId string) string {
	// create an authentication token as follows:
	// 1. takes machine id & unix time
	// 2. base 64 encode
	// 3. reverse string
	return reverse(
		base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s|%d", hostId, time.Now().Unix()))))
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
