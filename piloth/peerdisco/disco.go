package peerdisco

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"fmt"
	"github.com/gatblau/onix/pilotctl/types"
	"github.com/reugn/go-quartz/quartz"
	"time"
)

type Disco struct {
	scheduler *quartz.StdScheduler
	info      types.HostInfo
}

func NewDisco(info types.HostInfo) *Disco {
	return &Disco{
		scheduler: quartz.NewStdScheduler(),
		info:      info,
	}
}

func (d *Disco) Start() error {
	// creates a job to scan network searching for other pilot agents
	scanJob := NewScanJob(d.info)
	// start the scheduler
	d.scheduler.Start()
	// scheduler the scan job every 60 seconds
	err := d.scheduler.ScheduleJob(scanJob, quartz.NewSimpleTrigger(60*time.Second))
	if err != nil {
		return fmt.Errorf("cannot schedule image check job: %s\n", err)
	}
	return nil
}
