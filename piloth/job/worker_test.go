package job

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/pilotctl/core"
	"log"
	"testing"
	"time"
)

// test the worker
func TestWorker(t *testing.T) {
	// create a new job processing worker
	w := NewWorker(
		// define the processing logic
		func(data interface{}) (string, error) {
			// unbox the data
			c, _ := data.(core.CmdInfo)
			// print start of job message
			log.Printf("processing job %d, %s -> %s\n", c.JobId, c.Package, c.Function)
			// simulate process with delay
			time.Sleep(1 * time.Second)
			// return the job result
			return fmt.Sprintf("JOB %d => xxxxxxx\n", c.JobId), nil
		})
	// start the worker loop
	w.Start()
	// add a couple of jobs
	w.AddJob(core.CmdInfo{
		JobId:    1111,
		Package:  "package_01",
		Function: "function_01",
		Input:    &data.Input{},
	})
	w.AddJob(core.CmdInfo{
		JobId:    2222,
		Package:  "package_02",
		Function: "function_02",
		Input:    &data.Input{},
	})
	// starts a loop to process job results
	var (
		results = 2
		count   = 0
	)
	for count < results {
		// attempt to retrieve the next result
		r, gotOne := w.Result()
		if gotOne {
			count++
			status := func() string {
				if r.Success {
					return "successful"
				}
				return "failed"
			}()
			log.Printf("result for job %d: %s\n", r.JobId, status)
		} else {
			fmt.Printf(".")
			time.Sleep(time.Second)
		}
	}
}
