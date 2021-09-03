package merge

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"reflect"
)

// Context the merge context for artisan templates .art
type Context struct {
	loader loader
	// the selected variable group for a range
	currentGroup string
	// a list of variable sets
	Items []Set
}

func NewContext(env Envar) (*Context, error) {
	ctx := &Context{
		loader: NewLoader(env),
		Items:  []Set{},
	}
	return ctx, nil
}

// Var return the value of a variable
func (c *Context) Var(name reflect.Value) reflect.Value {
	return reflect.ValueOf(c.loader.vars[name.String()])
}

// Select select a specific variable group and populate all variable sets within the group
func (c *Context) Select(group reflect.Value) reflect.Value {
	c.currentGroup = group.String()
	ii := c.loader.indices(c.currentGroup)
	c.Items = []Set{}
	for i := 1; i <= ii; i++ {
		c.Items = append(c.Items, c.loader.set(c.currentGroup, fmt.Sprint(i), c))
	}
	return reflect.ValueOf("")
}

// Item return a grouped variable value using its name and the current iteration set
func (c *Context) Item(name reflect.Value, set reflect.Value) reflect.Value {
	s, ok := set.Interface().(Set)
	if !ok {
		panic("Item function requires a set for the first parameter\n")
	}
	return reflect.ValueOf(s.Value[name.String()])
}

func (c *Context) GroupExists(group reflect.Value) bool {
	ix := c.loader.indices(group.String())
	return ix != 0
}
