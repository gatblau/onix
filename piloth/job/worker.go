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
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/pilotctl/core"
	"log"
	"log/syslog"
	"os"
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
	// syslog writer
	logs *syslog.Writer
}

// NewWorker create new worker using the specified runnable function
// Runnable: the function that processes each job
func NewWorker(run Runnable, logger *syslog.Writer) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		status:  stopped,
		jobs:    list.List{},
		results: list.List{},
		ctx:     ctx,
		cancel:  cancel,
		run:     run,
		logs:    logger,
	}
}

// NewCmdRequestWorker create a new worker to process pilotctl command requests
func NewCmdRequestWorker(logger *syslog.Writer) *Worker {
	return NewWorker(run, logger)
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
						w.stdout("invalid job format")
						continue
					}
					w.stdout("starting job %d, %s -> %s", cmd.JobId, cmd.Package, cmd.Function)
					// dump env vars if in debug mode
					w.debug(cmd.PrintEnv())
					// execute the job
					out, err := w.run(cmd)
					if err != nil {
						w.stdout("job %d, %s -> %s failed: %s", cmd.JobId, cmd.Package, cmd.Function, mask(err.Error(), cmd.User, cmd.Pwd))
					} else {
						w.stdout("job %d, %s -> %s succeeded", cmd.JobId, cmd.Package, cmd.Function)
					}
					// remove job from the queue
					w.jobs.Remove(jobElement)
					// collect result
					var errorMsg string
					if err != nil {
						errorMsg = mask(err.Error(), cmd.User, cmd.Pwd)
					}
					result := &Result{
						JobId:   cmd.JobId,
						Success: err == nil,
						Log:     out,
						Err:     errorMsg,
						Time:    time.Now(),
					}
					// add the last result to the list
					w.results.PushBack(result)
				}
			}
		}()
	} else {
		w.stdout("worker has already started")
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
		w.stdout("cannot unbox result")
	}
	// no result are available
	return nil, false
}

func run(data interface{}) (string, error) {
	// unbox the data
	cmd, ok := data.(core.CmdValue)
	if !ok {
		return "", fmt.Errorf("Runnable data is not of the correct type\n")
	}
	// create the command to run
	cmdString := fmt.Sprintf("art exe -u %s:%s %s %s", cmd.User, cmd.Pwd, cmd.Package, cmd.Function)
	// capture the PATH variable
	path := os.Getenv("PATH")
	env := cmd.Envar()
	// inject the PATH into the process
	env.Merge(merge.NewEnVarFromSlice([]string{fmt.Sprintf("PATH=%s", path)}))
	// run and return
	return build.Exe(cmdString, ".", env, false)
}

// warn: write a warning in syslog
func (w *Worker) warn(msg string, args ...interface{}) {
	log.SetOutput(w.logs)
	w.logs.Warning(fmt.Sprintf(msg+"\n", args...))
}

// stdout: write a message to stdout
func (w *Worker) stdout(msg string, args ...interface{}) {
	log.SetOutput(os.Stdout)
	log.Printf(msg+"\n", args...)
}

func (w *Worker) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		w.stdout(fmt.Sprintf("DEBUG: %s", msg), a...)
	}
}

func mask(value, user, pwd string) string {
	str := strings.Replace(value, user, "****", -1)
	str = strings.Replace(str, pwd, "xxxx", -1)
	return str
}
