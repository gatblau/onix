package core

import "testing"

func Test(t *testing.T) {
	// LoadEnvFromFile("")
	name, err := ParseName("registry")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	println(name)
}
