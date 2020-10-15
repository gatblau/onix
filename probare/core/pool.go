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

import "github.com/gorilla/websocket"

type connectionPool struct {
	connections []*connection
}

func NewConnectionPool() *connectionPool {
	return &connectionPool{
		connections: make([]*connection, 0),
	}
}

func (pool *connectionPool) len() int {
	return len(pool.connections)
}

func (pool *connectionPool) add(c *connection) {
	pool.connections = append(pool.connections, c)
}

func (pool *connectionPool) removeInvalid() {
	for {
		ix := -1
		for i, conn := range pool.connections {
			if !conn.valid {
				ix = i
				break
			}
		}
		if ix != -1 {
			pool.connections = remove(pool.connections, ix)
		} else {
			break
		}
	}
}

func (pool *connectionPool) send(m *message) {
	// loop through the connections in the pool and send the message
	for _, conn := range pool.connections {
		if conn.valid {
			conn.msg <- *m
		}
	}
	// remove any invalid connections after an attempt has been made to send messages
	pool.removeInvalid()
}

// a connection to a WbeSocket client
type connection struct {
	// the channel to send messages to the client
	msg chan message
	// the websocket connection to the client
	ws *websocket.Conn
	// is the connection valid?
	valid bool
}
