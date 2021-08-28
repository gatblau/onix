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
	"log"
	"log/syslog"
	"os"
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
		logs:   logsWriter,
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
	log.Printf("launching...\n")
	collector, err := NewCollector("0.0.0.0", p.cfg.getSyslogPort())
	if err != nil {
		p.stdout("cannot create pilot syslog collector: %s\n", err)
		os.Exit(1)
	}
	collector.Start()
	if !commandExists("art") {
		p.stdout("cannot find artisan CLI, ensure it is installed before running pilot\n")
		os.Exit(127)
	}
	p.register()
	p.ping()
}

// warn: write a warning in syslog
func (p *Pilot) warn(msg string, args ...interface{}) {
	log.SetOutput(p.logs)
	p.logs.Warning(fmt.Sprintf(msg+"\n", args...))
}

// stdout: write a message to stdout
func (p *Pilot) stdout(msg string, args ...interface{}) {
	log.SetOutput(os.Stdout)
	log.Printf(msg+"\n", args...)
}

// register the host, keep retrying indefinitely until a registration is successful
func (p *Pilot) register() {
	// checks if the host is already registered
	if !IsRegistered() {
		p.stdout("host not registered, attempting registration")
		// starts a loop
		for {
			err := p.ctl.Register()
			// if no error then exit the loop
			if err == nil {
				p.stdout("registration successful")
				err = SetRegistered()
				p.connected = true
				if err != nil {
					p.stdout("failed to cache registration status: %s", err)
				}
				break
			} else {
				p.stdout("registration failed: %s", err)
				p.connected = false
			}
			// otherwise, waits for a period before retrying
			p.stdout("waiting 60s before attempting registration again")
			time.Sleep(1 * time.Minute)
		}
	} else {
		p.stdout("host is already registered")
	}
}

func (p *Pilot) ping() {
	p.stdout("starting ping loop")
	for {
		cmd, err := p.ctl.Ping()
		if err != nil {
			// write to the console output
			p.stdout("ping failed: %s", err)
			p.connected = false
		} else {
			if !p.connected {
				p.stdout("ping loop operational")
			}
			p.connected = true
			// verify the host identity and command value integrity using Pretty Good Privacy
			err = verify(cmd.Value, cmd.Signature)
			// if the verification fails, it is likely spoofing of pilotctl has happened
			if err != nil {
				p.stdout("invalid host signature, cannot trust the pilot control service => %s", err)
				// as it cannot trust the host writes a warning to syslog
				p.warn("invalid host signature, cannot trust the pilot control service => %s", err)
			} else { // if the host can be trusted
				// do we have a command to process?
				if cmd.Value.JobId > 0 {
					// execute the job
					p.stdout("starting execution of job #%v, package => '%s', fx => '%s'", cmd.Value.JobId, cmd.Value.Package, cmd.Value.Function)
					p.worker.AddJob(cmd.Value)
				}
			}
		}
		time.Sleep(15 * time.Second)
	}
}

func (p *Pilot) debug(msg string, a ...interface{}) {
	if len(os.Getenv("PILOT_DEBUG")) > 0 {
		p.stdout(fmt.Sprintf("DEBUG: %s", msg), a...)
	}
}
