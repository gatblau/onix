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
	"time"
)

// Pilot host
type Pilot struct {
	cfg  *Config
	info *HostInfo
	rem  *Rem
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
	r, err := NewRem()
	if err != nil {
		return nil, fmt.Errorf("cannot initialise remote control client: %s", err)
	}
	p := &Pilot{
		cfg:  cfg,
		info: info,
		rem:  r,
	}
	// return a new pilot
	return p, nil
}

func (p *Pilot) Start() {
	p.register()
	p.ping()
}

// register the host, keep retrying indefinitely until a registration is successful
func (p *Pilot) register() error {
	// checks if the host is already registered
	if !IsRegistered() {
		fmt.Println("attempting to registered host")
		// starts a loop
		for {
			err := p.rem.Register()
			// if no error then exit the loop
			if err == nil {
				fmt.Println("host registration successful")
				SetRegistered()
				break
			} else {
				fmt.Println(err)
			}
			// otherwise waits for a period before retrying
			fmt.Print(".")
			time.Sleep(1 * time.Minute)
		}
	}
	return nil
}

func (p *Pilot) ping() error {
	for {
		_, _ = p.rem.Ping()
		time.Sleep(15 * time.Second)
	}
	return nil
}
