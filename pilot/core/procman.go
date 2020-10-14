/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

// process manager used by pilot to manage the lifecycle of an applications
type procMan struct {
	pid          int
	process      *os.Process
	path         string
	cmd          string
	args         []string
	status       procStatus
	restartCount int
	startTime    int64
}

// start a process
func (proc *procMan) start(path string, cmd string, args []string) error {
	procAttr := &os.ProcAttr{
		Dir:   path,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	// create the arguments for the subprocess prepending the command name as the first argument
	logger.Info().Msgf("launching application %s", cmd)
	p, err := os.StartProcess(fmt.Sprintf("%s/%s", path, cmd), append([]string{cmd}, args...), procAttr)
	if err != nil {
		return err
	}
	proc.process = p
	proc.pid = proc.process.Pid
	proc.path = path
	proc.cmd = cmd
	proc.args = args
	proc.startTime = time.Now().Unix()
	proc.status = started
	logger.Info().Msgf("application %s launched successfully", cmd)
	return nil
}

// restart a process
func (proc *procMan) restart(timeOut time.Duration) error {
	// stops the application
	err := proc.stop(timeOut)
	if err != nil {
		return err
	}
	// starts the application
	err = proc.start(proc.path, proc.cmd, proc.args)
	if err != nil {
		return err
	}
	proc.restartCount++
	return nil
}

type procState struct {
	state *os.ProcessState
	err   error
}

// attempts to stop the process gracefully
func (proc *procMan) requestStop(timeOut time.Duration) (*os.ProcessState, error) {
	if proc.process != nil {
		logger.Info().Msgf("pilot is sending termination request signal")
		err := proc.process.Signal(syscall.SIGTERM)
		if err != nil {
			return nil, err
		}
		proc.status = stopRequested
		// new channel to handle process.Wait() response
		result := make(chan *procState)
		// launch a routine to wait for the process to finish
		go func(result chan *procState) {
			logger.Info().Msgf("pilot is waiting for the process to finish")
			// wait for the process to finish
			state, err := proc.process.Wait()
			// send the result back through the channel
			result <- &procState{
				state: state,
				err:   err,
			}
		}(result)
		select {
		// the process exited successfully
		case r := <-result:
			logger.Info().Msgf("application %s successfully terminated, state is %v", proc.cmd, r.state)
			proc.status = stopped
			// return the process state and / or any error
			return r.state, r.err
		// if the wait is longer than the specified timeOut
		case <-time.After(timeOut):
			logger.Info().Msgf(fmt.Sprintf("process did not terminate after %s, pilot will not wait any longer", timeOut))
			// wait no longer and return an error
			return nil, errors.New("process did not respond to termination request")
		}
	}
	return nil, errors.New("process does not exist")
}

// tries to stop the process gracefully but if it does not respond,
// then brutally kill it
func (proc *procMan) stop(timeOut time.Duration) error {
	// ask to stop the process politely ;)
	_, err := proc.requestStop(timeOut)
	// if it did not stop (have an error)
	if err != nil {
		// kill the process
		err := proc.kill()
		if err != nil {
			return err
		}
		proc.status = stopped
		proc.restartCount = 0
		proc.startTime = 0
		return nil
	}
	// if the programme exited successfully
	return nil
}

// terminate a process immediately
func (proc *procMan) kill() error {
	if proc.process != nil {
		logger.Info().Msgf("pilot is killing the process PID=%v", proc.pid)
		err := proc.process.Signal(syscall.SIGKILL)
		if err != nil {
			return err
		}
		err = proc.process.Release()
		if err != nil {
			return err
		}
		proc.startTime = 0
		proc.restartCount = 0
		proc.status = stopped
		logger.Info().Msgf("process has been killed successfully")
		return nil
	}
	return errors.New("process does not exist")
}

// true if the process is alive or false otherwise
// NOTE: this logic does not
func (proc *procMan) IsAlive() bool {
	// find the process
	p, err := os.FindProcess(proc.process.Pid)
	// if failed to find the process
	if err != nil {
		return false
	}
	// If  sig is 0, then no signal is sent, but error checking is still performed;
	// this can be used to check for the existence of a process ID or process group ID.
	return p.Signal(syscall.Signal(0)) == nil
}

// the the current uptime
func (proc *procMan) upTime() string {
	return proc.format(time.Now().Unix())
}

const (
	second = 1
	minute = second * 60
	hour   = minute * 60
	day    = hour * 24
	month  = day * 30
	year   = month * 12
)

// return a formatted representation of the gap between two times
func (proc *procMan) format(endTime int64) string {
	diff := endTime - proc.startTime
	if diff < 60 {
		return fmt.Sprintf("%ss", strconv.Itoa(int(diff)))
	} else if diff >= 60 && diff < hour {
		return fmt.Sprintf("%sm", strconv.Itoa(int(diff/minute)))
	} else if diff >= hour && diff < day {
		return fmt.Sprintf("%sh", strconv.Itoa(int(diff/hour)))
	} else if diff >= day && diff < month {
		return fmt.Sprintf("%sd", strconv.Itoa(int(diff/day)))
	} else if diff >= month && diff < year {
		return fmt.Sprintf("%sM", strconv.Itoa(int(diff/month)))
	}
	return fmt.Sprintf("%sy", strconv.Itoa(int(diff/year)))
}
