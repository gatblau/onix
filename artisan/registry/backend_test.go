package registry

import "testing"

func TestNexusGetFile(t *testing.T) {
	nexus := NewNexus3Backend(
		"http://localhost:8081",
	)
	_, err := nexus.GetRepositoryInfo("gatblau", "boot", "admin", "admin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRemoteApiGetFile(t *testing.T) {
	api := NewGenericAPI(
		"localhost:8082",
		false,
	)
	_, err := api.GetRepositoryInfo("gatblau", "boot", "admin", "admin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
