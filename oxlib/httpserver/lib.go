package httpserver

/*
  Onix Config Manager - Http Client
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gatblau/onix/oxlib/oxc"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"strings"
)

// Write writes the content of an object using the response writer in the format specified by the accept http header
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

// ParseBasicToken getUser retrieve the username from the basic authentication token in the http request
func ParseBasicToken(r http.Request) (user, pwd string) {
	// get the token from the authorization header
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", ""
	}
	decoded, err := base64.StdEncoding.DecodeString(token[len("Basic "):])
	if err != nil {
		log.Printf("WARNING: failed to decode Authorization header: %s, cannot retrieve username\n", err)
		return "", ""
	}
	parts := strings.Split(string(decoded[:]), ":")
	if len(parts) != 2 {
		log.Printf("WARNING: failed to parse Authorization header: invalid format '%s' assuming a Basic Authentication Token\n", string(decoded[:]))
		return "", ""
	}
	// retrieve the username part (i.e. #0: username:password => 0:1)
	return parts[0], parts[1]
}

// GetUserPrincipal get the user principal in the http request
func GetUserPrincipal(r *http.Request) *oxc.UserPrincipal {
	// get the user from the request context
	var principal interface{}
	if principal = r.Context().Value("User"); principal != nil {
		if value, ok := principal.(*oxc.UserPrincipal); ok {
			return value
		}
	}
	return nil
}
