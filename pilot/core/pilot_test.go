package core

import "testing"

func TestApp(t *testing.T) {
	p, err := NewPilot()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	p.Sidecar()
}
