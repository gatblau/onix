package cmd

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	c "github.com/gatblau/onix/pilotctl/core"
	"github.com/mattn/go-shellwords"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Worker responsible for queue serving.
type Worker struct {
	Queue *Queue
}

// NewWorker initializes a new Worker.
func NewWorker(queue *Queue) *Worker {
	return &Worker{
		Queue: queue,
	}
}

// Start processes jobs from the queue (jobs channel).
func (w *Worker) Start(user, pwd string) bool {
	for {
		select {
		// if context was canceled.
		case <-w.Queue.ctx.Done():
			log.Printf("Work done in queue %s: %s!", w.Queue.name, w.Queue.ctx.Err())
			return true
		// if job received.
		case job := <-w.Queue.jobs:
			log, err := run(job, user, pwd)
			var result Result
			if err != nil {
				result = Result{
					Success: false,
					Log:     log,
					Err:     &err,
					Time:    time.Now(),
				}
			} else {
				result = Result{
					Success: true,
					Log:     log,
					Err:     nil,
					Time:    time.Now(),
				}
			}
			w.Queue.Results <- result
		}
	}
}

// executes a command and returns its output
func run(cmd c.CmdValue, user, pwd string) (string, error) {
	// create a command parser
	p := shellwords.NewParser()
	// parse the command line
	cmdArr, err := p.Parse(cmd.Value(user, pwd))
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
