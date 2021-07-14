package job

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"container/list"
	"context"
	"fmt"
	"github.com/gatblau/onix/pilotctl/core"
	"github.com/mattn/go-shellwords"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// workerStatus the status of the job execution worker
type workerStatus int

const (
	// the worker loop is ready to process jobs
	ready workerStatus = iota
	// the worker is started and processing a specific job
	busy
	// the worker has not yet started
	stopped
)

// Runnable the function that carries out the job
type Runnable func(data interface{}) (string, error)

// Worker manage execution of jobs on a single job at a time basis
type Worker struct {
	// what is the current status of the worker?
	status workerStatus
	// the list of jobs to be processed
	jobs list.List
	// the results of the last job processed
	results list.List
	// the context to manage the worker loop go routine
	ctx context.Context
	// the function to cancel the worker loop go routine
	cancel context.CancelFunc
	// the logic that carries the instructions to process each job
	run Runnable
}

// NewWorker create new worker using the specified runnable function
// Runnable: the function that processes each job
func NewWorker(run Runnable) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		status:  stopped,
		jobs:    list.List{},
		results: list.List{},
		ctx:     ctx,
		cancel:  cancel,
		run:     run,
	}
}

// NewCmdRequestWorker create a new worker to process pilotctl command requests
func NewCmdRequestWorker() *Worker {
	return NewWorker(run)
}

// Start starts the worker execution loop
func (w *Worker) Start() {
	// if the worker is stopped then it can start
	if w.status == stopped {
		// changes the
		w.status = ready
		// launches the worker loop
		go func() {
			for {
				// if the worker is ready to process a job and there are jobs waiting to start
				if w.status == ready && w.jobs.Len() > 0 {
					// set the worker as busy
					w.status = busy
					// pick the next job from the queue
					jobElement := w.jobs.Front()
					// unbox the job
					cmd, ok := jobElement.Value.(core.CmdValue)
					if !ok {
						log.Printf("invalid job format")
						continue
					}
					log.Printf("starting job %d, %s -> %s\n", cmd.JobId, cmd.Package, cmd.Function)
					// execute the job
					out, err := w.run(cmd)
					// remove job from the queue
					w.jobs.Remove(jobElement)
					// collect result
					result := &Result{
						JobId:   cmd.JobId,
						Success: err == nil,
						Log:     out,
						Err:     &err,
						Time:    time.Now(),
					}
					// add the last result to the list
					w.results.PushBack(result)
				}
			}
		}()
	} else {
		log.Printf("worker has already started\n")
	}
}

// Stop stops the worker execution loop
func (w *Worker) Stop() {
	w.cancel()
	w.status = stopped
}

// AddJob add a new job for processing to the worker
func (w *Worker) AddJob(job core.CmdValue) {
	w.jobs.PushBack(job)
}

// Result get the next available result
func (w *Worker) Result() (*Result, bool) {
	e := w.results.Front()
	// if there is a result in the list
	if e != nil {
		r := e.Value.(*Result)
		if r != nil {
			// remove the result from the list after having read it
			w.results.Remove(e)
			// release the worker so that it can work on the next job
			w.status = ready
			// return the job status
			return r, true
		}
		log.Printf("cannot unbox result\n")
	}
	// no result are available
	return nil, false
}

// run executes a command and returns its output
func run(data interface{}) (string, error) {
	// unbox the data
	cmd, ok := data.(core.CmdValue)
	if !ok {
		return "", fmt.Errorf("Runnable data is not of the correct type\n")
	}
	// create a command parser
	p := shellwords.NewParser()
	// parse the command line
	cmdArr, err := p.Parse(cmd.Value())
	// if we are in windows
	if runtime.GOOS == "windows" {
		// prepend "cmd /C" to the command line
		cmdArr = append([]string{"cmd", "/C"}, cmdArr...)
	}
	name := cmdArr[0]
	var args []string
	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}
	command := exec.Command(name, args...)
	command.Env = cmd.Env()
	result, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(result), "\n"), nil
}
