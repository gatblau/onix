package plugins

import (
	"fmt"
	"testing"
)

func TestNewRootCfg(t *testing.T) {
	c := NewCache()
	fmt.Print(c.filename())
}
