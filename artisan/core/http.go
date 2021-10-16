package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Get make a GET HTTP request to the specified URL
func Get(url, user, pwd string) (*http.Response, error) {
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// add http request headers
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("Authorization", BasicToken(user, pwd))
	}
	// issue http request
	resp, err := http.DefaultClient.Do(req)
	// do we have a nil response?
	if resp == nil {
		return resp, errors.New(fmt.Sprintf("error: response was empty for resource: %s", url))
	}
	// check error status codes
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("error: response returned status: %s. resource: %s", resp.Status, url))
	}
	return resp, err
}

func BasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

func Curl(uri string, method string, token string, validCodes []int, payload string, file string, maxAttempts int, delaySecs int, timeoutSecs int, headers []string, outputFile string) {
	var (
		bodyBytes []byte    = nil
		body      io.Reader = nil
		attempts            = 0
	)
	if len(payload) > 0 {
		if len(file) > 0 {
			RaiseErr("use either payload or file options, not both\n")
		}
		bodyBytes = []byte(payload)
	} else {
		if len(file) > 0 {
			abs, err := filepath.Abs(file)
			if err != nil {
				RaiseErr("cannot obtain absolute path for file using %s: %s\n", file, err)
			}
			bodyBytes, err = os.ReadFile(abs)
		}
	}
	if bodyBytes != nil {
		body = bytes.NewReader(bodyBytes)
	}
	// create request
	req, err := http.NewRequest(strings.ToUpper(method), uri, body)
	if err != nil {
		RaiseErr("cannot create http request object: %s\n", err)
	}
	// add authorization token to http request headers
	if len(token) > 0 {
		req.Header.Add("Authorization", token)
	}
	// add custom headers
	if headers != nil {
		for _, header := range headers {
			parts := strings.Split(header, ":")
			if len(parts) != 2 {
				WarningLogger.Printf("wrong format of http header '%s'; format should be 'key:value', skipping it\n", header)
				continue
			}
			req.Header.Add(parts[0], parts[1])
		}
	}
	// create http client with timeout
	client := &http.Client{
		Timeout: time.Duration(int64(timeoutSecs)) * time.Second,
	}
	// issue http request
	resp, err := client.Do(req)
	// retry if error or invalid response code
	for err != nil || !validResponse(resp.StatusCode, validCodes) {
		if err != nil {
			if resp != nil {
				ErrorLogger.Printf("unexpected error with response code '%s'; error was: '%s', retrying attempt %d of %d in %d seconds, please wait...\n", resp.StatusCode, err, attempts+1, maxAttempts, delaySecs)
			} else {
				ErrorLogger.Printf("unexpected error with no response; error was: '%s', retrying attempt %d of %d in %d seconds, please wait...\n", err, attempts+1, maxAttempts, delaySecs)
			}
		} else {
			ErrorLogger.Printf("invalid response code %d, retrying attempt %d of %d in %d seconds, please wait...\n", resp.StatusCode, attempts+1, maxAttempts, delaySecs)
		}
		// wait for next attempt
		time.Sleep(time.Duration(int64(delaySecs)) * time.Second)
		// issue http request
		resp, err = client.Do(req)
		// increments the number of attempts
		attempts++
		// exits if max attempts reached
		if attempts >= maxAttempts {
			RaiseErr("%s request to '%s' failed after %d attempts\n", strings.ToUpper(method), uri, maxAttempts)
		}
	}
	// if there is a response body prints it to stdout
	if resp != nil && resp.Body != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			WarningLogger.Printf("cannot print response body: %s\n", err)
		} else {
			if len(outputFile) > 0 {
				// saves the response to a file
				abs, err := filepath.Abs(outputFile)
				if err != nil {
					RaiseErr("cannot save response body to %s: %s\n", outputFile, err)
				}
				err = os.WriteFile(abs, b, 644)
				if err != nil {
					RaiseErr("cannot save response body to %s: %s\n", abs, err)
				}
			} else {
				// prints the response to sdt out
				fmt.Println(string(b[:]))
			}
		}
	}
}

func validResponse(responseCode int, validCodes []int) bool {
	for _, validCode := range validCodes {
		if responseCode == validCode {
			return true
		}
	}
	return false
}
