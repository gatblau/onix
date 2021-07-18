package core

import "testing"

func TestMerge(t *testing.T) {
	e := &Envar{map[string]string{}}
	e.Vars["PORT_VALUE_1"] = "80"
	e.Vars["PORT_VALUE_2"] = "8080"
	e.Vars["PORT_VALUE_3"] = "443"
	m, _ := NewTemplMerger()
	m.LoadTemplates([]string{"test/test.yaml.t"})
	err := m.Merge(e)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
