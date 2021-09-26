package core

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
	"github.com/gatblau/onix/pilotctl/types"
	"log"
	"os"
	"testing"
	"time"
)

// test the worker
func TestWorker(t *testing.T) {
	os.Setenv("PILOT_HOME", "../")
	// create a new job processing worker
	w := NewWorker(
		// define the processing logic
		func(data interface{}) (string, error) {
			// unbox the data
			c, ok := data.(types.CmdInfo)
			if !ok {
				panic("CmdInfo type casting failed")
			}
			// print start of job message
			log.Printf("processing job %d, %s -> %s\n", c.JobId, c.Package, c.Function)
			// simulate process with delay
			time.Sleep(1 * time.Second)
			// return the job result
			return fmt.Sprintf("JOB %d => complete\n", c.JobId), nil
		})
	// start the worker loop
	w.Start()
	// add a couple of jobs
	w.AddJob(types.CmdInfo{
		JobId:    1010,
		Package:  "list",
		Function: "list2",
		Input:    &data.Input{},
	})
	w.AddJob(types.CmdInfo{
		JobId:    1020,
		Package:  "list",
		Function: "list2",
		Input:    &data.Input{},
	})
	// wait until no more jobs to process
	for w.Jobs() > 0 {
		time.Sleep(1 * time.Second)
	}

	// starts a loop to process job results
	var (
		results = 2
		count   = 0
	)
	for count < results {
		// attempt to retrieve the next result
		r, _ := peekJobResult()
		if r != nil {
			count++
			status := func() string {
				if r.Success {
					return "successful"
				}
				return "failed"
			}()
			if status == "successful" {
				removeJobResult(*r)
			}
			log.Printf("result for job %d: %s\n", r.JobId, status)
		} else {
			fmt.Printf(".")
			time.Sleep(time.Second)
		}
	}
}
