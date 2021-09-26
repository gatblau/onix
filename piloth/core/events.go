package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	ctl "github.com/gatblau/onix/pilotctl/types"
	"os"
	"path/filepath"
)

// getEvents retrieve event log entries
func getEvents(max int) (*ctl.Events, error) {
	dir := submitDir("")
	files, err := ls(dir)
	if err != nil {
		return nil, err
	}
	// collect event file names up to a max number
	var (
		names  []string
		events = &ctl.Events{Events: []ctl.Event{}}
	)
	// loop through the files in submit directory
	for _, file := range files {
		// if the file is an event (*.ev)
		if !file.IsDir() && filepath.Ext(file.Name()) == ".ev" {
			// append its name to the event list
			names = append(names, file.Name())
			// read the event bytes
			bytes, err := os.ReadFile(submitDir(file.Name()))
			if err != nil {
				return nil, err
			}
			// unmarshal the event bytes
			var entry ctl.Event
			err = json.Unmarshal(bytes, &entry)
			if err != nil {
				return nil, err
			}
			// append the event to the event list
			events.Events = append(events.Events, entry)
			if len(names) >= max {
				break
			}
		}
	}
	// if there are no events
	if names == nil {
		// return
		return nil, nil
	}
	// otherwise, save the list of events being processed
	bytes, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}
	// write to file
	err = os.WriteFile(dataDir("events.json"), bytes, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// remove events that have been submitted
func removeEvents() error {
	// work out the file path where events
	dir := dataDir("events.json")
	bytes, err := os.ReadFile(dir)
	if err != nil {
		return nil
	}
	var names []string
	err = json.Unmarshal(bytes, &names)
	if err != nil {
		return nil
	}
	// remove the respective event files
	for i := 0; i < len(names); i++ {
		err = os.Remove(submitDir(names[i]))
		if err != nil {
			return err
		}
	}
	// remove the events.json file marker
	return os.Remove(dataDir("events.json"))
}
