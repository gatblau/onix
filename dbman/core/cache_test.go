package core

import (
	"fmt"
	"testing"
)

func TestNewRootCfg(t *testing.T) {
	c := NewCache()
	fmt.Print(c.filename())
}
