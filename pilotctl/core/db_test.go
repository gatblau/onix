package core

import (
	"testing"
)

func TestToHStoreString(t *testing.T) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	s := toHStoreString(m)
	s1 := "\"key1\"=>\"value1\", \"key2\"=>\"value2\", \"key3\"=>\"value3\""
	if s != s1 {
		t.FailNow()
	}
}
