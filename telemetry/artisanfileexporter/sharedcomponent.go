/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package artisanfileexporter

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// SharedComponents a map that keeps reference of all created instances for a given configuration,
// and ensures that the shared state is started and stopped only once.
type SharedComponents struct {
	comps map[interface{}]*SharedComponent
}

// NewSharedComponents returns a new empty SharedComponents.
func NewSharedComponents() *SharedComponents {
	return &SharedComponents{
		comps: make(map[interface{}]*SharedComponent),
	}
}

// GetOrAdd returns the already created instance if exists, otherwise creates a new instance
// and adds it to the map of references.
func (scs *SharedComponents) GetOrAdd(key interface{}, create func() component.Component) *SharedComponent {
	if c, ok := scs.comps[key]; ok {
		return c
	}
	newComp := &SharedComponent{
		Component: create(),
		removeFunc: func() {
			delete(scs.comps, key)
		},
	}
	scs.comps[key] = newComp
	return newComp
}

// SharedComponent ensures that the wrapped component is started and stopped only once.
// When stopped it is removed from the SharedComponents map.
type SharedComponent struct {
	component.Component

	startOnce  sync.Once
	stopOnce   sync.Once
	removeFunc func()
}

// Unwrap returns the original component.
func (r *SharedComponent) Unwrap() component.Component {
	return r.Component
}

// Start implements component.Component.
func (r *SharedComponent) Start(ctx context.Context, host component.Host) error {
	var err error
	r.startOnce.Do(func() {
		err = r.Component.Start(ctx, host)
	})
	return err
}

// Shutdown implements component.Component.
func (r *SharedComponent) Shutdown(ctx context.Context) error {
	var err error
	r.stopOnce.Do(func() {
		err = r.Component.Shutdown(ctx)
		r.removeFunc()
	})
	return err
}
