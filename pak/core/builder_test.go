package core

import "testing"

func TestPack(t *testing.T) {
	p := NewBuilder()
	p.Build("https://github.com/gatblau/boot")
}
