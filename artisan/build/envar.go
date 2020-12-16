package build

import (
	"fmt"
	"strings"
)

type envar struct {
	vars map[string]string
}

func NewEnVarFromMap(v map[string]string) *envar {
	return &envar{
		vars: v,
	}
}

func NewEnVarFromSlice(v []string) *envar {
	ev := &envar{
		vars: make(map[string]string),
	}
	for _, s := range v {
		kv := strings.Split(s, "=")
		ev.add(kv[0], kv[1])
	}
	return ev
}

func (e *envar) add(key, value string) {
	e.vars[key] = value
}

func (e *envar) slice() []string {
	var result []string
	for k, v := range e.vars {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

func (e *envar) append(v map[string]string) *envar {
	var result = make(map[string]string)
	result = e.vars
	for k, v := range v {
		result[k] = v
	}
	return NewEnVarFromMap(result)
}
