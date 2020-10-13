/*
*    Onix Probare - Demo Application for reactive config management
*    Copyright (c) 2020 by www.gatblau.org
*
*    Licensed under the Apache License, Version 2.0 (the "License");
*    you may not use this file except in compliance with the License.
*    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*    Unless required by applicable law or agreed to in writing, software distributed under
*    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
*    either express or implied.
*    See the License for the specific language governing permissions and limitations under the License.
*
*    Contributors to this project, hereby assign copyright in this code to the project,
*    to be licensed under the same terms as the rest of the code.
 */
package core

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"
)

type MessageType int

const (
	// terminal log
	Terminal MessageType = iota
	// config file
	File
	// env vars
	Vars
	// control UI
	Control
)

// get all application specific environment variables
// as a JSON string
func getEnv() []string {
	var result = make([]string, 0)
	env := os.Environ()
	for _, envVar := range env {
		// filters out any variable that is not prefixed
		if strings.HasPrefix(envVar, "PROBARE_") {
			result = append(result, envVar)
		}
	}
	return result
}

// sends messages to the browser using web sockets
// ws: the websocket connection
// msg: the channel containing the messages to send
func send(conn *connection) {
Loop:
	for {
		select {
		// receive a new message to be sent to the browser
		case m, more := <-conn.msg:
			if more {
				// send message
				err := conn.ws.WriteMessage(websocket.TextMessage, marshal(m))
				if err != nil {
					log.Error().Msgf("cannot write message to websocket: %v, closing connection", err)
					// invalidate the connection
					conn.valid = false
					// close the connection
					conn.ws.Close()
					break Loop
				}
			} else {
				// closes the websocket
				log.Info().Msgf("closing WebSocket connection")
				conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
				conn.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				time.Sleep(closeGracePeriod)
				conn.ws.Close()
				break Loop
			}
		}
	}
	log.Info().Msg("message sender loop finishing")
}

func remove(s []*connection, i int) []*connection {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// calls the specified function and returns the Stdout of its execution
// f: the function to execute
func peekStdout(f func()) string {
	// lock this thread
	lock.Lock()
	// first assign the current Stdout writer to a temp variable
	stdoutWriter := os.Stdout
	// create a reader / writer linked pair
	r, w, _ := os.Pipe()
	// assign the writer to the Stdout
	os.Stdout = w
	// invoke the function for which the Stdout is needed
	f()
	// close the writer
	w.Close()
	// restore the original writer
	os.Stdout = stdoutWriter
	// unlock the thread
	lock.Unlock()
	// copy the content of the reader into a buffer
	var buf bytes.Buffer
	io.Copy(&buf, r)
	// write the peeked logs to the original Stdout
	os.Stdout.Write(buf.Bytes())
	// return the peeked logs
	return buf.String()
}

// send a WebSocket message to the client
func sendMsg(msgType MessageType, msgValue []string) {
	// create the message structure
	m := &message{
		Type: msgType,
		Body: msgValue,
	}
	// send the message to all the channels / websocket connections
	pool.send(m)

	log.Info().Msgf("sending websocket message of type %v to %v clients", msgType, pool.len())
}

// convert the message into a json []byte
func marshal(msg message) []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Error().Msgf("cannot marshal message: %v", err)
		return []byte{}
	}
	return bytes
}
