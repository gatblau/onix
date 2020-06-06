package util

import "testing"

var dbman *DbMan

func init() {
	dbm, err := NewDbMan()
	if err != nil {
		panic(err)
	}
	dbman = dbm
}

func TestFetchReleasePlan(t *testing.T) {
	plan, err := dbman.GetReleasePlan()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(plan.Releases) == 0 {
		t.Errorf("no releases found in getPlan")
		t.Fail()
	}
}

func TestSaveConfig(t *testing.T) {
	dbman.SetConfig("Schema.URI", "AAAA")
	dbman.SaveConfig()
}

func TestUseConfig(t *testing.T) {
	dbman.UseConfigSet("", "myapp")
}

func TestDbMan_InitialiseDb(t *testing.T) {
	err := dbman.InitialiseDb()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_Deploy(t *testing.T) {
	err := dbman.Deploy("0.0.4")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_CheckConfigSet(t *testing.T) {
	err := dbman.CheckConfigSet()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_Serve(t *testing.T) {
	dbman.Serve()
}
