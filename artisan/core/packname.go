package core

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"net"
	"regexp"
	"strings"
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

func splitName(packageName string) (domain, group, name, tag string, err error) {
	parts := strings.Split(packageName, "/")
	switch len(parts) {
	case 1:
		name, tag, err = parseNameTag(parts[0])
		domain = "artr.gdn"
		group = "lib"
	case 2:
		domain = "artr.gdn"
		group = parts[0]
		name, tag, err = parseNameTag(parts[1])
	case 3:
		domain = parts[0]
		group = parts[1]
		name, tag, err = parseNameTag(parts[2])
	default:
		l := len(parts)
		domain = parts[0]
		group = strings.Join(parts[1:l-1], "/")
		name, tag, err = parseNameTag(parts[l-1])
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
