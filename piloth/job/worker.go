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
	ctl "github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"log/syslog"
	"os"
	"path/filepath"
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
func NewCmdRequestWorker(logger *syslog.Writer) *Worker {
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
				// peek the next job to be processed
				job, err := peekJob()
				// if it can't peek the next job
				if err != nil {
					// write an error
					ErrorLogger.Printf("cannot peek next job to process: %s\n", err)
					// sleep before trying again
					time.Sleep(1 * time.Minute)
					// restart the loop
					continue
				}
				// if the worker is ready to process a job and there are jobs waiting to start
				if w.status == ready && job != nil {
					// set the worker as busy
					w.status = busy
					InfoLogger.Printf("starting job %d, %s -> %s", job.cmd.JobId, job.cmd.Package, job.cmd.Function)
					// dump env vars if in debug mode
					w.debug(job.cmd.PrintEnv())
					// execute the job
					out, err := w.run(*job.cmd)
					if err != nil {
						InfoLogger.Printf("job %d, %s -> %s failed: %s", job.cmd.JobId, job.cmd.Package, job.cmd.Function, mask(err.Error(), job.cmd.User, job.cmd.Pwd))
					} else {
						InfoLogger.Printf("job %d, %s -> %s succeeded", job.cmd.JobId, job.cmd.Package, job.cmd.Function)
					}
					// collect result
					var errorMsg string
					if err != nil {
						errorMsg = mask(err.Error(), job.cmd.User, job.cmd.Pwd)
					}
					result := &ctl.JobResult{
						JobId:   job.cmd.JobId,
						Success: err == nil,
						Log:     out,
						Err:     errorMsg,
						Time:    time.Now(),
					}
					// add the last result to the submit queue
					err = submitJobResult(*result)
					// if the job result could not be saved
					if err != nil {
						SyslogWriter.Err(fmt.Sprintf("cannot persist result for Job Id = %d: %s\n", job.cmd.JobId, err))
					}
					// remove job from the queue
					err = removeJob(*job)
					// if the job could not be removed
					if err != nil {
						SyslogWriter.Err(fmt.Sprintf("cannot remove Job Id = %d from local queue: %s\n", job.cmd.JobId, err))
					}
					w.status = ready
				}
			}
		}()
	} else {
		InfoLogger.Printf("worker has already started\n")
	}
}

// Stop stops the worker execution loop
func (w *Worker) Stop() {
	w.cancel()
	w.status = stopped
}

func (w *Worker) Jobs() int {
	dir := processDir("")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("cannot read directory: %s", err))
	}
	count := 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".job" {
			count = count + 1
		}
	}
	return count
}

// AddJob add a new job for processing to the worker
func (w *Worker) AddJob(job ctl.CmdInfo) {
	err := addJob(Job{cmd: &job})
	if err != nil {
		ErrorLogger.Printf("%s\n", err)
	}
}

// Result returns the next
func (w *Worker) Result() (*ctl.JobResult, error) {
	return peekJobResult()
}

func run(data interface{}) (string, error) {
	// unbox the data
	cmd, ok := data.(ctl.CmdInfo)
	if !ok {
		return "", fmt.Errorf("Runnable data is not of the correct type\n")
	}

	// create the command to run
	var (
		artCmd = "exe"
	)
	// get the variables in the host environment
	hostEnv := merge.NewEnVarFromSlice(os.Environ())
	// get the variables in the command
	cmdEnv := merge.NewEnVarFromSlice(cmd.Env())
	// if the execution is containerised
	if cmd.Containerised {
		// use the exec command instead
		artCmd = fmt.Sprintf("%sc", artCmd)
	} else {
		// if not containerised add PATH to execution environment
		hostEnv.Merge(cmdEnv)
		cmdEnv = hostEnv
	}
	// if running in verbose mode
	if cmd.Verbose {
		// add ARTISAN_DEBUG to execution environment
		cmdEnv.Vars["ARTISAN_DEBUG"] = "true"
	}
	// create the command statement to run
	cmdString := fmt.Sprintf("art %s -u %s:%s %s %s", artCmd, cmd.User, cmd.Pwd, cmd.Package, cmd.Function)
	// run and return
	return build.ExeAsync(cmdString, ".", cmdEnv, false)
}

func (w *Worker) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		DebugLogger.Printf(fmt.Sprintf("DEBUG: %s", msg), a...)
	}
}

func (w *Worker) RemoveResult(result *ctl.JobResult) error {
	return removeJobResult(*result)
}

func mask(value, user, pwd string) string {
	str := strings.Replace(value, user, "****", -1)
	str = strings.Replace(str, pwd, "xxxx", -1)
	return str
}
