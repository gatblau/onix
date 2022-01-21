//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-Present by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.

package core

import "testing"

func TestNewScriptManager(t *testing.T) {
	// create an instance of the current configuration set
	cfg := NewConfig("", "")
	// create an instance of the script manager
	sm, err := NewScriptManager(cfg)
	if err != nil {
		t.Fatal(err)
	}
	info, manifest, err := sm.fetchManifest("0.0.4")
	if err != nil {
		t.Fatal(err)
	}
	if info == nil || manifest == nil {
		t.FailNow()
	}

	info, err = sm.getReleaseInfo("0.0.4")
	if err != nil {
		t.Fatal(err)
	}
	if info == nil {
		t.FailNow()
	}

	plan, err := sm.fetchPlan()
	if err != nil {
		t.Fatal(err)
	}
	if plan == nil {
		t.FailNow()
	}
}
