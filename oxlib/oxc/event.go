package oxc

/*
   Onix Configuration Manager - HTTP Client
   Copyright (c) 2018-2021 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/

import (
	"crypto/tls"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTT client for change notifications
type EventManager struct {
	done   chan bool
	cfg    *EventConfig
	client MQTT.Client
}

// creates a new event manager subscribed to a specific topic
// cfg: the mqtt server configuration
func NewEventManager(cfg *EventConfig) (*EventManager, error) {
	// check the configuration is valid (preconditions)
	if ok, err := cfg.isValid(); !ok {
		return nil, err
	}
	m := new(EventManager)
	// create connection configuration
	connOpts := MQTT.NewClientOptions().AddBroker(cfg.Server).SetClientID(cfg.clientId()).SetCleanSession(true)
	// add credentials if provided
	if cfg.hasCredentials() {
		connOpts.SetUsername(cfg.Username)
		connOpts.SetPassword(cfg.Password)
	}
	// setup tls configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		ClientAuth:         cfg.ClientAuthType,
	}
	connOpts.SetTLSConfig(tlsConfig)
	// subscribe to the topic on connection
	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(cfg.topic(), byte(cfg.Qos), cfg.OnMsgReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	// finally create the client
	client := MQTT.NewClient(connOpts)
	// set up the manager
	m.client = client
	m.cfg = cfg
	// return a new setup manager
	return m, nil
}

// connect to the message broker
func (m *EventManager) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Connected to %s\n", m.cfg.Server)
	return nil
}

// disconnect from the message broker
func (m *EventManager) Disconnect(timeoutMilSecs uint) {
	m.client.Disconnect(timeoutMilSecs)
}
