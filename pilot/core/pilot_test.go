package core

import "testing"

func TestPilot(t *testing.T) {
	p, err := NewPilot(Sidecar, "", nil)
	check(t, err)
	p.Start()
}
