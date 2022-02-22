/*
  Onix Config Manager - Http Server
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package httpserver

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gatblau/onix/oxlib/oxc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
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

// ParseBasicToken retrieve the username and password from the basic authentication token in the http request
func ParseBasicToken(r http.Request) (user, pwd string) {
	// get the token from the authorization header
	token := r.Header.Get("Authorization")
	return ReadBasicToken(token)
}

// ReadBasicToken retrieve the username and password from the passed-in basic authentication token
func ReadBasicToken(token string) (user, pwd string) {
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

// BasicToken creates a basic authentication token
func BasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
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

func IsErr(w http.ResponseWriter, err error, statusCode int, msg string) bool {
	if err != nil {
		msg = fmt.Sprintf("%s: %s\n", msg, err)
		log.Printf(msg)
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
		return true
	}
	return false
}

func Err(w http.ResponseWriter, statusCode int, msg string) {
	log.Printf("%s\n", msg)
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}

func Unmarshal(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot read input data: %s\n", err)
	}
	switch contentType {
	case "application/json":
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("cannot unmarshal input data: %s\n", err)
		}
	case "application/yaml":
		err = yaml.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("cannot unmarshal input data: %s\n", err)
		}
	case "application/xml":
		err = xml.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("cannot unmarshal input data: %s\n", err)
		}
	default:
		return fmt.Errorf("invalid Content-Type %s\n", contentType)
	}
	return nil
}

// FindRealIP find the real IP of the http requester
// uses X-Forwarded-For and X-Real-Ip http headers to discover the IP of the sender
// as otherwise, it is likely that the seen IP is the one of the load balancer that sits in front of the http service
func FindRealIP(r *http.Request) string {
	remoteIP := ""
	// the default is the originating ip, but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// otherwise, parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}
