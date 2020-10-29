package core

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// defines the name of an artefact
// domain/repository/name:tag
type ArtieName struct {
	Domain     string
	Repository string
	Name       string
	Tag        string
}

func (a *ArtieName) String() string {
	return fmt.Sprintf("%s/%s/%s:%s", a.Domain, a.Repository, a.Name, a.Tag)
}

func (a *ArtieName) FullyQualifiedRepository() string {
	return fmt.Sprintf("%s/%s", a.Domain, a.Repository)
}

func (a *ArtieName) FullyQualifiedName() string {
	return fmt.Sprintf("%s/%s/%s", a.Domain, a.Repository, a.Name)
}

func (a *ArtieName) Path() string {
	return fmt.Sprintf("%s/%s", a.Repository, a.Name)
}

func ParseName(name string) *ArtieName {
	n := &ArtieName{}
	components := strings.Split(name, "/")
	// validate component elements
	for _, component := range components {
		if !domainComponentRegexp.MatchString(component) {
			log.Fatal(fmt.Sprintf("artefact name %s: component part %s is invalid", name, component))
		}
	}
	if len(components) == 1 {
		n.Name = components[0]
	}
	if len(components) == 2 {
		log.Fatal(errors.New("invalid artefact URI"))
	}
	if len(components) == 3 {
		n.Domain = components[0]
		n.Repository = components[1]
		parts := strings.Split(components[2], ":")
		if len(parts) == 2 {
			n.Tag = parts[1]
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
		n.Domain = "artie.library"
	}
	if len(n.Repository) == 0 {
		n.Repository = "library"
	}
	if len(n.Tag) == 0 {
		n.Tag = "latest"
	}
	return n
}
