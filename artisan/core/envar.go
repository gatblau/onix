package core

import (
	"fmt"
	"strings"
)

type Envar struct {
	Vars map[string]string
}

func NewEnVarFromMap(v map[string]string) *Envar {
	return &Envar{
		Vars: v,
	}
}

func NewEnVarFromSlice(v []string) *Envar {
	ev := &Envar{
		Vars: make(map[string]string),
	}
	for _, s := range v {
		kv := strings.Split(s, "=")
		ev.Add(kv[0], kv[1])
	}
	return ev
}

func (e *Envar) Add(key, value string) {
	e.Vars[key] = value
}

func (e *Envar) Slice() []string {
	var result []string
	for k, v := range e.Vars {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

func (e *Envar) Append(v map[string]string) *Envar {
	var result = make(map[string]string)
	result = e.Vars
	for k, v := range v {
		result[k] = v
	}
	return NewEnVarFromMap(result)
}
