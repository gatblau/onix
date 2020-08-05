package core

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/dbman/plugin"
	"os/exec"
	"testing"
	"time"
)

func init() {
	dbm, err := NewDbMan()
	if err != nil {
		panic(err)
	}
	DM = dbm
}

func TestFetchReleasePlan(t *testing.T) {
	plan, err := DM.GetReleasePlan()
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
	_, _, _ = DM.GetReleaseInfo("0.0.4")
}

// func TestDbMan_RunQuery(t *testing.T) {
// 	newDb()
// 	dbman.Create()
// 	dbman.Deploy()
// 	_, manifest, _ := dbman.GetReleaseInfo("0.0.4")
// 	results, _, _, err := dbman.Query(manifest, manifest.GetQuery("version-history"), nil)
// 	if err != nil {
// 		t.Error(err)
// 		t.Fail()
// 	}
// 	if len(results.Rows) == 0 {
// 		t.Error(err)
// 		t.Fail()
// 	}
// 	yaml := results.Sprint("yaml")
// 	print(yaml)
// }

func TestSaveConfig(t *testing.T) {
	DM.SetConfig("Schema.URI", "AAAA")
	DM.SaveConfig()
}

func TestUseConfig(t *testing.T) {
	DM.UseConfigSet("", "myapp")
}

func TestDbMan_CheckConfigSet(t *testing.T) {
	err := DM.CheckConfigSet()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_Serve(t *testing.T) {
	DM.Serve()
}

func TestDbMan_Create_Deploy(t *testing.T) {
	output, err, _ := DM.Create()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	output, err, _ = DM.Deploy()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_Upgrade(t *testing.T) {
	newDb()
	DM.Cfg.Set("AppVersion", "0.0.1")
	_, err, _ := DM.Create()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	_, err, _ = DM.Deploy()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	DM.Cfg.Set("AppVersion", "0.0.4")
	output, err, _ := DM.Upgrade()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDbMan_QueryWithParam(t *testing.T) {
	params := make(map[string]string)
	params["svc"] = "etcd"
	_, _, _, err := DM.Query("svc-down-instance", params)
	if err != nil {
	}
}

func TestDbMan_MergeTable(t *testing.T) {
	table, _, _, err := DM.Query("version-history", nil)
	if err != nil {
	}
	theme := DM.getTheme("basic")
	writer := &bytes.Buffer{}
	err = table.AsHTML(writer, &plugin.HtmlTableVars{
		Title:       "query.Name",
		Description: "query.Description",
		QueryURI:    "query.uri",
		Style:       theme.Style,
		Header:      theme.Header,
		Footer:      theme.Footer,
	})
	fmt.Print(err)
}

func newDb() {
	exec.Command("docker", "rm", "oxdb", "-f").Run()
	exec.Command("docker", "run", "--name", "oxdb", "-itd", "-p", "5432:5432", "-e", "POSTGRESQL_ADMIN_PASSWORD=onix", "centos/postgresql-12-centos7").Run()
	time.Sleep(2 * time.Second)
}
