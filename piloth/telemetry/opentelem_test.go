/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package telemetry

import (
	"context"
	"testing"
)

func TestCollectorRun(t *testing.T) {
	collector := NewOpenTelemetry([]string{"./telem.yaml"}, "0.0.0", nil)
	err := collector.Run(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}
	status := <-collector.Status()
	if !status.Running {
		t.Fatal("service should be running")
	}
	if status.Err != nil {
		t.Fatal(status.Err.Error())
	}
	collector.Stop()
	status = <-collector.Status()
	if status.Running {
		t.Fatal("service should be stopped")
	}
}
