package core

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
	ctl "github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Job struct {
	file os.FileInfo
	cmd  *ctl.CmdInfo
}

// peekJob return the oldest job waiting to be processing without removing it from the queue
// fail-safe: peeking a job that has a start mark means the job started but for some reason no completion could be sent
//   to pilot control in this case, add a job result to the submit queue warning job might have not been completed
//   The start mark is removed when the job result has been submitted
//   If a submitted mark is found, the remove job is called and the next jon is peeked
func peekJob() (job *Job, err error) {
	var bytes []byte
	dir := processDir("")
	files, err := ls(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ".job" {
			bytes, err = ioutil.ReadFile(processDir(file.Name()))
			if err != nil {
				return nil, fmt.Errorf("cannot read job file %s: %s", file.Name(), err)
			}
			if len(bytes) == 0 {
				return nil, fmt.Errorf("job file %s is empty: %s", file.Name(), err)
			}
			var cmdInfo ctl.CmdInfo
			err = json.Unmarshal(bytes, &cmdInfo)
			if err != nil {
				// try and extract job Id from the file name with the following regex
				re := regexp.MustCompile(`\d[\d,]*`)
				// extract the job id
				jobIdString := re.FindString(file.Name())
				// parse string to int64
				jobId, err2 := strconv.ParseInt(jobIdString, 10, 64)
				// if we have an error we cannot get the Id
				if err2 != nil {
					// return a job with the file reference, so it can be dealt with by the caller
					return &Job{file: file}, fmt.Errorf(fmt.Sprintf("cannot unmarshal file '%s', possibly due to a corruption: %s, also failed to retrieve the job Id from file name: %s, and therefore cannot report back to pilot control; the file content was '%s'\n", file.Name(), err, err2, string(bytes)))
				}
				// if we managed to get a jobId, return a job with the Job Id and file reference
				return &Job{
					file: file,
					cmd:  &ctl.CmdInfo{JobId: jobId},
				}, fmt.Errorf("cannot unmarhsal file '%s', possibly due to a corruption: %s; the file content was '%s'\n", file.Name(), err, string(bytes))
			}
			job = &Job{
				file: file,
				cmd:  &cmdInfo,
			}
			if submittedMarkerExists(job.cmd.JobId) {
				// it means that the host halted after submitting job result but could not remove job from the queue
				// therefore removes job from the queue
				err = removeJob(job.cmd.JobId)
				if err != nil {
					return nil, err
				}
				// peek next job
				return peekJob()
			}
			// returns the found job and creates a started marker for the job in the file system
			return job, startedMarker(job)
		}
	}
	// no job found
	return nil, nil
}

// removeJob remove the specified job from the directory it is in
// failsafe: removes the submitted marker
func removeJob(jobId int64) error {
	dir := dataDir(fmt.Sprintf("job_%d.submitted", jobId))
	// remove submitted marker
	err := os.Remove(dir)
	if err != nil {
		return err
	}
	// remove job from queue
	dir = processDir(fmt.Sprintf("job_%d.job", jobId))
	return os.Remove(dir)
}

// addJob add a new job to the process queue
func addJob(job Job) error {
	bytes, err := json.Marshal(job.cmd)
	if err != nil {
		return err
	}
	dir := processDir(fmt.Sprintf("job_%d.job", job.cmd.JobId))
	return os.WriteFile(dir, bytes, os.ModePerm)
}

// ls files in a folder by date (oldest modified time first)
func ls(dirname string) ([]os.FileInfo, error) {
	// read files from folder
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("cannot read process directory for a list of jobs: %s", err)
	}
	// sort the file slice by ModTime()
	// ensuring older job is processed first
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().UnixNano() <= files[j].ModTime().UnixNano()
	})
	return files, nil
}

func startedMarker(job *Job) error {
	if job == nil {
		return nil
	}
	dir := dataDir(fmt.Sprintf("job_%d.started", job.cmd.JobId))
	return ioutil.WriteFile(dir, []byte{}, os.ModePerm)
}

func processDir(file string) string {
	fp := os.Getenv("PILOT_HOME")
	fp, _ = filepath.Abs(fp)
	if len(fp) == 0 {
		fp, _ = os.Executable()
	}
	return filepath.Join(fp, "data", "process", file)
}

func dataDir(file string) string {
	fp := os.Getenv("PILOT_HOME")
	fp, _ = filepath.Abs(fp)
	if len(fp) == 0 {
		fp, _ = os.Executable()
	}
	return filepath.Join(fp, "data", file)
}

func submitDir(file string) string {
	fp := os.Getenv("PILOT_HOME")
	fp, _ = filepath.Abs(fp)
	if len(fp) == 0 {
		fp, _ = os.Executable()
	}
	return filepath.Join(fp, "data", "submit", file)
}
