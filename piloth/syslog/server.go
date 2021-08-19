package syslog

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/piloth/core"
	"gopkg.in/mcuadros/go-syslog.v2"
	"log"
)

// Server syslog log collection service that wraps a syslog server
type Server struct {
	server *syslog.Server
}

// NewServer creates an instance of a syslog collection service
func NewServer(bindIP, port string) (*Server, error) {
	// create local cache folder in pilot's current location
	core.CheckCachePath()
	channel := make(syslog.LogPartsChannel)
	sysServ := syslog.NewServer()
	sysServ.SetHandler(syslog.NewChannelHandler(channel))
	// uses RFC3164 because it is default for rsyslog
	sysServ.SetFormat(syslog.RFC3164)
	err := sysServ.ListenUDP(fmt.Sprintf("%s:%s", bindIP, port))
	if err != nil {
		return nil, err
	}
	go func(channel syslog.LogPartsChannel) {
		for logEntry := range channel {
			info, err := core.NewHostInfo()
			if err != nil {
				info = &core.HostInfo{}
			}
			event, err := NewEvent(logEntry, *info)
			if err != nil {
				log.Printf("cannot format syslog enrty: %s\n", err)
			}
			err = event.Save()
			if err != nil {
				log.Printf("cannot save syslog entry to file: %s\n", err)
			}
		}
	}(channel)
	return &Server{
		server: sysServ,
	}, nil
}

// Start the server
func (s *Server) Start() error {
	return s.server.Boot()
}

// Wait the server
func (s *Server) Wait() {
	s.server.Wait()
}

// Stop the server
func (s *Server) Stop() error {
	return s.server.Kill()
}
