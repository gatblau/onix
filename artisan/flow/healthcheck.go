/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"github.com/gatblau/onix/artisan/core"
	"strings"
)

// run a health check of the flow configuration to determine if there is no conflicting information
// and Artisan can identify the sources to survey inputs
func flowHealthCheck(flow *Flow, step *Step) {
	// check no reserved vars are used
	for _, s := range flow.Steps {
		if s.Input != nil && s.Input.Var != nil {
			for _, v := range s.Input.Var {
				if strings.HasPrefix(v.Name, "OXART_") {
					core.RaiseErr("variable name %s is reserved for Artisan use, choose a different name")
				}
			}
		}
		if s.Input != nil && s.Input.Secret != nil {
			for _, s := range s.Input.Secret {
				if strings.HasPrefix(s.Name, "OXART_") {
					core.RaiseErr("secret name %s is reserved for Artisan use, choose a different name")
				}
			}
		}
	}
	// if there are sources there must be a create
	sourcesCount, createCount := 0, 0
	createFirstIx, mergeFirstIx, readFirstIx := -1, -1, -1
	invalidSourceValue, invalidSourceStep := "", ""
	for ix, s := range flow.Steps {
		if len(s.PackageSource) > 0 {
			sourcesCount++
			if s.PackageSource == "create" {
				if createFirstIx < 0 {
					createFirstIx = ix
				}
				createCount++
			} else if s.PackageSource == "read" {
				if readFirstIx < 0 {
					readFirstIx = ix
				}
			} else if s.PackageSource == "merge" {
				if mergeFirstIx < 0 {
					mergeFirstIx = ix
				}
			} else {
				// catch a wrong source name
				invalidSourceStep = s.Name
				invalidSourceValue = s.PackageSource
			}
		}
	}
	// invalid source value
	if len(invalidSourceValue) > 0 {
		core.RaiseErr("step '%s' has an invalid source value '%s': permitted values are 'create', 'merge' or 'read'", invalidSourceStep, invalidSourceValue)
	}
	// check that if there are sources, then all sources are set
	if sourcesCount > 0 && createCount == 0 {
		core.RaiseErr("missing a create value for step source, steps with read or merge source values require at least one step before them with a source set to create")
	}
	// read before create
	if readFirstIx > -1 && readFirstIx < createFirstIx {
		core.RaiseErr("a step defines source=read before another step defines source=create, a step with source=create must be set before a step with source=read")
	}
	// merge before create
	if mergeFirstIx > -1 && mergeFirstIx < createFirstIx {
		core.RaiseErr("a step defines source=merge before another step defines source=create, a step with source=create must be set before a step with source=merge")
	}
	// check for a package source
	packageSource := false
	packageSourceStepName := ""
	for _, s := range flow.Steps {
		if s.PackageSource == "create" {
			packageSource = true
			packageSourceStepName = s.Name
			break
		}
	}
	// check for git source
	gitSource := false
	gitSourceStepName := ""
	for _, s := range flow.Steps {
		if len(s.Package) == 0 && len(s.Function) > 0 {
			gitSource = true
			gitSourceStepName = s.Name
			break
		}
	}
	// missing package name
	missingPackage := false
	missingPackageStepName := ""
	for _, s := range flow.Steps {
		// source & function without package
		if len(s.PackageSource) > 0 && len(s.Function) > 0 && len(s.Package) == 0 {
			missingPackage = true
			missingPackageStepName = s.Name
			break
		}
	}
	// if missing a package name on a step
	if missingPackage {
		core.RaiseErr(`the flow step '%s' is missing a package name.`, missingPackageStepName)
	}
	// missing function name
	missingFunction := false
	missingFunctionStepName := ""
	for _, s := range flow.Steps {
		// missing function when package does not exists and source is not merge
		if len(s.Function) == 0 && len(s.Package) > 0 && strings.ToLower(s.PackageSource) != "merge" {
			missingFunction = true
			missingFunctionStepName = s.Name
			break
		}
	}
	// if missing a function name on a step
	if missingFunction {
		core.RaiseErr(`the flow step '%s' is missing a function name.`, missingFunctionStepName)
	}
	// check if two steps using the same package require both transient and persistent sources
	transient := false
	transientStepName := ""
	transientPackageName := ""
	// check if transient source is used on a package
	for _, s := range flow.Steps {
		if len(s.PackageSource) == 0 && len(s.Function) > 0 && len(s.Package) > 0 {
			transient = true
			transientStepName = s.Name
			transientPackageName = s.Package
			break
		}
	}
	// check if persistent source is used on the same package that was used with a transient source
	persistent := false
	persistentStepName := ""
	for _, s := range flow.Steps {
		if len(s.PackageSource) > 0 && len(s.Function) > 0 && len(s.Package) > 0 && s.Package == transientPackageName {
			persistent = true
			persistentStepName = s.Name
			break
		}
	}
	if transient && persistent {
		core.RaiseErr(`the flow suggest that the same package '%s' is used both with a transient source 
(i.e. step '%s') and a persistent source (i.e. step '%s')
Which of the two is correct?`, transientPackageName, transientStepName, persistentStepName)
	}
	// a flow must not have two orthogonal sources such as package and git
	if packageSource && gitSource {
		core.RaiseErr(`the flow suggests that a GIT source and a PACKAGE source are required.

A flow can read files in one of three ways:

1) Every step uses the transient files in the package defined for the step 
   => each step defines package and function names

2) All steps use the files in a git repository 
   => the step defines only function name

3) The first step provides a persistent source and the other steps either read or merge onto that source 
   each step defines package and function name and source property as follows:
   => the first step defines a source=create
   => subsequent steps define source=read or source=merge
   NOTE: use 'merge' option if you want to merge additional files from a different package into the source,
     otherwise use 'read' option

You must ensure that the flow is setup in one of the ways described above.

Debugging Information:
- the step '%s' suggests a GIT source is required as only a function name is provided.
- the step '%s' suggests a PACKAGE source is required as a source attribute is provided.

Which of the above two scenarios is correct?
`, gitSourceStepName, packageSourceStepName)
	}
}
