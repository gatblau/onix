/*
  Onix Config Manager - Warden
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

// writes the content of an object using the response writer in the format specified by the accept http header
// supporting content negotiation for json, yaml, and xml formats
func Write(w http.ResponseWriter, r *http.Request, obj interface{}) {
	var (
		bs  []byte
		err error
	)
	// gets the accept http header
	accept := r.Header.Get("Accept")
	switch accept {
	case "*/*":
		fallthrough
	case "application/json":
		fallthrough
	default:
		w.Header().Set("Content-Type", "application/json")
		bs, err = json.Marshal(obj)
	case "application/yaml":
		w.Header().Set("Content-Type", "application/yaml")
		bs, err = yaml.Marshal(obj)
	case "application/xml":
		w.Header().Set("Content-Type", "application/xml")
		bs, err = xml.Marshal(obj)
	}
	if err != nil {
		WriteError(w, err, 500)
	}
	_, err = w.Write(bs)
	if err != nil {
		log.Printf("error writing data to response: %s\n", err)
		WriteError(w, err, 500)
	}
}

func WriteError(w http.ResponseWriter, err error, errorCode int) {
	fmt.Printf(fmt.Sprintf("%s\n", err))
	w.WriteHeader(errorCode)
}
