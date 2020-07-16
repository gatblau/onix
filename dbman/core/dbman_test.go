package core

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
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

func TestDbMan_Upgrade(t *testing.T) {
	newDb()
	dbman.Cfg.Set("AppVersion", "0.0.1")
	_, err, _ := dbman.Create()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	_, err, _ = dbman.Deploy()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	dbman.Cfg.Set("AppVersion", "0.0.4")
	output, err, _ := dbman.Upgrade()
	fmt.Print(output.String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func newDb() {
	exec.Command("docker", "rm", "oxdb", "-f").Run()
	exec.Command("docker", "run", "--name", "oxdb", "-itd", "-p", "5432:5432", "-e", "POSTGRESQL_ADMIN_PASSWORD=onix", "centos/postgresql-12-centos7").Run()
	time.Sleep(2 * time.Second)
}
