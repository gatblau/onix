/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog/log"
)

// pilot monitors and reloads application configuration
type pilot struct {
	// the natural key of the item to track
	itemKey string
	// the http client for the web api
	ox *oxc.Client
	// the parser for item content
	parser *parser
	// file manager
	fileManager *fileman
	// pilot mode
	mode opMode
	// the command to launch the application
	cmd string
	// the arguments to launch the application
	args []string
	// the application process
	proc *procMan
	// the mqtt broker client
	events *oxc.EventManager
	// the pilot base configuration
	cfg *Config
}

// NewPilot create a new pilot
func NewPilot(mode opMode, cmd string, args []string) (*pilot, error) {
	// read configuration
	cfg := &Config{}
	err := cfg.Load()
	if err != nil {
		return nil, err
	}
	// create an onix web api client
	ox, err := oxc.NewClient(cfg.OxConf)
	if err != nil {
		return nil, err
	}
	pilot := &pilot{
		itemKey:     cfg.EmConf.ItemInstance,
		ox:          ox,
		parser:      NewParser(),
		mode:        mode,
		cmd:         cmd,
		args:        args,
		cfg:         cfg,
		fileManager: NewFileManager(),
		proc:        NewProcessManager(),
	}
	// set the notification's handler to use the one in pilot
	cfg.EmConf.OnMsgReceived = pilot.onNotification
	// create the event manager
	em, err := oxc.NewEventManager(cfg.EmConf)
	if err != nil {
		return nil, err
	}
	// allocate the event manager instance to pilot
	pilot.events = em
	// return a new pilot
	return pilot, nil
}

// refreshFile ensures the passed-in configuration is written to the file system and
// monitored for unsolicited changes
func (p *pilot) refreshFile(appCf *appCfg) error {
	// first check the reload is permitted based on the trigger type and pilot mode
	if err, ok := p.canReload(appCf.reloadTrigger); !ok {
		return err
	}
	// pre-condition check for file configurations only
	if appCf.meta.typeVal() != TypeFile {
		return errors.New("can only accept a Type='file' to reload")
	}
	// if the file manager is not managing the configuration
	if !p.fileManager.isManaged(appCf) {
		// start managing the file configuration and monitoring the file system for unsolicited changes
		p.fileManager.add(appCf)
	} else {
		// if already managed then update content in case it has changed
		file := p.fileManager.get(appCf.meta.Path)
		if file == nil {
			return errors.New(fmt.Sprintf("file manager cannot retrieve file by path = %s", appCf.meta.Path))
		}
		// updates the file content
		file.content = []byte(appCf.config)
		// updates the front matter
		file.meta = appCf.meta
		// save the configuration to the file system
		file.save()
	}
	return nil
}

// query the configuration manager and retrieve all defined application configurations
func (p *pilot) retrieveConfigurations() ([]*appCfg, error) {
	// retrieve configuration information
	log.Info().Msgf("fetching configuration for application with key %s", p.itemKey)
	item, err := p.ox.GetItem(&oxc.Item{Key: p.itemKey})
	if err != nil {
		log.Warn().Msgf("cannot fetch application configuration: %s", err)
		log.Info().Msgf("the application configuration will be unmanaged until it is created in Onix")
	} else {
		log.Info().Msgf("application configuration retrieved successfully")
	}
	if item != nil {
		// parse the item content into one or more configurations and their metadata
		configs, err := p.parser.parse(item.Txt)
		// if the parser fails return the error
		if err != nil {
			return nil, err
		}
		// return the configuration matching the required path
		return configs, nil
	}
	return nil, errors.New(fmt.Sprintf("no information could be retrieved from configuration manager for item '%s", p.itemKey))
}

// retrieve the application that matches the specified configuration path
func (p *pilot) retrieveConfigurationByPath(path string) (*appCfg, error) {
	configs, err := p.retrieveConfigurations()
	if err != nil {
		return nil, err
	}
	// return the configuration matching the required path
	return p.findCfgByPath(path, configs), nil
}

// find the configuration that matches the passed in path (file or http destination)
func (p *pilot) findCfgByPath(path string, configs []*appCfg) *appCfg {
	for _, config := range configs {
		if config.meta.Path == path {
			return config
		}
	}
	return nil
}

