package syslog

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

type Server struct {
	bindIP string
	port   string
}

const timeFormat = "2006-01-02T15:04:05.0000"

func NewServer(bindIP, port string) *syslog.Server {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)

	server.SetHandler(handler)
	server.ListenUDP(fmt.Sprintf("%s:%s", bindIP, port))
	server.ListenTCP(fmt.Sprintf("%s:%s", bindIP, port))
	go func(channel syslog.LogPartsChannel) {
		for logEntry := range channel {
			event, err := format(logEntry)
			if err != nil {
				log.Printf("failed to format log enrty: %s\n", err)
			}
			err = save(event)
			if err != nil {
				log.Printf("cannot save syslog event: %s\n", "")
			}
		}
	}(channel)
	server.Wait()
	server.Boot()
	return server
}

func format(logPart interface{}) (elog EventLog, err error) {
	logRFC3164 := RsyslogLogRFC3164{}
	logByte, err := json.Marshal(logPart)
	if err != nil {
		return elog, fmt.Errorf("failed to convert interface to byte due to error %w", err)
	}

	if err = json.Unmarshal(logByte, &logRFC3164); err != nil {
		return elog, fmt.Errorf("failed to convert byte to struct due to error %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	randEventID := rand.Intn(99999) * 999999
	elog.EventID = strconv.Itoa(randEventID)
	elog.Client = logRFC3164.Client
	elog.CreateTimeStamp = time.Now().Format(timeFormat)
	elog.Hostname = logRFC3164.Hostname
	elog.HostID = ""
	elog.HostAddress = logRFC3164.Hostname
	elog.Location = ""
	elog.Facility = logRFC3164.Facility
	elog.Priority = logRFC3164.Priority
	elog.Severity = logRFC3164.Severity
	elog.Tag = logRFC3164.Tag
	elog.EventTimestamp = logRFC3164.Timestamp
	elog.Content = logRFC3164.Content
	elog.Details = ""
	return elog, nil
}

func save(ev EventLog) error {
	bytes, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename(ev.EventID), bytes, os.ModePerm)
}

func filename(id string) string {
	// TODO: save to specific logs folder
	return fmt.Sprintf("%s.json", id)
}
