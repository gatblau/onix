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
	return pilot, nil
}

func (p *Pilot) Host() {
	h, err := NewHostInfo()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	log.Info().Msgf(h.String())
}

func (p *Pilot) InitC() {
	p.fetchCfg()
}

func (p *Pilot) Sidecar() {
	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)

	log.Info().Msgf("pilot sidecar launching\n")
	// attempt to fetch the application configuration in the first place
	if p.fetchCfg() {
		p.reload()
	}
	// subscribe to configuration change notifications
	p.subscribe()

	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop
}

func (p *Pilot) onNotification(mqtt.Client, mqtt.Message) {
	// implement locking
	if p.fetchCfg() {
		p.reload()
	}
	// give it some time to reload
	time.Sleep(2 * time.Second)
	// check if app is ok after reload
	p.checkRestore()
}
