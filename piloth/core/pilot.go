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
	ctl "github.com/gatblau/onix/pilotctl/types"
	"log/syslog"
	"math"
	"os"
	"strings"
	"time"
)

// Pilot host
type Pilot struct {
	cfg       *Config
	info      *ctl.HostInfo
	ctl       *PilotCtl
	logs      *syslog.Writer
	worker    *Worker
	connected bool
	// A duration string is a possibly signed sequence of
	// decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	pingInterval time.Duration
}

func NewPilot(hostInfo *ctl.HostInfo) (*Pilot, error) {
	// read configuration
	cfg := &Config{}
	err := cfg.Load()
	if err != nil {
		return nil, err
	}
	logsWriter, err := syslog.New(syslog.LOG_ALERT, "onix-pilot")
	if err != nil {
		return nil, err
	}
	// create a new job worker
	worker := NewCmdRequestWorker(logsWriter)
	// start the worker
	worker.Start()
	r, err := NewPilotCtl(worker, hostInfo)
	if err != nil {
		return nil, err
	}
	p := &Pilot{
		cfg:    cfg,
		info:   hostInfo,
		ctl:    r,
		worker: worker,
	}
	// return a new pilot
	return p, nil
}

func (p *Pilot) Start() {
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Lean&text=PILOT%0A
	fmt.Println(`+-----------------| ONIX CONFIG MANAGER |-----------------+
|      _/_/_/    _/_/_/  _/          _/_/    _/_/_/_/_/   |
|     _/    _/    _/    _/        _/    _/      _/        |
|    _/_/_/      _/    _/        _/    _/      _/         |
|   _/          _/    _/        _/    _/      _/          |
|  _/        _/_/_/  _/_/_/_/    _/_/        _/           |
|                     Host Controller                     | 
+---------------------------------------------------------+`)
	InfoLogger.Printf("launching pilot version %s\n", Version)
	InfoLogger.Printf("using Host UUID = '%s'\n", p.info.HostUUID)
	// creates a new SysLog collector
	collector, err := NewCollector("0.0.0.0", p.cfg.getSyslogPort())
	if err != nil {
		ErrorLogger.Printf("cannot create pilot syslog collector: %s\n", err)
		os.Exit(1)
	}
	// starts the collector service
	collector.Start()
	// check artisan cli is installed
	if !commandExists("art") {
		ErrorLogger.Printf("cannot find artisan CLI, ensure it is installed before running pilot\n")
		os.Exit(127)
	}
	// registers the host
	p.register()
	// initiates the ping loop
	p.ping()
}

// register the host, keep retrying indefinitely until a registration is successful
func (p *Pilot) register() {
	var failures float64 = 0
	// starts a loop
	for {
		op, err := p.ctl.Register()
		// if no error then exit the loop
		if err == nil {
			switch strings.ToUpper(op.Operation) {
			case "I":
				InfoLogger.Printf("new host registration created successfully\n")
			case "U":
				InfoLogger.Printf("host registration updated successfully\n")
			case "N":
				InfoLogger.Printf("host already registered\n")
			}

			// set the fallback ping interval, this will be automatically adjusted with the first ping response
			p.pingInterval, _ = time.ParseDuration("15s")

			// break the loop
			break
		}
		// calculates next retry interval
		interval := p.nextInterval(failures)

		// set the ping interval
		// the registration call failed, need to retry
		ErrorLogger.Printf("registration failed: %s, waiting %.2f minutes before attempting registration again\n", err, interval.Seconds()/60)

		// sleep until next ping
		time.Sleep(interval)

		// increment count
		failures = failures + 1
	}
}

// nextInterval calculates the next retry interval using exponential backoff strategy
// exponential backoff interval for registration retries
// waitInterval = base * multiplier ^ n
//   - base is the initial interval, ie, wait for the first retry
//   - n is the number of failures that have occurred
//   - multiplier is an arbitrary multiplier that can be replaced with any suitable value
func (p *Pilot) nextInterval(failureCount float64) time.Duration {
	// multiplier 2.0 yields 15s, 60s, 135s, 240s, 375s, 540s, etc
	interval := 15 * math.Pow(2.0, failureCount)
	// puts a maximum limit of 1 hour
	if interval > 3600 {
		interval = 3600
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%fs", interval))
	if err != nil {
		ErrorLogger.Printf(err.Error())
	}
	return duration
}

func (p *Pilot) ping() {
	for {
		resp, err := p.ctl.Ping()
		if err != nil {
			// write to the console output
			InfoLogger.Printf("ping failed: %s\n", err)
			p.connected = false
		} else {
			if !p.connected {
				InfoLogger.Printf("ping loop operational\n")
			}
			p.connected = true
			// verify the host identity and response integrity using Pretty Good Privacy (PGP
			err = verify(resp.Envelope, resp.Signature)
			// if the verification fails, it is likely spoofing of pilotctl has happened
			if err != nil {
				WarningLogger.Printf("invalid host signature, cannot trust the pilot control service => %s\n", err)
			} else { // if the host can be trusted
				cmd := resp.Envelope.Command
				// do we have a command to process?
				if cmd.JobId > 0 {
					// execute the job
					InfoLogger.Printf("starting execution of job #%v, package => '%s', fx => '%s'\n", cmd.JobId, cmd.Package, cmd.Function)
					p.worker.AddJob(cmd)
				}
			}
		}
		// if the  pilot interval is different from the interval requested by pilot control
		if resp.Envelope.Interval.Seconds() > 0 && p.pingInterval != resp.Envelope.Interval {
			// issue a notice about the ping interval adjustment
			InfoLogger.Printf("adjusting ping interval to %.0f seconds\n", resp.Envelope.Interval.Seconds())
			// update the local interval value
			p.pingInterval = resp.Envelope.Interval
		}
		// waits for the requested interval
		time.Sleep(p.pingInterval)
	}
}

func (p *Pilot) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		DebugLogger.Printf(msg, a...)
	}
}
