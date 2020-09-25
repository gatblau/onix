/*
*    Onix ProtoTip - Demo Application for reactive config management
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
	// a channel to send messages to the client via web sockets
	msg = make(chan message)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type message struct {
	mtype int
	body  string
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
	quit := make(chan bool)

	// runs the HTTP server asynchronously
	go func() {
		log.Info().Msgf("Onix ProtoApp is listening on :%v", port)
		log.Info().Msgf("it took %v to start", time.Since(svr.start))
		if err := server.ListenAndServe(); err != nil {
			log.Info().Msgf("stopping the server: %v", err)
		}
	}()

	// loop to send WebSocket messages to the client
	go func(quitCh chan bool, msgCh chan message) {
	Loop:
		for {
			select {
			case <-quitCh:
				log.Info().Msg("closing message channel")
				close(msgCh)
				break Loop
			default:
				m := &message{
					mtype: 0,
					body:  "Ping!",
				}
				msgCh <- *m
				time.Sleep(3 * time.Second)
			}
		}
		// signal the routine is done
		log.Info().Msg("message generator routine finishing")
	}(quit, msg)

	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)

	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop
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
func serveWs(w http.ResponseWriter, r *http.Request) {
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
	// launch the subroutine to send WebSocket messages to the client
	go send(ws, msg)
}

// starts the server
func (svr *server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWs)
	// NOTE: add always as last handler!
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	svr.listen(r)
}

// sends messages to the browser using web sockets
// ws: the websocket connection
// msg: the channel containing the messages to send
func send(ws *websocket.Conn, msg <-chan message) {
Loop:
	for {
		select {
		// receive a new message to be sent to the browser
		case m, more := <-msg:
			if more {
				// send message
				err := ws.WriteMessage(websocket.TextMessage, []byte(m.body))
				if err != nil {
					log.Error().Msgf("cannot write message to websocket: %v, closing connection", err)
					ws.Close()
					break Loop
				}
			} else {
				// closes the websocket
				log.Info().Msgf("closing WebSocket connection")
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				time.Sleep(closeGracePeriod)
				ws.Close()
				break Loop
			}
		}
	}
	log.Info().Msg("message sender loop finishing")
}
