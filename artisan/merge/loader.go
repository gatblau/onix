package merge

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"strconv"
	"strings"
)

type loader struct {
	items []item
	vars  map[string]string
}

type item struct {
	group string
	name  string
	index string
	value string
}

func NewLoader(env Envar) loader {
	l := &loader{
		items: []item{},
		vars:  map[string]string{},
	}
	for key, value := range env.Vars {
		// any variable added to a key-value map
		l.vars[key] = value
		// now processes grouped variables (i.e. following naming convention GROUP__NAME__IX)
		ix1 := strings.Index(key, "__")
		ix2 := strings.LastIndex(key, "__")
		if ix1 > 0 && ix2 > 0 {
			group := key[:ix1]
			name := key[ix1+2 : ix2]
			index := key[ix2+2:]
			l.items = append(l.items, item{
				group: group,
				name:  name,
				index: index,
				value: value,
			})
		}
	}
	return *l
}

func (l *loader) set(group, index string, ctx *Context) Set {
	vars := make(map[string]string)
	for _, i := range l.items {
		if i.group == group && i.index == index {
			vars[i.name] = i.value
		}
	}
	return Set{
		Value:   vars,
		Context: ctx,
	}
}

func (l *loader) indices(group string) int {
	var result int = 0
	for _, i := range l.items {
		ii, _ := strconv.Atoi(i.index)
		if ii > result && i.group == group {
			result = ii
		}
	}
	return result
}
