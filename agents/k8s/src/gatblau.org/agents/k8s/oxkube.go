/*
   Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*
	oxkube is a simple http server process which exposes a webhook
 */
func main() {
	// determines the operation mode i.e. webhook or msg broker
	mode := getMode()

	switch mode {
	case "webhook":
		runWebhook()
	case "broker":
		panic("Mode 'broker' is not implemented.")
	default:
		panic(fmt.Sprintf("Mode '%s' is not implemented.", getMode()))
	}
}

// launch a webhook on a TCP port listening for events
func runWebhook() {
	// creates an http server listening on the specified TCP port
	server := &http.Server{Addr: fmt.Sprintf(":%s", getPort()), Handler: nil}

	// registers a handler for the /webhook endpoint
	http.HandleFunc("/webhook", webhook)

	// runs the server asynchronously
	go func() {
		log.Println(fmt.Sprintf("oxkube listening on :%s/webhook", getPort()))
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)

	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop

	// gets a context with some delay to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// releases resources if main completes before the delay period elapses
	defer cancel()

	// on error shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
		log.Println("shutting down oxkube")
	}
}

func getPort() string {
	return getEnv("OX_KUBE_PORT", "8080")
}

func getMode() string {
	return getEnv("OX_KUBE_MODE", "webhook")
}

func getEnv(key string, defaultValue string) string {
	value := ""
	if os.Getenv(key) != "" {
		value = os.Getenv(key)
	} else {
		value = defaultValue
	}
	return value
}

func webhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "GET":
		io.WriteString(w, "OxKube is ready to process events posted to this endpoint")
	case "POST":
		
	}
}
