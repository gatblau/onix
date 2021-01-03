package registry

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"io/ioutil"
	"log"
	"testing"
)

func TestUpload(t *testing.T) {
	name := core.ParseName("localhost:8082/gatblau/boot")
	local := NewLocalRegistry()
	local.Push(name, "admin:admin", false)
}

func TestDownload(t *testing.T) {
	name := core.ParseName("localhost:8082/gatblau/artie")
	local := NewLocalRegistry()
	local.Pull(name, "admin:admin", false)
}

func TestRemove(t *testing.T) {
	l := NewLocalRegistry()
	l.Remove([]*core.ArtieName{core.ParseName("879")})
}

func TestTag(t *testing.T) {
	l := NewLocalRegistry()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Tag(core.ParseName("localhost:8081/gatblau/boot"), core.ParseName("boot:11"))
}

func TestOpen2(t *testing.T) {
	l := NewLocalRegistry()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Open(core.ParseName("gatblau/boot"), "admin:admin", false, "", "", true)
}

func TestOpen(t *testing.T) {
	back := NewNexus3Backend("http://localhost:8081")
	file, err := back.Download("gatblau", "artie", "161120190537714-38c2222fe7.json", "admin", "admin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	b, err := ioutil.ReadAll(file)
	fmt.Print(string(b))
}

func TestGetRepoInfo(t *testing.T) {
	back := NewNexus3Backend("http://localhost:8081")
	repo, err := back.GetRepositoryInfo("gatblau", "boot", "admin", "admin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	bytes, _ := json.Marshal(repo)
	fmt.Print(string(bytes))
}

func TestGetManifest(t *testing.T) {
	// validate the name
	artie := core.ParseName("localhost:8082/gatblau/boot")
	// create a local registry
	local := NewLocalRegistry()
	// find the artefact in the local registry
	a := local.FindArtefact(artie)
	if a == nil {
		core.RaiseErr("artefact not found")
	}
	// get the artefact manifest
	m := local.GetManifest(a)
	// marshal the manifest
	bytes, err := json.MarshalIndent(m, "", "   ")
	core.CheckErr(err, "cannot marshal manifest")
	fmt.Printf(string(bytes) + "\n")
}

func TestUnzip(t *testing.T) {
	err := unzip("../images/bin/output/art.zip", ".")
	core.CheckErr(err, "cannot marshal manifest")
}
