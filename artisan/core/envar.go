package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func NewEnVarFromFile(envFile string) (*Envar, error) {
	var outMap = make(map[string]string)
	file := ToAbs(envFile)
	data, err := ioutil.ReadFile(file)
	// if it managed to find the env file load it
	// otherwise skip it
	content := strings.Split(string(data), "\n")
	if err == nil {
		for _, line := range content {
			// skips comments
			if strings.HasPrefix(strings.Trim(line, " "), "#") || len(strings.Trim(line, " ")) == 0 {
				continue
			}
			keyValue := strings.Split(line, "=")
			if len(keyValue) != 2 {
				return nil, fmt.Errorf("invalid env file format")
			}
			outMap[keyValue[0]] = keyValue[1]
		}
	} else {
		Debug("cannot load env file: %s", err.Error())
	}
	return &Envar{
		Vars: outMap,
	}, nil
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

func (e *Envar) Merge(env *Envar) {
	for key, value := range env.Vars {
		e.Vars[key] = value
	}
}

func (e *Envar) String() string {
	buffer := bytes.Buffer{}
	for key, value := range e.Vars {
		buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	return buffer.String()
}
