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
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// sidecar mode behaviour
type sidecar struct {
	// the mqtt broker client
	events *oxc.EventManager
	// config file watcher
	watcher *fsnotify.Watcher
	// the path to the configuration file
	cfgFile string
	// the natural key of the item to track
	itemKey string
	// the http client for the web api
	ox *oxc.Client
	// the config file MD5 integrity checksum
	checksum [16]byte
	// command to instruct the application to reload its configuration (if defined, it is exclusive with reloadURI)
	reloadCmd string
	// URI to post new configuration (if defined, it is exclusive with reloadCmd)
	reloadURI string
	// the content type passed to the http post reloadURI to describe the type of content submitted
	reloadURIContentType string
	// the username to authenticate with the reload URI endpoint
	reloadURIUsername string
	// the password to authenticate with the reload URI endpoint
	reloadURIPassword string
}

// create a new sidecar
func NewSidecar() (*sidecar, error) {
	// read configuration
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	s := &sidecar{
		cfgFile:              cfg.CfgFile,
		itemKey:              cfg.EmConf.ItemInstance,
		reloadCmd:            cfg.ReloadCmd,
		reloadURI:            cfg.ReloadURI,
		reloadURIContentType: cfg.ReloadURIContentType,
		reloadURIUsername:    cfg.ReloadURIUser,
		reloadURIPassword:    cfg.ReloadURIPwd,
	}
	// set the notification's handler
	cfg.EmConf.OnMsgReceived = s.onNotification
	// create the event manager
	em, err := oxc.NewEventManager(cfg.EmConf)
	if err != nil {
		return nil, err
	}
	s.events = em
	return s, nil
}

// start the sidecar
func (s *sidecar) Start() {
	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)
	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)
	log.Info().Msgf("pilot sidecar launching\n")
	// refresh the application configuration
	s.refresh()
	// subscribe to configuration change notifications
	s.subscribe()
	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop
	// close all pilot resources
	s.bye()
}

// fetch and save or post configuration
func (s *sidecar) refresh() {
	// attempt to fetch the application configuration in the first place
	if ok, cfg := s.fetch(); ok {
		// if a configuration file is defined
		if len(s.cfgFile) > 0 {
			// save the configuration to the file
			s.save(cfg)
		}
		// reload the configuration
		s.reload(cfg)
	}
}

// fetch configuration
func (s *sidecar) fetch() (bool, string) {
	// retrieve configuration information
	log.Info().Msgf("fetching configuration for application with key '%s'\n", s.itemKey)
	item, err := s.item()
	if err != nil {
		log.Warn().Msgf("cannot fetch application configuration: %s\n", err)
		log.Info().Msgf("the application configuration will be unmanaged until it is created in Onix")
	} else {
		log.Info().Msgf("application configuration retrieved successfully\n")
	}
	if item != nil {
		// compute the configuration file MD5 checksum
		s.checksum = checksum(item.Txt)
		return true, item.Txt
	}
	return false, ""
}

// get the item the sidecar is managing
func (s *sidecar) item() (*oxc.Item, error) {
	return s.ox.GetItem(&oxc.Item{Key: s.itemKey})
}

// save the passed in configuration to disk
func (s *sidecar) save(cfg string) error {
	log.Info().Msgf("backing up current configuration")
	err := copyFile(s.cfgFile, fmt.Sprintf("%s.bak", s.cfgFile))
	if err != nil {
		log.Warn().Msgf("cannot backup configuration: %s", err)
	}
	// write retrieved configuration to disk
	if len(cfg) > 0 {
		err = ioutil.WriteFile(s.cfgFile, []byte(cfg), 0644)
	} else {
		log.Warn().Msg("cannot write configuration to file, configuration is empty")
	}
	if err != nil {
		log.Error().Msgf("failed to write application configuration file: %s\n", err)
	} else {
		log.Info().Msgf("writing application configuration to '%s'\n", s.cfgFile)
	}
	return err
}

// instigate an application configuration reload
func (s *sidecar) reload(cfg string) {
	// if a reload command is defined
	if len(s.reloadCmd) > 0 {
		// execute the command
		execute(s.reloadCmd)
	} else
	// if a reload URI is defined
	if len(s.reloadURI) > 0 {
		// post the configuration to the URI
		s.postConfig(cfg)
	} else {
		// not reloading
		log.Info().Msg("skipping reloading")
	}
}

// post the app configuration to the reload URI
func (s *sidecar) postConfig(cfg string) {
	// gets a reader for the payload
	reader := bytes.NewReader([]byte(cfg))
	// constructs the request
	req, err := http.NewRequest("POST", s.reloadURI, reader)
	if err != nil {
		log.Error().Msgf("failed to create http request for reload URI: %s", err)
	}
	// if a content type is defined
	if len(s.reloadURIContentType) > 0 {
		req.Header.Add("Content-Type", s.reloadURIContentType)
	}
	// if a user name is provided then add basic authentication token and content type
	if len(s.reloadURIUsername) > 0 {
		req.Header.Add("Authorization", basicToken(s.reloadURIUsername, s.reloadURIPassword))
	}
	// sets a request timeout
	http.DefaultClient.Timeout = 6 * time.Second
	// issue the request
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("failed to post configuration to reload URI: %s", err)
	} else {
		log.Info().Msgf("application configuration successfully posted to '%s'", s.reloadURI)
	}
}

// connect to the MQTT broker and subscribe for notifications
func (s *sidecar) subscribe() {
	err := s.events.Connect()
	if err != nil {
		log.Error().Msgf("failed to connect to the notification broker: %s\n", err)
	} else {
		log.Info().Msgf("connected to notification broker, subscribed to '/II_%s' topic\n", s.itemKey)
	}
}

// refresh the application configuration when a change notification is received
func (s *sidecar) onNotification(mqtt.Client, mqtt.Message) {
	// refresh the configuration
	s.refresh()
	// give it some time to reload
	time.Sleep(2 * time.Second)
	// check if app is ok after reload
	s.checkRestore()
}

// if application is not ready after configuration reload then restore previous configuration
func (s *sidecar) checkRestore() {

}

// initialises a watcher for application configuration file changes
func (s *sidecar) newWatcher() {
	// set the config file watcher if a config file has been defined
	if len(s.cfgFile) > 0 {
		log.Info().Msgf("monitoring configuration file '%s' for unsolicited changes", s.cfgFile)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error().Msgf("cannot create a watcher for file '%s': %s", s.cfgFile, err)
		}
		// launch go routine to watch for file changes
		go s.monitorCfgFile()
		// add the file watcher
		err = watcher.Add(s.cfgFile)
		if err != nil {
			log.Error().Msgf("failed to add configuration file watcher: %s", err)
		}
		s.watcher = watcher
	}
}

// monitor the configuration file for changes
func (s *sidecar) monitorCfgFile() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			log.Warn().Msgf("configuration file event: '%s'", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				// check that the modified file checksum matches the original
				content, err := ioutil.ReadFile(s.cfgFile)
				if err != nil {
					log.Error().Msgf("cannot read modified configuration file: %s", err)
				}
				// if the files are different
				if s.checksum != checksum(string(content)) {
					log.Warn().Msgf("modified file has unauthorised content, proceeding to revoke any changes")
					s.refresh()
					log.Info().Msgf("configuration file changes successfully revoked")
				}
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Error().Msgf("file watcher error: %s", err)
		}
	}
}

// dispose the sidecar resources
func (s *sidecar) bye() {
	// if a file watcher has been defined
	if s.watcher != nil {
		// removes all watches and closes the events channel
		s.watcher.Close()
	}
}
