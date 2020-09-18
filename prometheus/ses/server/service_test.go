package server

import "testing"

func TestModelExists(t *testing.T) {
	ses, err := NewSeS()
	if err != nil {
		t.Fatal(err)
	}
	ses.Start()
}
