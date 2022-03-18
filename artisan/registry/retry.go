/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package registry

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"math/rand"
	"strings"
	"time"
)

func init() {
	// randomises the seed generator
	rand.Seed(time.Now().UnixNano())
}

// downloadFileRetry download the file specified by the passed-in info (either package seal or zip)
// using retry with exponential back-off
// note: downloadInfo parameter should be a pointer so that downloadedFilename can be passed back to the calling code
func downloadFileRetry(downloadInfo interface{}, attempts int) error {
	startingInterval := 5 * time.Second
	// retry attempts number of times applying exponential back-off intervals
	// adds jitter to the interval to prevent creating a Thundering Herd effect
	if err := retry(attempts, startingInterval, downloadFile, downloadInfo); err != nil {
		// if the retry failed returns the error
		return fmt.Errorf("cannot download package after %d attempts", attempts)
	}
	return nil
}

// downloadFile download the file specified by the passed-in info (either package seal or zip)
// it is meant to be used by the core.Retry() function within downloadFileRetry()
func downloadFile(info interface{}) error {
	i, ok := info.(*downloadInfo)
	if !ok {
		return fmt.Errorf("input parameter should be of type downloadInfo")
	}
	downFilename, err, status := i.api.Download(i.name.Group, i.name.Name, i.filename, i.uname, i.pwd, i.tls)
	// if the error is not recoverable
	if err != nil && status > 299 {
		// stop the retry
		return stop{err}
	}
	i.downloadedFilename = downFilename
	return err
}

// downloadInfo the information required by the file download function
// downloadedFilename: get populated with the location of the file downloaded
type downloadInfo struct {
	name               core.PackageName
	filename           string
	uname              string
	pwd                string
	tls                bool
	api                Api
	downloadedFilename string
}

func getRepositoryInfoRetry(repoInfo interface{}) error {
	attempts := 3
	startingInterval := 5 * time.Second
	if err := retry(attempts, startingInterval, getRepositoryInfo, repoInfo); err != nil {
		// if the retry failed returns the error
		return fmt.Errorf("cannot retrieve repository information after %d attempts: %s", attempts, err)
	}
	return nil
}

func getRepositoryInfo(info interface{}) error {
	i, ok := info.(*repositoryInfo)
	if !ok {
		return fmt.Errorf("input parameter should be of type repositoryInfo")
	}
	repo, err, code := i.api.GetRepositoryInfo(i.name.Group, i.name.Name, i.uname, i.pwd, i.tls)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP response to HTTPS client") {
			return &stop{err}
		}
	} else if code > 299 {
		return stop{err}
	}
	i.repo = repo
	return err
}

type repositoryInfo struct {
	name  core.PackageName
	uname string
	pwd   string
	tls   bool
	api   Api
	repo  *Repository
}

// stop is a wrapper for an error to tell the retry function that the retry loop should finish before
// it has reached the total number of retries
type stop struct {
	error
}

// retry the execution of a function if an error is returned to provide an easy way to increase resiliency.
// This is especially useful when making HTTP requests or doing anything else that has to reach out across the network.
// There are two options for stopping the retry loop before all the attempts are made:
//   1) return nil; or
//   2) return a wrapped error: stop{err}
// Choose option #2 when an error occurs where retrying would be futile.
// Consider most 4XX HTTP status codes. They indicate that the client has done something wrong and subsequent retries,
// without any modification to the request will result in the same response.
// In this case we still want to return an error, so we wrap the error in the stop type.
// The actual error that is returned by the retry function will be the original non-wrapped error.
// This allows for later checks like err == ErrUnauthorized.
func retry(attempts int, sleep time.Duration, f func(input interface{}) error, input interface{}) error {
	core.InfoLogger.Printf("retry in progress...\n")
	if err := f(input); err != nil {
		str := fmt.Sprintf("%T\n", err)
		if strings.Contains(str, "stop") {
			// Return the original error for later checking
			return err
		}
		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd:
			// The retry function recursively calls itself, counting down attempts and sleeping
			// for twice as long each time (i.e. exponential backoff).
			// This technique works well until the situation arises where a good number of clients start
			// their retry loops at roughly the same time.
			// This could happen if a lot of connections get dropped at once.
			// The retry attempts would then be in sync with each other, creating what is known as
			// the Thundering Herd problem. To prevent this, we can add some randomness by inserting the following
			// lines before we call time.Sleep:
			jitter := time.Duration(rand.Int63n(int64(sleep))) / 3

			// add the jitter to sleep period
			sleep = sleep + jitter

			go core.ErrorLogger.Printf("ERROR: %s; retrying in %.0f seconds\n", err, sleep.Seconds())

			// wait until next attempt
			time.Sleep(sleep)

			return retry(attempts, 2*sleep, f, input)
		}
		return err
	}
	return nil
}
