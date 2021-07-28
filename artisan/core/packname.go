/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"fmt"
	"strings"
)

// defines the name of a package
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

func ParseName(packageName string) (*PackageName, error) {
	domain, group, name, tag, err := breakName(packageName)
	n := &PackageName{
		Domain: domain,
		Group:  group,
		Name:   name,
		Tag:    tag,
	}
	if err != nil {
		return nil, err
	}
	// set defaults if there are missing values
	if len(n.Domain) == 0 {
		n.Domain = "art1san.net"
	}
	if len(n.Group) == 0 {
		n.Group = "root"
	}
	if len(n.Tag) == 0 {
		n.Tag = "latest"
	}
	return n, nil
}

func breakName(packageName string) (domain, group, name, tag string, err error) {
	components := strings.Split(packageName, "/")
	for i, component := range components {
		// name:tag
		if i == len(components)-1 {
			parts := strings.Split(component, ":")
			if len(parts) == 2 {
				tag = parts[1]
				if !tagRegexp.MatchString(tag) {
					err = fmt.Errorf("package name %s: tag %s is invalid", name, tag)
					return
				}
			}
			name = parts[0]
			if !domainComponentRegexp.MatchString(name) {
				err = fmt.Errorf("package name %s: component part %s is invalid", packageName, name)
				return
			}
		} else if i == 0 {
			domain = component
			if !domainComponentRegexp.MatchString(domain) {
				err = fmt.Errorf("package name %s: domain %s is invalid", packageName, domain)
				return
			}
		} else {
			group += fmt.Sprintf("%s/", component)
			if !domainComponentRegexp.MatchString(component) {
				err = fmt.Errorf("package name %s: component part %s is invalid", packageName, component)
				return
			}
		}
	}
	if len(group) > 0 {
		group = group[:len(group)-1]
	}
	return
}
