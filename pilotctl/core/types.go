package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"time"
)

// Cmd command information for remote host execution
type Cmd struct {
	// the natural key uniquely identifying the command
	Key string `json:"key"`
	// description of the command
	Description string `json:"description"`
	// the package to use
	Package string `json:"package"`
	// the function in the package to call
	Function string `json:"function"`
	// the function input information
	Input *data.Input `json:"input"`
	// the package registry user
	User string `json:"user"`
	// the package registry password
	Pwd string `json:"pwd"`
	// enables verbose output
	Verbose bool `json:"verbose"`
	// run command in runtime
	Containerised bool `json:"containerised"`
}

// NewPingResponse creates a new ping response
func NewPingResponse(cmdInfo CmdInfo, pingInterval time.Duration) (*PingResponse, error) {
	// create a signature for the envelope
	envelope := PingResponseEnvelope{
		Command:  cmdInfo,
		Interval: pingInterval,
	}
	signature, err := sign(envelope)
	if err != nil {
		return nil, fmt.Errorf("cannot sign ping response: %s", err)
	}
	return &PingResponse{
		Signature: signature,
		Envelope:  envelope,
	}, nil
}

// PingResponse a command for execution with a job reference
type PingResponse struct {
	// the envelope signature
	Signature string `json:"signature"`
	// the signed content sent to pilot
	Envelope PingResponseEnvelope `json:"envelope"`
}

// PingResponseEnvelope contains the signed content sent to pilot
type PingResponseEnvelope struct {
	// the information about the command to execute
	Command CmdInfo `json:"value"`
	// the ping interval
	Interval time.Duration `json:"interval"`
}

type CmdInfo struct {
	JobId         int64       `json:"job_id"`
	Package       string      `json:"package"`
	Function      string      `json:"function"`
	User          string      `json:"user"`
	Pwd           string      `json:"pwd"`
	Verbose       bool        `json:"verbose"`
	Containerised bool        `json:"containerised"`
	Input         *data.Input `json:"input,omitempty"`
}

func (c *CmdInfo) Value() string {
	var artCmd string
	// if command is to run in a runtime
	if c.Containerised {
		// use art exec
		artCmd = "exec"
	} else {
		// otherwise use art exe
		artCmd = "exe"
	}
	// if user credentials for the Artisan registry were provided
	if len(c.User) > 0 && len(c.Pwd) > 0 {
		// pass the credentials to the art cli
		return fmt.Sprintf("art %s -u %s:%s %s %s", artCmd, c.User, c.Pwd, c.Package, c.Function)
	}
	// otherwise run the command without credentials (assume public registry)
	return fmt.Sprintf("art %s %s %s", artCmd, c.Package, c.Function)
}

func (c *CmdInfo) Env() []string {
	var vars []string
	// append vars
	for _, v := range c.Input.Var {
		vars = append(vars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	// append secrets
	for _, s := range c.Input.Secret {
		vars = append(vars, fmt.Sprintf("%s=%s", s.Name, s.Value))
	}
	return vars
}

func (c *CmdInfo) Envar() *merge.Envar {
	return merge.NewEnVarFromSlice(c.Env())
}

func (c *CmdInfo) PrintEnv() string {
	var vars bytes.Buffer
	vars.WriteString("printing variables passed to the shell\n{\n")
	for _, v := range c.Input.Var {
		vars.WriteString(fmt.Sprintf("%s=%s\n", v.Name, v.Value))
	}
	vars.WriteString("}\n")
	return vars.String()
}

// Host monitoring information
type Host struct {
	HostUUID  string `json:"host_uuid"`
	OrgGroup  string `json:"org_group"`
	Org       string `json:"org"`
	Area      string `json:"area"`
	Location  string `json:"location"`
	Connected bool   `json:"connected"`
	LastSeen  int64  `json:"last_seen"`
	Since     int    `json:"since"`
	SinceType string `json:"since_type"`
}

// Registration information for host registration
type Registration struct {
	Hostname    string  `json:"hostname"`
	HostIP      string  `json:"host_ip"`
	MachineId   string  `json:"machine_id"`
	OS          string  `json:"os"`
	Platform    string  `json:"platform"`
	Virtual     bool    `json:"virtual"`
	TotalMemory float64 `json:"total_memory"`
	CPUs        int     `json:"cpus"`
}

// Reader Get a JSON bytes reader for the Serializable
func (r *Registration) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Bytes Get a []byte representing the Serializable
func (r *Registration) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}

// Event host events
type Event struct {
	// 0: host up, 1: host down, 2: network up, 3: network down
	Type int       `json:"type"`
	Time time.Time `json:"time"`
}

type Events []Event

// Job a job to be executed on one or more hosts
type Job struct {
	Id         int64    `json:"id"`
	HostUUID   string   `json:"host_uuid"`
	JobBatchId int64    `json:"job_batch_id"`
	FxKey      string   `json:"fx_key"`
	FxVersion  int64    `json:"fx_version"`
	Created    string   `json:"created"`
	Started    string   `json:"started"`
	Completed  string   `json:"completed"`
	Log        string   `json:"log"`
	Error      bool     `json:"error"`
	OrgGroup   string   `json:"org_group"`
	Org        string   `json:"org"`
	Area       string   `json:"area"`
	Location   string   `json:"location"`
	Tag        []string `json:"tag"`
}

// JobBatchInfo information required to create a new job batch
type JobBatchInfo struct {
	// the name of the batch (not unique, a user-friendly name)
	Name string `json:"name"`
	// a description for the batch (not mandatory)
	Description string `json:"description,omitempty"`
	// one or more search labels
	Label []string `json:"label,omitempty"`
	// the universally unique host identifier created by pilot
	HostUUID []string `json:"host_uuid"`
	// the unique key of the function to run
	FxKey string `json:"fx_key"`
	// the version of the function to run
	FxVersion int64 `json:"fx_version"`
}

type JobBatch struct {
	// the id of the job batch
	BatchId int64 `json:"batch_id"`
	// the name of the batch (not unique, a user-friendly name)
	Name string `json:"name"`
	// a description for the batch (not mandatory)
	Description string `json:"description,omitempty"`
	// creation time
	Created time.Time `json:"created"`
	// one or more search labels
	Label []string `json:"label,omitempty"`
	// owner
	Owner string `json:"owner"`
	// jobs
	Jobs int `json:"jobs"`
}

// ToJson convert the passed-in object to a JSON byte slice
// NOTE: json.Marshal is purposely not used as it will escape any < > characters
func ToJson(object interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// switch off the escaping!
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}

type Area struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Org struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Location struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Admission struct {
	HostUUID string   `json:"host_uuid"`
	OrgGroup string   `json:"org_group"`
	Org      string   `json:"org"`
	Area     string   `json:"area"`
	Location string   `json:"location"`
	Label    []string `json:"label"`
}

// Result
// note: ensure it is aligned with the same struct in piloth
type Result struct {
	// the unique job id
	JobId int64
	// indicates of the job was successful
	Success bool
	// the execution log for the job
	Log string
	// the error if any
	Err string
	// the completion time
	Time time.Time
}

func (r *Result) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

func (r *Result) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}

// PackageInfo describes a package and all its tags
type PackageInfo struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
	Ref  string   `json:"ref"`
}

// InitialConfig data returned to pilot upon registration
type InitialConfig struct {
	// the status of the registration - I: created, U: updated, N: already exist
	Operation string `json:"operation"`
}
