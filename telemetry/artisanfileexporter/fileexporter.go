/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package artisanfileexporter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"encoding/binary"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/oxlib/resx"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	timeformat  = "2006_01_02_15_04_05_000000000"
	maxfilesize = int64(60) //60kb
	ext         = "inprocess"
	final       = "processed"
)

// Marshaler configuration used for marhsaling Protobuf to JSON.
var tracesMarshaler = ptrace.NewJSONMarshaler()
var metricsMarshaler = pmetric.NewJSONMarshaler()
var logsMarshaler = plog.NewJSONMarshaler()

// fileExporter is the implementation of file exporter that writes telemetry data to a file
// in Protobuf-JSON format.
type fileExporter struct {
	path  string
	mutex sync.Mutex
}

func (e *fileExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *fileExporter) ConsumeTraces(_ context.Context, td ptrace.Traces) error {
	buf, err := tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	return exportAsLine(e, buf, "traces")
}

func (e *fileExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	buf, err := metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return exportAsLine(e, buf, "metrics")
}

func (e *fileExporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	buf, err := logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return exportAsLine(e, buf, "logs")
}

func exportAsLine(e *fileExporter, buf []byte, exporttype string) error {
	// Ensure only one write operation happens at a time.
	e.mutex.Lock()
	defer e.mutex.Unlock()
	path := core.ToAbs(e.path)
	path = filepath.Join(path, exporttype)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			core.ErrorLogger.Printf("failed to create path %s, error %s \n", path, err)
		}
	}
	//check if there is already a file with extension .inprocess, if yes use it else create new
	files, err := filepath.Glob(filepath.Join(path, fmt.Sprintf("*.%s", ext)))
	if err != nil {
		core.ErrorLogger.Printf("failed to find inprocess file at path %s, error %s \n", path, err)
		return err
	}
	msize := int64(binary.Size(buf) / 1024)
	if len(files) == 0 {
		err = writeToNewFile(path, buf)
		return err
	} else {
		//get the size of .inprocess file
		f := files[0]
		core.InfoLogger.Printf("current inprocess file found, %s \n", f)
		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			core.ErrorLogger.Printf("failed to open inprocess file, %s , error %s \n", f, err)
			return err
		}
		defer file.Close()
		core.DebugLogger.Printf("finding the size of current inprocess file ")
		stat, err := file.Stat()
		if err != nil {
			core.ErrorLogger.Printf("failed to get stats for inprocess file, %s , error %s \n", f, err)
			return err
		}
		kb := (stat.Size() / 1024)
		core.DebugLogger.Printf("before writing to inprocess file %s, file size is [ %d ]kb and metrics data size is [ %d ]kb \n", f, kb, msize)
		total := (kb + msize)
		// after adding current data to existing inprocess file, if the size of in process file exceeds
		// the maxfilesize, then close the current inprocess file and delete the extension .inprocess
		// so it will be treated as completed and ready for upload, and the current data will be written
		//to new inprocess file
		if total > maxfilesize {
			core.DebugLogger.Printf("closing the current inprocess file ")
			err = file.Close()
			if err != nil {
				core.ErrorLogger.Printf("failed to close inprocess file, %s, error %s \n", f, err)
				return err
			}

			fnew := strings.Replace(f, ext, final, 1)
			err = os.Rename(f, fnew)
			if err != nil {
				core.ErrorLogger.Printf("failed to rename inprocess file, %s \n to new file name %s \n", f, fnew)
				core.ErrorLogger.Printf("error %s ", err)
				return err
			}
			core.DebugLogger.Printf("size of current inprocess file and metrics data size exceeds the max file size, so creating new file to write metrics data")
			err = writeToNewFile(path, buf)
			core.DebugLogger.Printf("metrics data wrote to new inprocess file ")
			return err
		} else {
			core.InfoLogger.Printf("writing metricss data to exitsing inprocess file, %s \n", f)
			newline := string("\n")
			if _, err := file.WriteString(newline); err != nil {
				core.ErrorLogger.Printf("failed to append new line to inprocess file, %s, error %s \n", f, err)
				return err
			}
			if _, err := file.Write(buf); err != nil {
				core.ErrorLogger.Printf("failed to write data to inprocess file, %s, error %s \n", f, err)
				return err
			}
			core.DebugLogger.Printf("size of current inprocess file and metrics data size is less than the max file size, so writing to the same file")
		}
	}

	return nil
}

func writeToNewFile(path string, buf []byte) error {
	core.InfoLogger.Printf("current inprocess file not found, so creating one \n")
	currentTime := time.Now().UTC()
	t := currentTime.Format(timeformat)
	filename := fmt.Sprintf("%s.%s", t, ext)
	path = filepath.Join(path, filename)
	err := resx.WriteFile(buf, path, "")
	if err != nil {
		core.ErrorLogger.Printf("failed to write to file, %s , error %s \n", path, err)
		return err
	}
	return err
}

func (e *fileExporter) Start(context.Context, component.Host) error {
	//var err error
	//e.file, err = os.OpenFile(e.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	//return err
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *fileExporter) Shutdown(context.Context) error {
	return nil
}
