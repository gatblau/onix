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
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// how to use the event manager
func TestReceiver(t *testing.T) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	// create a new instance of the event manager
	m, err := NewEventManager(&EventConfig{
		Server:             "tcp://127.0.0.1:1883",
		ItemInstance:       "TEST_APP_01",
		Qos:                2,
		InsecureSkipVerify: true,
		OnMsgReceived:      onMsgReceived,
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	// connect and subscribe
	m.Connect()
	<-done
}

// a handler to process received messages
func onMsgReceived(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", msg.Topic(), msg.Payload())
}
