package server

import "testing"

func Test(t *testing.T) {
	_, err := getImgInfo("quay.io/gatblau/dummy-app", "", "")
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
}

func TestJob(t *testing.T) {
	// _ = parseTime("2020-07-21T12:15:56.18126Z")
	job, err := NewCheckImageJob()
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	job.Execute()
}
