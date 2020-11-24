package registry

import (
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io/ioutil"
	"log"
	"testing"
)

func TestUpload(t *testing.T) {
	name := core.ParseName("localhost:8082/gatblau/artie:v10")
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
	l.Remove([]*core.ArtieName{core.ParseName("874484")})
}

func TestTag(t *testing.T) {
	l := NewLocalRegistry()
	// l.Tag(core.ParseName("boot"), core.ParseName("gatblau/boot:v1"))
	l.Tag(core.ParseName("localhost:8081/gatblau/boot"), core.ParseName("boot:11"))
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
