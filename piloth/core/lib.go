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
	"math"
	"os"
	"os/exec"
	"os/user"
	"strconv"
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

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// nextInterval calculates the next retry interval using exponential backoff strategy
// exponential backoff interval for registration retries
// waitInterval = base * multiplier ^ n
//   - base is the initial interval, ie, wait for the first retry
//   - n is the number of failures that have occurred
//   - multiplier is an arbitrary multiplier that can be replaced with any suitable value
func nextInterval(failureCount float64) time.Duration {
	// multiplier 2.0 yields 15s, 60s, 135s, 240s, 375s, 540s, etc
	interval := 15 * math.Pow(2.0, failureCount)
	// puts a maximum limit of 1 hour
	if interval > 3600 {
		interval = 3600
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%fs", interval))
	if err != nil {
		ErrorLogger.Printf(err.Error())
	}
	return duration
}

// collectorEnabled determine if the log collector should be enabled
// uses PILOT_LOG_COLLECTION, if its value is not set then the collector is enabled by default
// to disable the collector set PILOT_LOG_COLLECTION=false (possible values "0", "f", "F", "false", "FALSE", "False")
func collectorEnabled() (enabled bool) {
	var err error
	collection := os.Getenv("PILOT_LOG_COLLECTION")
	if len(collection) > 0 {
		enabled, err = strconv.ParseBool(collection)
		if err != nil {
			WarningLogger.Printf("invalid format for PILOT_LOG_COLLECTION variable: %s\n; log collection is enabled by default", err)
			enabled = true
		}
	} else {
		enabled = true
	}
	return enabled
}

func (p *Pilot) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		DebugLogger.Printf(msg, a...)
	}
}
