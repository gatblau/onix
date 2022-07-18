/*
  Onix Config Manager - Warden
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package mode

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/gatblau/onix/warden/handler"

	"github.com/elazarl/goproxy"
)

// connectionStopListener return stoppableConn, it also tracks the lifetime of connection to notify
// when it is safe to terminate the application.
type connectionStopListener struct {
	net.Listener
	sync.WaitGroup
}

type stoppableConn struct {
	net.Conn
	wg *sync.WaitGroup
}

func newConnectionStopListener(l net.Listener) *connectionStopListener {
	return &connectionStopListener{l, sync.WaitGroup{}}
}

func (sl *connectionStopListener) Accept() (net.Conn, error) {
	c, err := sl.Listener.Accept()
	if err != nil {
		return c, err
	}
	sl.Add(1)
	return &stoppableConn{c, &sl.WaitGroup}, nil
}

func (sc *stoppableConn) Close() error {
	sc.wg.Done()
	return sc.Conn.Close()
}

func Tap(uri, credential, address string, verbose, insecureskipverify bool) {
	if len(uri) == 0 {
		panic("URI is mandatory when running warden in a tap mode")
	}

	if len(address) == 0 || !strings.Contains(address, ":") {
		panic("either port number is missing or invalid port number format, port number format must be :PORT_NUMBER")
	}
	// set env so that api conf can retrieve it from environment
	os.Setenv(string(handler.TapDomain), uri)
	os.Setenv(string(handler.TapInsecureSkipVerify), strconv.FormatBool(insecureskipverify))
	if len(credential) > 0 {
		creds := strings.Split(credential, ":")
		if len(creds) > 0 {
			os.Setenv(string(handler.TapUserName), creds[0])
			os.Setenv(string(handler.TapPassword), creds[1])
		} else {
			panic("Invalid credentials, valid format of credentials is username:password")
		}
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	rd := handler.NewRequestHandler()
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		rd.Handle(req, ctx)
		return req, nil
	})

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error listen:", err)
	}
	csl := newConnectionStopListener(l)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		csl.Add(1)
		csl.Close()
		csl.Done()
	}()

	log.Println("going to start proxy server")
	http.Serve(csl, proxy)
	csl.Wait()
	log.Println("exiting after closing connection")
}
