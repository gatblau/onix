package core

/*
  Onix Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"github.com/gatblau/onix/pilotctl/types"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConf(t *testing.T) {
	conf := &types.EventReceivers{EventReceivers: []types.EventReceiver{
		{
			URI:  "AAA",
			User: "BBB",
			Pwd:  "CCC",
		},
	}}
	bytes, _ := json.Marshal(conf)
	path, _ := filepath.Abs("ev_receive.json")
	os.WriteFile(path, bytes, os.ModePerm)
}
