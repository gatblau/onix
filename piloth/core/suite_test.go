package core

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"testing"
)

func TestNewRem(t *testing.T) {
	r, err := NewPilotCtl()
	if err != nil {
		t.FailNow()
	}
	err = r.Register()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestDmiDecode(t *testing.T) {
	id, err := machineid.ID()
	if err != nil {

	}
	fmt.Println(id)
}
