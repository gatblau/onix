/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

var P *Pilot

type Pilot struct {
	// configuration
	Cfg *Config
	// operation mode
	Mode PilotMode
	// onix client
	Ox *oxc.Client
	// event manager
	EM *oxc.EventManager
	// config file watcher
	W *fsnotify.Watcher
	// the last known config file MD5 checksum
	Checksum [16]byte
}

func NewPilot() (*Pilot, error) {
	pilot := new(Pilot)
	// read configuration
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	pilot.Cfg = cfg

	// create the Onix client
	client, err := oxc.NewClient(cfg.OxConf)
	if err != nil {
		return nil, err
	}
	pilot.Ox = client

	// set the notification's handler
	cfg.EmConf.OnMsgReceived = pilot.onNotification

	// create the event manager
	em, err := oxc.NewEventManager(cfg.EmConf)
	if err != nil {
		return nil, err
	}
	pilot.EM = em

	// initialises a configuration file watcher
	pilot.createWatcher()

	return pilot, nil
}

// launch pilot in host mode
func (p *Pilot) Host() {
	h, err := NewHostInfo()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	log.Info().Msgf(h.String())
}

// launches pilot in sidecar mode
func (p *Pilot) Sidecar() {
	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)
	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)
	log.Info().Msgf("pilot sidecar launching\n")
	// refresh the application configuration
	p.refreshCfg()
	// subscribe to configuration change notifications
	p.subscribe()
	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop
	// close all pilot resources
	p.bye()
}

// fetch and save or post configuration
func (p *Pilot) refreshCfg() {
	// attempt to fetch the application configuration in the first place
	if ok, cfg := p.fetch(); ok {
		// if a configuration file is defined
		if len(p.Cfg.CfgFile) > 0 {
			// save the configuration to the file
			p.save(cfg)
		}
		// reload the configuration
		p.reload(cfg)
	}
}

// refresh the application configuration when a change notification is received
func (p *Pilot) onNotification(mqtt.Client, mqtt.Message) {
	// refresh the configuration
	p.refreshCfg()
	// give it some time to reload
	time.Sleep(2 * time.Second)
	// check if app is ok after reload
	p.checkRestore()
}

// dispose pilot resources
func (p *Pilot) bye() {
	// if a file watcher has been defined
	if p.W != nil {
		// removes all watches and closes the events channel
		p.W.Close()
	}
}
