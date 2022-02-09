/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"math/rand"
	"time"
)

// StopRetry is a wrapper for an error to tell the Retry function that the retry loop should finish before
// it has reached the total number of retries
type StopRetry struct {
	error
}

// Retry the execution of a function if an error is returned to provide an easy way to increase resiliency.
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
func Retry(attempts int, sleep time.Duration, f func(input interface{}) error, input interface{}) error {
	if err := f(input); err != nil {
		if s, ok := err.(StopRetry); ok {
			// Return the original error for later checking
			return s.error
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

			go ErrorLogger.Printf("ERROR: %s; retrying in %.0f seconds\n", err, sleep.Seconds())

			// wait until next attempt
			time.Sleep(sleep)

			return Retry(attempts, 2*sleep, f, input)
		}
		return err
	}

	return nil
}
