package server

import "testing"

func TestK8S(t *testing.T) {
	k, err := NewK8S()
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	err = k.NewImagePipeline("dummy-app", "default")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
}
