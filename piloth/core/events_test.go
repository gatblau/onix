package core

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"os"
	"testing"
)

func TestGetEvents(t *testing.T) {
	os.Setenv("PILOT_HOME", "../")
	events, err := getEvents(2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("processing %d events\n", len(events))
	err = removeEvents()
	if err != nil {
		t.Fatal(err)
	}
}
