package util

import (
	"fmt"
	"testing"
)

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

func TestDbMan_GetReleaseInfo(t *testing.T) {
	_, _ = dbman.GetReleaseInfo("0.0.4")
}

func TestDbMan_RunQuery(t *testing.T) {
	// results, _, err := dbman.RunQuery("db-version", "0.0.4", []string{"0.0.4"})
	// if err != nil {
	// 	t.Error(err)
	// 	t.Fail()
	// }
	// if len(results.Rows) == 0 {
	// 	t.Error(err)
	// 	t.Fail()
	// }
	// csv, _ := dbman.TableTo(results, "csv")
	// print(csv)
}

func TestSaveConfig(t *testing.T) {
	dbman.SetConfig("Schema.URI", "AAAA")
	dbman.SaveConfig()
}

func TestUseConfig(t *testing.T) {
	dbman.UseConfigSet("", "myapp")
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

func TestDbMan_Create_Deploy(t *testing.T) {
	output, err, _ := dbman.Create()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	output, err, _ = dbman.Deploy()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
