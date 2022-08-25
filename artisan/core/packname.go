/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	regDomain = "artr.gdn"
	regGroup  = "lib"
)

// PackageName defines the name of a package following the format: domain/group/name:tag
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

func (a *PackageName) NormalString() string {
	return strings.ReplaceAll(fmt.Sprintf("%s_%s_%s_%s", a.Domain, a.Group, a.Name, a.Tag), ".", "_")
}

func (a *PackageName) FullyQualifiedGroup() string {
	return fmt.Sprintf("%s/%s", a.Domain, a.Group)
}

func (a *PackageName) FullyQualifiedName() string {
	return fmt.Sprintf("%s/%s/%s", a.Domain, a.Group, a.Name)
}

func (a *PackageName) FullyQualifiedNameTag() string {
	return fmt.Sprintf("%s/%s/%s:%s", a.Domain, a.Group, a.Name, a.Tag)
}

func (a *PackageName) Repository() string {
	return fmt.Sprintf("%s/%s", a.Group, a.Name)
}

func ParseName(packageName string) (*PackageName, error) {
	domain, group, name, tag, err := splitName(packageName)
	n := &PackageName{
		Domain: domain,
		Group:  group,
		Name:   name,
		Tag:    tag,
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

// ValidateNames get a list of qualified package names from their string representation
func ValidateNames(packages []string) ([]PackageName, error) {
	var names []PackageName
	for _, p := range packages {
		name, err := ParseName(p)
		if err != nil {
			return nil, err
		}
		names = append(names, *name)
	}
	return names, nil
}

// splitName splits the package name string into domain, group, name and tag
// implements the logic for deciding if the package name is correct and if certain portions are missing it
// automatically insert them using the global registry path
func splitName(packageName string) (domain, group, name, tag string, err error) {
	parts := strings.Split(packageName, "/")
	switch len(parts) {
	// it assumes the part is the package name and potentially a tag
	case 1:
		name, tag, err = parseNameTag(parts[0])
		// check the name is not a domain
		if isDomain(name) {
			return "", "", "", "", fmt.Errorf("missing package group and name")
		}
		domain = regDomain
		group = regGroup
		// it assumes the first part is the group and the second part is the name
	case 2:
		domain = regDomain
		group = parts[0]
		// if the group is a domain insert the default group
		if isDomain(group) {
			domain = group
			group = regGroup
		}
		name, tag, err = parseNameTag(parts[1])
		if isDomain(name) {
			return "", "", "", "", fmt.Errorf("the second portion of the package name %s should not be a domain", name)
		}
	// it assumes the first part is a domain, the second part is a group and the third part is a name
	case 3:
		domain = parts[0]
		group = parts[1]
		name, tag, err = parseNameTag(parts[2])
		if isDomain(group) {
			return "", "", "", "", fmt.Errorf("the second portion of the package name %s should not be a domain", group)
		}
		if isDomain(name) {
			return "", "", "", "", fmt.Errorf("the third portion of the package name %s should not be a domain", name)
		}
		// if the domain part does not pass validation then assume it is part of the group
		if !isDomain(domain) {
			group = fmt.Sprintf("%s/%s", domain, group)
			domain = regDomain
		}
	// it assumes the first part is a domain, the last part is a name and the parts in the middle are the group
	default:
		l := len(parts)
		domain = parts[0]
		group = strings.Join(parts[1:l-1], "/")
		name, tag, err = parseNameTag(parts[l-1])
		// if the domain part does not pass validation then assume it is part of the group
		if !isDomain(domain) {
			group = fmt.Sprintf("%s/%s", domain, group)
			domain = regDomain
		}
		if isDomain(name) {
			return "", "", "", "", fmt.Errorf("the last portion of the package name %s should not be a domain", name)
		}
	}
	if !validTag(tag) {
		err = fmt.Errorf("package name %s: tag %s is invalid", name, tag)
		return
	}
	if !validName(name) {
		err = fmt.Errorf("package name %s: name %s is invalid", packageName, name)
		return
	}
	if !validGroup(group) {
		err = fmt.Errorf("package name %s: group %s is invalid", packageName, name)
		return
	}
	if !validDomain(domain) {
		err = fmt.Errorf("package name %s: domain %s is invalid", packageName, domain)
		return
	}
	return
}

func parseNameTag(nameTag string) (name, tag string, err error) {
	parts := strings.Split(nameTag, ":")
	switch len(parts) {
	case 1:
		name = parts[0]
		tag = "latest"
	case 2:
		name = parts[0]
		tag = parts[1]
	default:
		err = fmt.Errorf("invalid name:tag %s", nameTag)
	}
	return
}

// validTag validates an artisan tag
// A tag name must be valid ASCII and may contain lowercase and uppercase letters, digits, underscores, periods and dashes.
// A tag name may not start with a period or a dash and may contain a maximum of 128 characters.
func validTag(tag string) (valid bool) {
	valid, _ = regexp.MatchString(`^([a-zA-Z0-9\._-]{0,62})*$`, tag)
	return
}

func validDomain(domain string) bool {
	// breaks down port
	parts := strings.Split(domain, ":")
	switch len(parts) {
	// no port
	case 1:
		return validDNSIP(parts[0])
	// check ip + port
	case 2:
		numeric, _ := regexp.MatchString("^[0-9]+$", parts[1])
		return validDNSIP(parts[0]) && numeric
	}
	return true
}

func validDNSIP(domain string) bool {
	return (len(domain) > 2 &&
		len(domain) < 64 &&
		validDNS(domain) &&
		domain[0] != '-' &&
		domain[len(domain)-1] != '-' &&
		!strings.Contains(domain, ":") &&
		!hasScheme(domain)) ||
		(net.ParseIP(domain) != nil)
}

func hasScheme(name string) (hasScheme bool) {
	hasScheme, _ = regexp.MatchString(`((ftp|tcp|udp|wss?|https?):\/\/)`, name)
	return
}

func validDNS(dns string) (valid bool) {
	valid, _ = regexp.MatchString(`^([a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})(\.[a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})*[\._]?$`, dns)
	return
}

func validGroup(group string) (valid bool) {
	valid, _ = regexp.MatchString(`^([a-zA-Z0-9_-]{0,62}\/?)*$`, group)
	return
}

func validName(name string) (valid bool) {
	valid, _ = regexp.MatchString(`^([a-zA-Z0-9_-]{0,62})*$`, name)
	return
}

func isDomain(name string) (valid bool) {
	valid, _ = regexp.MatchString(`^.*:\d*$|^.*\.[a-zA-Z].*$|^\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b$`, name)
	return valid
}
