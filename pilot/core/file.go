/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"crypto/md5"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/google/renameio"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

// contains all information required by Pilot to manage a configuration file
type file struct {
	// the file metadata
	meta *frontMatter
	// configuration content
	content []byte
	// config file watcher
	watcher *fsnotify.Watcher
}

// create a new file data object and launch a file monitoring routine
func NewFile(cf *appCfg) *file {
	f := &file{
		meta:    cf.meta,
		content: []byte(cf.config),
	}
	log.Info().Msgf("monitoring configuration file '%s' for unsolicited changes", f.path())
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Msgf("cannot create a watcher for file '%s': %s", f.path(), err)
	}
	// start the monitor routine
	go f.monitor()
	// add the file watcher
	err = watcher.Add(f.path())
	if err != nil {
		log.Error().Msgf("failed to add configuration file watcher: %s", err)
	}
	// stores the watcher reference
	f.watcher = watcher
	// returns the file struct
	return f
}

// monitor the configuration file for changes
// note: launched from a go routine, to stop just close the watcher Events channel
func (f *file) monitor() {
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}
			log.Warn().Msgf("configuration file event: '%s'", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				// check that the modified file checksum matches the original
				// load the modified file
				newContent, err := ioutil.ReadFile(f.path())
				if err != nil {
					log.Error().Msgf("cannot read modified configuration file: %s", err)
				}
				// compare the checksums to see if they are different
				if md5.Sum(f.content) != md5.Sum(newContent) {
					log.Warn().Msgf("modified file has unauthorised content, proceeding to revoke any changes")
					// restore the content to the file
					f.save()
					// TODO: trigger reload of content?
					log.Info().Msgf("configuration file changes successfully revoked")
				}
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			log.Error().Msgf("file watcher error: %s", err)
		}
	}
}

// save the configuration file to disk
func (f *file) save() error {
	log.Info().Msgf("backing up current configuration")
	err := f.copy(fmt.Sprintf("%s.bak", f.path()))
	if err != nil {
		log.Warn().Msgf("cannot backup configuration: %s", err)
	}
	// write configuration to disk
	if len(f.content) > 0 {
		err = ioutil.WriteFile(f.path(), []byte(f.content), 0644)
	} else {
		log.Warn().Msg("cannot write configuration to file, configuration is empty")
	}
	if err != nil {
		log.Error().Msgf("failed to write application configuration file: %s\n", err)
	} else {
		log.Info().Msgf("writing application configuration to '%s'\n", f.path())
	}
	return err
}

// copy the file to a destination
func (f *file) copy(dest string) error {
	return renameio.WriteFile(dest, f.content, 0644)
}

// the path to the file
func (f *file) path() string {
	return f.meta.Path
}

// stop the monitoring process
func (f *file) stop() {
	close(f.watcher.Events)
}
