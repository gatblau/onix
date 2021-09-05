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
	"github.com/gatblau/onix/piloth/job"
	"log/syslog"
	"os"
	"strings"
	"time"
)

// Pilot host
type Pilot struct {
	cfg       *Config
	info      *HostInfo
	ctl       *PilotCtl
	logs      *syslog.Writer
	worker    *job.Worker
	connected bool
}

func NewPilot() (*Pilot, error) {
	// read configuration
	cfg := &Config{}
	err := cfg.Load()
	if err != nil {
		return nil, err
	}
	info, err := NewHostInfo()
	if err != nil {
		return nil, err
	}
	logsWriter, err := syslog.New(syslog.LOG_ALERT, "onix-pilot")
	if err != nil {
		return nil, err
	}
	// create a new job worker
	worker := job.NewCmdRequestWorker(logsWriter)
	// start the worker
	worker.Start()
	r, err := NewPilotCtl(worker)
	if err != nil {
		return nil, err
	}
	p := &Pilot{
		cfg:    cfg,
		info:   info,
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
	InfoLogger.Printf("launching...\n")
	collector, err := NewCollector("0.0.0.0", p.cfg.getSyslogPort())
	if err != nil {
		ErrorLogger.Printf("cannot create pilot syslog collector: %s\n", err)
		os.Exit(1)
	}
	collector.Start()
	if !commandExists("art") {
		ErrorLogger.Printf("cannot find artisan CLI, ensure it is installed before running pilot\n")
		os.Exit(127)
	}
	p.register()
	p.ping()
}

// register the host, keep retrying indefinitely until a registration is successful
func (p *Pilot) register() {
	// starts a loop
	for {
		op, err := p.ctl.Register()
		// if no error then exit the loop
		if err == nil {
			switch strings.ToUpper(op) {
			case "I":
				InfoLogger.Printf("new host registration created successfully\n")
			case "U":
				InfoLogger.Printf("host registration updated successfully\n")
			case "N":
				InfoLogger.Printf("host already registered\n")
			}
			p.connected = true
			// break the loop
			break
		}
		// assume connectivity is down
		p.connected = false
		// the registration call failed, need to retry
		ErrorLogger.Printf("registration failed: %s, waiting 60 secs before attempting registration again\n", err)
		time.Sleep(1 * time.Minute)
	}
}

func (p *Pilot) ping() {
	InfoLogger.Printf("starting ping loop\n")
	for {
		cmd, err := p.ctl.Ping()
		if err != nil {
			// write to the console output
			InfoLogger.Printf("ping failed: %s\n", err)
			p.connected = false
		} else {
			if !p.connected {
				InfoLogger.Printf("ping loop operational\n")
			}
			p.connected = true
			// verify the host identity and command value integrity using Pretty Good Privacy
			err = verify(cmd.Value, cmd.Signature)
			// if the verification fails, it is likely spoofing of pilotctl has happened
			if err != nil {
				WarningLogger.Printf("invalid host signature, cannot trust the pilot control service => %s\n", err)
			} else { // if the host can be trusted
				// do we have a command to process?
				if cmd.Value.JobId > 0 {
					// execute the job
					InfoLogger.Printf("starting execution of job #%v, package => '%s', fx => '%s'\n", cmd.Value.JobId, cmd.Value.Package, cmd.Value.Function)
					p.worker.AddJob(cmd.Value)
				}
			}
		}
		time.Sleep(15 * time.Second)
	}
}

func (p *Pilot) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		DebugLogger.Printf(msg, a...)
	}
}
