/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
    "encoding/base64"
    "fmt"
    "log"
    "math"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "time"
)

// HomeDir pilot's home directory
func HomeDir() string {
    defer TRA(CE())
    usr, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }
    return usr.HomeDir
}

// reverse the passed-in string
func reverse(str string) (result string) {
    defer TRA(CE())
    for _, v := range str {
        result = string(v) + result
    }
    return
}

func newToken(hostUUID, hostIP, hostName string) string {
    defer TRA(CE())
    // create an authentication token as follows:
    // 1. takes host uuid (i.e. machine Id + hostname hash), host ip, name and unix time
    // 2. base 64 encode
    // 3. reverse string
    return reverse(
        base64.StdEncoding.EncodeToString(
            []byte(fmt.Sprintf("%s|%s|%s|%d", hostUUID, hostIP, hostName, time.Now().Unix()))))
}

func commandExists(cmd string) bool {
    defer TRA(CE())
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
    defer TRA(CE())
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

func (p *Pilot) debug(msg string, a ...interface{}) {
    if len(os.Getenv("PILOT_DEBUG")) > 0 {
        DebugLogger.Printf(msg, a...)
    }
}

func Abs(path string) string {
    defer TRA(CE())
    if !filepath.IsAbs(path) {
        p, err := filepath.Abs(path)
        if err != nil {
            fmt.Printf("cannot work out absolute path for %s: %s\n", path, err)
            os.Exit(1)
        }
        path = p
    }
    return path
}

// Proxy returns a proxy func or nil if PILOT_HTTP_PROXY is not set
// Nil value means no proxy used in http.Transport
func Proxy() func(*http.Request) (*url.URL, error) {
    // constructs a proxy func if a proxy URL has been defined
    var proxy func(*http.Request) (*url.URL, error)
    // if a proxy url is set
    if v := os.Getenv("PILOT_HTTP_PROXY"); len(v) > 0 {
        proxyURL, err := url.Parse(v)
        if err != nil {
            log.Printf("invalid proxy URL: %s, proxy will be disabled\n", err)
        } else {
            // set the value of proxy function
            proxy = http.ProxyURL(proxyURL)
        }
    }
    return proxy
}
