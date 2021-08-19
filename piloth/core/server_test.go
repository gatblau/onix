package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"log/syslog"
	"os"
	"testing"
	"time"
)

func TestSyslogCollector(t *testing.T) {
	// set the config path to this folder
	os.Setenv("PILOT_CFG_PATH", ".")
	// create the syslog collector on port 534
	collector, err := NewCollector("0.0.0.0", "534")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// start listening
	err = collector.Start()
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// get a log writer
	logWriter, err := syslog.Dial("udp", "127.0.0.1:534", syslog.LOG_ERR, "Test")
	defer logWriter.Close()
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// write directly to the collector
	// note: in a real scenario, it should write to syslogs and let syslog forward to the collector
	err = logWriter.Err("error information here")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// wait for the collector to write the event to disk
	time.Sleep(time.Millisecond * 1500)
	// check event has been written
}

func TestNewLog(t *testing.T) {
	logWriter, err := syslog.Dial("udp", "127.0.0.1:534", syslog.LOG_ERR, "Test")
	defer logWriter.Close()
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	// write directly to the collector
	// note: in a real scenario, it should write to syslogs and let syslog forward to the collector
	err = logWriter.Err("error information here")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
}
