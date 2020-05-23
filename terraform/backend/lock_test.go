package main

import (
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/states/statemgr"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var (
	config *Config
	client *oxc.Client
)

func init() {
	// load the configuration file
	if cfg, err := NewConfig(); err == nil {
		config = cfg
	} else {
		panic(err)
	}
	// create a Web API client
	c, err := oxc.NewClient(config.Ox)
	if err != nil {
		panic(err)
	}
	client = c
	// check a model exists
	model := NewTerraModel(client)
	err = model.create()
	if err != nil {
		panic(err)
	}
}

func TestLockSuccess(t *testing.T) {
	// ensures no lock exist
	result, err := client.DeleteItem(lockKeyItem("foo"))
	if result.Error {
		t.Error(result.Message)
	}
	if err != nil {
		t.Error(err)
	}
	// try create a lock
	lock := NewLock("foo", client)
	id, _ := uuid.GenerateUUID()
	err, code := lock.lock(&statemgr.LockInfo{
		ID:        id,
		Operation: "O",
		Info:      "xxx",
		Who:       "a@a.com",
		Version:   "1.1.1",
		Created:   time.Now(),
		Path:      "a/b/c",
	})
	if err != nil {
		t.Error(err)
	}
	if code != http.StatusOK {
		t.Error(fmt.Sprintf("http status code was expected to be 200 OK but %s was found instead", strconv.Itoa(code)))
	}
}

func TestLockConflict(t *testing.T) {
	id1, _ := uuid.GenerateUUID()
	id2, _ := uuid.GenerateUUID()
	_, _ = client.DeleteItem(lockKeyItem("foo"))
	// acquires lock 1
	lock1 := NewLock("foo", client)
	_, _ = lock1.lock(&statemgr.LockInfo{
		ID:        id1,
		Operation: "O",
		Info:      "xxx",
		Who:       "a@a.com",
		Version:   "1.1.1",
		Created:   time.Now(),
		Path:      "a/b/c",
	})
	// tries and acquire another lock
	_, code := lock1.lock(&statemgr.LockInfo{
		ID:        id2,
		Operation: "O",
		Info:      "xxx",
		Who:       "b@b.com",
		Version:   "1.1.1",
		Created:   time.Now(),
		Path:      "a/b/c",
	})
	if code != http.StatusLocked {
		t.Error(fmt.Sprintf("http status code was expected to be 423 conflict but %s was found instead", strconv.Itoa(code)))
	}
}

func TestUnlockSuccess(t *testing.T) {
	id, _ := uuid.GenerateUUID()
	// clean any existing lock
	_, _ = client.DeleteItem(lockKeyItem("foo"))
	// acquires lock
	lock := NewLock("foo", client)
	_, _ = lock.lock(&statemgr.LockInfo{
		ID:        id,
		Operation: "O",
		Info:      "xxx",
		Who:       "a@a.com",
		Version:   "1.1.1",
		Created:   time.Now(),
		Path:      "a/b/c",
	})
	// now release the lock
	err, code := lock.unlock()
	if err != nil {
		t.Error(fmt.Sprintf("could not remove lock: %s", err))
	}
	if code != http.StatusOK {
		t.Error(fmt.Sprintf("http status code was expected to be 200 OK but %s was found instead", strconv.Itoa(code)))
	}
}
