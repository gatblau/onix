/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/radovskyb/watcher"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type CVEExporter struct {
	delay       time.Duration
	submit      func(cveReportFile string, delay time.Duration, ctl *PilotCtl) error
	ctl         *PilotCtl
	w           *watcher.Watcher
	pathToWatch string
}

func NewCVEExporter(ctl *PilotCtl, pathToWatch string) *CVEExporter {
	pathToWatch, _ = filepath.Abs(pathToWatch)
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile("^*.json$"), false))
	return &CVEExporter{
		submit:      postReport,
		ctl:         ctl,
		w:           w,
		pathToWatch: pathToWatch,
	}
}

func (r *CVEExporter) Start() error {
	core.InfoLogger.Printf("starting CVE exporter, listening for reports at %s\n", r.pathToWatch)
	go func() {
		for {
			select {
			case event := <-r.w.Event:
				// randomise the post over a 5-minute window to prevent all pilots hitting pilot-ctl at the same time
				err := r.submit(event.Path, time.Duration(int64(rand.Intn(5*60)))*time.Second, r.ctl)
				if err != nil {
					core.ErrorLogger.Printf("cannot submit CVE report: %s\n", err)
				}
			case err := <-r.w.Error:
				core.WarningLogger.Println(err.Error())
			case <-r.w.Closed:
				return
			}
		}
	}()
	if _, err := os.Stat(r.pathToWatch); os.IsNotExist(err) {
		if err = os.MkdirAll(r.pathToWatch, 0755); err != nil {
			core.RaiseErr("cannot create cve folder: %s", err)
		}
	}
	// watch this folder for changes
	if err := r.w.Add(r.pathToWatch); err != nil {
		log.Fatalln(err)
	}
	files, err := ioutil.ReadDir(r.pathToWatch)
	core.CheckErr(err, "cannot read CVE path")
	for _, file := range files {
		err = r.submit(file.Name(), time.Duration(0), r.ctl)
		if err != nil {
			core.ErrorLogger.Printf("cannot submit CVE report: %s\n", err)
		}
	}
	core.InfoLogger.Printf("watching for new CVE (*.json) reports at %s\n", r.pathToWatch)
	// Start the watching process - it'll check for changes every 15 secs.
	return r.w.Start(time.Second * 15)
}

func (r *CVEExporter) Close() {
	r.w.Close()
}

func postReport(cveReportFile string, delay time.Duration, ctl *PilotCtl) error {
	core.InfoLogger.Printf("new CVE report detected: %s\n", cveReportFile)
	core.InfoLogger.Printf("staggering publication by %v\n", delay)
	time.Sleep(delay)
	content, err := os.ReadFile(cveReportFile)
	if err != nil {
		return err
	}
	err = ctl.SubmitCveReport(content)
	if err != nil {
		return err
	}
	core.InfoLogger.Printf("CVE report %s posted successfully\n", cveReportFile)
	return nil
}
