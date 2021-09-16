package job

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"os"
	"path"
)

// submitJobResult persist the result of executing a Job in the file system
// fail-safe: remove start mark created by peek
//   create a submitted mark in case host goes before removing the job from the queue
func submitJobResult(result types.JobResult) error {
	dir := submitDir(fmt.Sprintf("job_%d.result", result.JobId))
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dir, bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return submittedMarker(result.JobId)
}

func submittedMarker(jobId int64) error {
	dir := dataDir(fmt.Sprintf("job_%d.submitted", jobId))
	// creates a submitted marker
	err := ioutil.WriteFile(dir, []byte{}, os.ModePerm)
	if err != nil {
		return err
	}
	// remove the started marker
	dir = dataDir(fmt.Sprintf("job_%d.started", jobId))
	return os.Remove(dir)
}

func submittedMarkerExists(jobId int64) bool {
	dir := dataDir(fmt.Sprintf("job_%d.submitted", jobId))
	_, err := os.Stat(dir)
	return err == nil
}

func peekJobResult() (jobResult *types.JobResult, err error) {
	var bytes []byte
	dir := submitDir("")
	files, err := ls(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ".result" {
			bytes, err = ioutil.ReadFile(submitDir(file.Name()))
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bytes, &jobResult)
			if err != nil {
				return nil, err
			}
			// returns the found job and creates a started marker for the job in the file system
			return jobResult, nil
		}
	}
	// no job found
	return nil, nil
}

func removeJobResult(result types.JobResult) error {
	// remove job from queue
	dir := submitDir(fmt.Sprintf("job_%d.result", result.JobId))
	return os.Remove(dir)
}
