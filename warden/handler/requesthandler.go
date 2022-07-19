/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/elazarl/goproxy"
)

type ConfKey string

const (
	TapDomain             ConfKey = "TAP_DOMAIN"
	TapUserName           ConfKey = "TAP_USERNAME"
	TapPassword           ConfKey = "TAP_PWD"
	TapInsecureSkipVerify ConfKey = "TAP_INSECURE_SKIP_VERIFY"
	TapBearerToken        ConfKey = "TAP_BEARER_TOKEN"
)

type RequestDetails struct {
	req *http.Request
}

type RequestHandler struct {
	c chan *RequestDetails
}

func NewRequestHandler() *RequestHandler {
	rd := &RequestHandler{make(chan *RequestDetails)}
	go func() {
		for m := range rd.c {
			m.send()
		}
	}()
	return rd
}

func (m *RequestDetails) send() {

	if len(os.Getenv(string(TapBearerToken))) > 0 {
		bearer := "Bearer " + os.Getenv(string(TapBearerToken))
		m.req.Header.Add("Authorization", bearer)
	} else if len(os.Getenv(string(TapUserName))) > 0 && len(os.Getenv(string(TapPassword))) > 0 {
		m.req.SetBasicAuth(os.Getenv(string(TapUserName)), os.Getenv(string(TapPassword)))
	}

	resp, err := http.DefaultClient.Do(m.req)
	if err != nil {
		log.Println(" failed to send request ", err)
	} else {
		log.Println(" successfully sent request ", resp.Status)
	}
}

func (rd *RequestHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) {
	if req == nil {
		fmt.Println("empty http request")
	} else {
		//read body
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln("error while reading request body")
		}
		// copy body to original request as once request body is ready, the body is not available again
		// so adding back to original request
		req.Body = ioutil.NopCloser(bytes.NewReader(body))

		//build new url from original request
		url, err := url.Parse(req.URL.String())
		if err != nil {
			log.Fatalln("error while parsing original url")
		}

		//change the host of new url, this will be the host which is used to tap the request
		url.Host = os.Getenv(string(TapDomain))
		log.Println(" new URL is ", url.String())
		// build duplicate request using cloned body
		dupReq, err := http.NewRequest(req.Method, url.String(), ioutil.NopCloser(bytes.NewReader(body)))
		if err != nil {
			log.Fatalln("error while creating duplicate request")
		}
		// copy the headers from origin request to duplicate request
		dupReq.Header = req.Header
		r := &RequestDetails{req: dupReq}
		rd.addRequest(r)
	}
}

func (rd *RequestHandler) addRequest(m *RequestDetails) {
	rd.c <- m
}
