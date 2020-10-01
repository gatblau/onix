/*
*    Onix ProtoApp - Demo Application for reactive config management
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
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var (
	port = 3000
	// a channel to indicate the web socket send channel should close
	done = make(chan bool)
	// a pool of websocket connections
	pool = make([]*connection, 0)

	// create a WebSocket connection by upgrading the original HTTP one
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// used to lock access to the Stdout whilst peeking
	lock sync.Mutex
)

// a WebSocket message
type message struct {
	Type int      `json:"type"`
	Body []string `json:"body"`
}

// a connection to a WbeSocket client
type connection struct {
	// the channel to send messages to the client
	msg chan message
	// the websocket connection to the client
	ws *websocket.Conn
}

type server struct {
	start time.Time
}

func NewServer() *server {
	return &server{
		start: time.Now(),
	}
}

func (svr *server) listen(handler http.Handler) {
	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)
	hangup := make(chan os.Signal, 1)

	// runs the HTTP server asynchronously
	go func() {
		log.Info().Msgf("Onix ProtoApp is listening on :%v", port)
		log.Info().Msgf("it took %v to start", time.Since(svr.start))
		if err := server.ListenAndServe(); err != nil {
			log.Info().Msgf("stopping the server: %v", err)
		}
	}()

	// sends any interrupt signal (SIGINT) to the stop channel
	signal.Notify(stop, os.Interrupt)
	// sends any hang-up signal (SIGHUP) to the hangup channel
	signal.Notify(hangup, syscall.SIGHUP)

Loop:
	for {
		select {
		case <-hangup:
			sendMsg(0, []string{"reloading configuration"})
			svr.LoadCfg()
		case <-stop:
			sendMsg(0, []string{"shutting down server"})
			time.Sleep(2 * time.Second)
			svr.Stop(server)
			break Loop
		}
	}
}

func (svr *server) Stop(server *http.Server) {
	log.Info().Msg("interrupt signal received")

	// gets a context with some delay to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// releases resources if main completes before the delay period elapses
	defer cancel()

	// on error shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("shutting down due to an error: %v", err)
	}
}

// http handler for the websocket connection
func (svr *server) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("opening WebSocket connection")
	// get a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	// if it failed
	if err != nil {
		// log the error and return
		log.Error().Msgf("failed to upgrade HTTP server connection to the WebSocket protocol: %v", err)
		return
	}
	log.Info().Msgf("launching message sender")

	// create a channel to send messages
	messageCh := make(chan message)

	// create the connection info reference
	conn := &connection{
		msg: messageCh,
		ws:  ws,
	}

	// add the connection info to the pool
	pool = append(pool, conn)

	// launch the subroutine to send WebSocket messages to the client
	// note: there are as many subroutines as web socket connections
	go send(ws, messageCh)

	// send configuration to the clients
	sendMsg(0, []string{"loading configuration"})
	svr.LoadCfg()
}

// starts the server
func (svr *server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", svr.serveWs)
	// NOTE: add always as last handler!
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	svr.listen(r)
}

func (svr *server) LoadCfg() {
	sendMsg(2, getEnv())
}
