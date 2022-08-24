/*
  Onix Config Manager - OpenTelemetry collector for managed hosts
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

// Collector interface for host metrics
type Collector interface {
	Run(context.Context) error
	Stop()
	Restart(context.Context) error
	Status() <-chan *Status
}

// NewCollector returns a new collector.
func NewCollector(configPaths []string, version string, logOpts []zap.Option) Collector {
	return &otCollector{
		configPaths: configPaths,
		version:     version,
		logOpts:     logOpts,
		status:      make(chan *Status, 10), // buffered channel blocks after 10 messages
		wg:          new(sync.WaitGroup),
	}
}

// Status is the status of a collector.
type Status struct {
	Running bool
	Err     error
}

// otCollector the implementation of the Collector interface for OpenTelemetry.
type otCollector struct {
	configPaths []string
	version     string
	mux         sync.Mutex
	svc         *service.Collector
	logOpts     []zap.Option
	status      chan *Status
	wg          *sync.WaitGroup
}

func (c otCollector) Run(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.svc != nil {
		return errors.New("service already running")
	}
	settings, err := NewSettings(c.configPaths, c.version, c.logOpts)
	if err != nil {
		return err
	}
	// must create settings instance for every run
	svc, err := service.New(*settings)
	if err != nil {
		err = fmt.Errorf("failed to create service: %w", err)
		c.setStatus(false, err)
		return err
	}
	startupErr := make(chan error, 1)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	c.wg = wg
	c.svc = svc

	go func() {
		defer wg.Done()
		err = svc.Run(ctx)
		c.setStatus(false, err)
		if err != nil {
			startupErr <- err
		}
	}()
	// avoid race condition in OT collector if its shutdown channel is not initialised before the shutdown func is called
	return c.waitForStartBeforeShutdown(ctx, startupErr)
	return nil
}

func (c otCollector) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.svc == nil {
		return
	}
	c.svc.Shutdown()
	c.wg.Wait()
	c.svc = nil
}

func (c otCollector) Restart(ctx context.Context) error {
	c.Stop()
	return c.Run(ctx)
}

func (c otCollector) Status() <-chan *Status {
	return c.status
}

// setStatus will set the status of the collector
func (c *otCollector) setStatus(running bool, err error) {
	select {
	case c.status <- &Status{running, err}:
	default:
	}
}

// waitForStartBeforeShutdown waits for the service to startup before exiting.
func (c *otCollector) waitForStartBeforeShutdown(ctx context.Context, startupErr chan error) error {
	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()
	for {
		if c.svc.GetState() == service.Running {
			c.setStatus(true, nil)
			return nil
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			c.svc.Shutdown()
			return ctx.Err()
		case err := <-startupErr:
			return err
		}
	}
}