// reload all configurations
func (p *pilot) reloadAll() error {
	// retrieve all defined application configurations
	configs, err := p.retrieveConfigurations()
	if err != nil {
		return err
	}
	// deploy and reload
	for _, config := range configs {
		switch config.confType {
		// if the configuration is a file
		case TypeFile:
			{
				// refresh the local file
				err = p.refreshFile(config)
				if err != nil {
					return err
				}
			}
		}
		// trigger the reloading mechanism
		err = p.reload(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// reload the application configuration
func (p *pilot) reload(cf *appCfg) error {
	switch cf.reloadTrigger {
	case TriggerRestart:
		{
			// TODO: issue a process restart signal
			return errors.New(fmt.Sprintf("process restart trigger not implemented"))
		}
	case TriggerGet:
		{
			logger.Info().Msgf("reloading configuration resource (%s) using HTTP GET", cf.meta.Path)
			_, err := p.http("GET", cf.meta.Uri, "", nil)
			return err
		}
	case TriggerPost:
		{
			logger.Info().Msgf("reloading configuration resource (%s) using HTTP POST", cf.meta.Path)
			err := p.submitConfiguration(cf, "POST")
			if err != nil {
				log.Error().Msgf(err.Error())
				return err
			}
		}
	case TriggerPut:
		{
			logger.Info().Msgf("reloading configuration resource (%s) using HTTP PUT", cf.meta.Path)
			return p.submitConfiguration(cf, "PUT")
		}
	case TriggerSignal:
		{
			// TODO: issue a POSIX signal
			return errors.New(fmt.Sprintf("signal trigger not implemented"))
		}
	}
	return errors.New(fmt.Sprintf("reload trigger %s not supported for file configuration: use either signal, get, put, post or restart", cf.reloadTrigger))
}

// launch the application
func (p *pilot) launch() error {
	// if p.mode != Controller {
	// 	return errors.New(fmt.Sprintf("pilot cannot launch application in '%s' mode", p.mode))
	// }
	// parts := strings.Split(p.cmd, " ")
	// if len(parts) == 0 {
	// 	return errors.New(fmt.Sprintf("pilot cannot launch application as no command has been specified"))
	// }
	// procAttr := new(os.ProcAttr)
	// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	// appProc, err := os.StartProcess(parts[0], parts[1:], procAttr)
	// if err != nil {
	// 	return err
	// }
	// p.proc = appProc
	return nil
}

// restart an application
func (p *pilot) restart() error {
	return nil
}

// check if the pilot current operation mode can reload the configuration using the specified trigger
func (p *pilot) canReload(t trigger) (error, bool) {
	switch p.mode {
	case Sidecar:
		{
			return nil, t == TriggerRestart || t == TriggerPut || t == TriggerGet || t == TriggerPost
		}
	case Controller:
		{
			return nil, t == TriggerRestart || t == TriggerPut || t == TriggerGet || t == TriggerPost || t == TriggerSignal
		}
	}
	return errors.New(fmt.Sprintf("pilot in %s mode cannot reload configuration using %s trigger", p.mode, t)), false
}

// creates a new Basic Authentication Token
func (p *pilot) basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// refresh the application configuration when a change notification is received
func (p *pilot) onNotification(mqtt.Client, mqtt.Message) {
	// refresh the configuration
	p.reloadAll()
}

// connect to the MQTT broker and subscribe for notifications
func (p *pilot) subscribe() {
	err := p.events.Connect()
	if err != nil {
		log.Error().Msgf("failed to connect to the notification broker: %s\n", err)
	} else {
		log.Info().Msgf("connected to notification broker, subscribed to '/II_%s' topic\n", p.itemKey)
	}
}

func (p *pilot) Start() {
	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)
	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)
	log.Info().Msgf("pilot launching in %s mode", p.mode)
	// refresh the application configuration
	err := p.reloadAll()
	if err != nil {
		logger.Error().Msgf("cannot reload configurations: %v", err)
		logger.Info().Msgf("exiting")
		os.Exit(-1)
	}
	// if in controller mode the launch the app
	if p.mode == Controller {

	}
	// subscribe to configuration change notifications
	p.subscribe()
	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop
	// close all pilot resources
	p.Stop()
}

func (p *pilot) Stop() {
	p.fileManager.stop()
}

// send a configuration to the application via HTTP
func (p *pilot) submitConfiguration(cf *appCfg, method string) error {
	if method != "POST" || method != "PUT" {
		return errors.New("configuration can only be posted or put to a resource URI")
	}
	headers := http.Header{}
	// if authentication credentials exists
	if len(cf.meta.User) > 0 && len(cf.meta.Pwd) > 0 {
		// add Authorization header (with basic authentication token)
		headers.Set("Authorization", basicToken(cf.meta.User, cf.meta.Pwd))
	}
	// add Content-Type header
	headers.Set("Content-Type", cf.meta.ContentType)
	// submits the configuration
	_, err := p.http(method, cf.meta.Uri, cf.config, headers)
	return err
}

// Make a generic HTTP request
func (p *pilot) http(method string, url string, payload string, headers http.Header) (*http.Response, error) {
	// creates the request
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// add the http headers to the request
	req.Header = headers

	// submits the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// do we have a nil response?
	if resp == nil {
		return resp, errors.New(fmt.Sprintf("error: response was empty for resource: %s, check the service is up and running", url))
	}

	// check for response status
	if resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		err = errors.New(fmt.Sprintf("error: response returned status='%s', body='%s'", resp.Status, body))
		if err != nil {
			return resp, err
		}
	}

	return resp, err
}
