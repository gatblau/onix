package server

import "testing"

func Test(t *testing.T) {
	_, err := getImgInfo("quay.io/gatblau/art-java", "", "")
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
}

func TestJob(t *testing.T) {
	job, err := NewCheckImageJob()
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	job.Execute()
}
