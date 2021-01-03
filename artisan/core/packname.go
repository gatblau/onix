/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"fmt"
	"log"
	"strings"
)

// defines the name of an artefact
// domain/group/name:tag
type PackageName struct {
	Domain string
	Group  string
	Name   string
	Tag    string
}

func (a *PackageName) IsInTheSameRepositoryAs(name *PackageName) bool {
	return a.FullyQualifiedName() == name.FullyQualifiedName()
}

func (a *PackageName) String() string {
	return fmt.Sprintf("%s/%s/%s:%s", a.Domain, a.Group, a.Name, a.Tag)
}

func (a *PackageName) FullyQualifiedGroup() string {
	return fmt.Sprintf("%s/%s", a.Domain, a.Group)
}

func (a *PackageName) FullyQualifiedName() string {
	return fmt.Sprintf("%s/%s/%s", a.Domain, a.Group, a.Name)
}

func (a *PackageName) Repository() string {
	return fmt.Sprintf("%s/%s", a.Group, a.Name)
}

func ParseName(name string) *PackageName {
	n := &PackageName{}
	components := strings.Split(name, "/")
	// validate component elements
	for _, component := range components {
		if !domainComponentRegexp.MatchString(component) {
			log.Fatal(fmt.Sprintf("artefact name %s: component part %s is invalid", name, component))
		}
	}
	if len(components) == 1 {
		parts := strings.Split(components[0], ":")
		if len(parts) == 2 {
			n.Name = parts[0]
			n.Tag = parts[1]
		} else {
			n.Name = components[0]
		}
	} else if len(components) == 2 {
		parts := strings.Split(components[1], ":")
		n.Group = components[0]
		if len(parts) == 2 {
			n.Name = parts[0]
			n.Tag = parts[1]
		} else {
			n.Name = components[1]
		}
	} else if len(components) == 3 {
		n.Domain = components[0]
		n.Group = components[1]
		parts := strings.Split(components[2], ":")
		switch len(parts) {
		case 2:
			n.Tag = parts[1]
		case 1:
			n.Tag = "latest"
		}
		n.Name = parts[0]
	}
	// validate
	if len(n.Domain) > 0 {
		if !domainComponentRegexp.MatchString(n.Domain) {
			log.Fatal(fmt.Sprintf("artefact name %s: domain %s is invalid", name, n.Domain))
		}
	}
	if len(n.Tag) > 0 {
		if !tagRegexp.MatchString(n.Tag) {
			log.Fatal(fmt.Sprintf("artefact name %s: tag %s is invalid", name, n.Tag))
		}
	}
	// set defaults if there are missing values
	if len(n.Domain) == 0 {
		n.Domain = "artisan.library"
	}
	if len(n.Group) == 0 {
		n.Group = "library"
	}
	if len(n.Tag) == 0 {
		n.Tag = "latest"
	}
	return n
}
