package types

/*
  Onix Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"fmt"
	"os"
)

type EventReceiver struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri"`
	// optional credentials if authentication is required
	User string `json:"user,omitempty"`
	Pwd  string `json:"pwd,omitempty"`
}

type EventReceivers struct {
	EventReceivers []EventReceiver `json:"event_receivers"`
}

func NewEventPubConf() *EventReceivers {
	confFile := receiverConfigFile()
	if len(confFile) > 0 {
		bytes, err := os.ReadFile(confFile)
		if err != nil {
			return nil
		}
		var conf EventReceivers
		err = json.Unmarshal(bytes, &conf)
		if err != nil {
			fmt.Printf("ERROR: cannot unmarshal event reciever configuration: %s; event receivers have been disabled\n", err)
			return nil
		}
		return &conf
	}
	return nil
}
