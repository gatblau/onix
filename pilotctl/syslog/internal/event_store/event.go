package event_store

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Event struct {
	EventService Service
}

func (e *Event) SyslogServer(listener SyslogListener) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)

	server.SetHandler(handler)
	server.ListenUDP(listener.BindIP + ":" + listener.Port)
	server.ListenTCP(listener.BindIP + ":" + listener.Port)
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			event, err := Reformat(logParts)
			if err != nil {
				log.Printf("failed to reformat log due to error %s", err)
			}
			fmt.Println(event)
			eventID, err := e.EventService.Create(context.Background(), event)
			if err != nil {
				log.Printf("failed to save events due to error %s", err)
			}
			log.Printf("Successfully saved event with eventID: %s", eventID)
		}
	}(channel)
	server.Wait()
}

func Reformat(logPart interface{}) (elog EventLog, err error) {
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
	elog.EventID = "ENV-" + strconv.Itoa(randEventID)
	elog.Client = logRFC3164.Client
	elog.CreateTimeStamp = time.Now().Format(layout)
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
