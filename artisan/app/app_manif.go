/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Manifest the application manifest that is made up of one or more service manifests
type Manifest struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Version     string   `yaml:"version"`
	Profiles    Profiles `yaml:"profiles"`
	Services    Services `yaml:"services"`
	Var         Vars     `yaml:"var,omitempty"`
}

type Profile struct {
	Name     string   `yaml:"name"`
	Services []string `yaml:"services"`
}

type Profiles []Profile

func (p *Profiles) Get(name string) *Profile {
	for _, profile := range *p {
		if profile.Name == name {
			return &profile
		}
	}
	return nil
}

// servicesSlice returns a string slice with service names in the profile
func (p *Profile) servicesSlice() []string {
	result := make([]string, 0)
	for _, svc := range p.Services {
		result = append(result, svc)
	}
	return result
}

type Services []SvcRef

func (s *SvcRef) InProfile(profile []string) bool {
	for _, p := range profile {
		if s.Name == p {
			return true
		}
	}
	return false
}

type SvcRef struct {
	// the name of the service
	Name string `yaml:"name,omitempty"`
	// the service description
	Description string `yaml:"description"`
	// the uri of the service manifest
	URI string `yaml:"uri,omitempty"`
	// the uri of the database schema definition (if any)
	SchemaURI string `yaml:"schema_uri,omitempty"`
	// the URI of the service image containing the service manifest
	Image string `yaml:"image,omitempty"`
	// whether this service should not be publicly exposed, by default is false
	Private bool `yaml:"private,omitempty"`
	// the service port, if not specified, the application port (in the service manifest) is used
	Port string `yaml:"port,omitempty"`
	// the service manifest loaded from remote image
	Info *SvcManifest `yaml:"service,omitempty"`
	// the other services it depends on
	DependsOn []string `yaml:"depends_on,omitempty"`
	// the other services using it
	UsedBy []string `yaml:"used_by_count,omitempty"`
}

// NewAppMan creates a new application manifest from an URI (supported schemes are http(s):// and file://
func NewAppMan(uri, profile string) (man *Manifest, err error) {
	if ok, path := isFile(uri); ok {
		man, err = loadFromFile(path)
	} else if isURL(uri) {
		man, err = loadFromURL(uri)
	}
	if err != nil {
		return
	}
	if man == nil {
		return nil, fmt.Errorf("invalid URI value '%s': should start with either file://, http:// or https://\n", uri)
	}
	// first trim the services in the manifest to only the one defined in the requested profile
	if man, err = man.trim(profile); err != nil {
		return
	}
	//  then fetches remote service manifests
	if man, err = man.explode(); err != nil {
		return
	}
	// finally, evaluates functions and bindings (wire all service dependencies by evaluating dependent variables)
	if man, err = man.wire(); err != nil {
		return
	}
	return
}

// trim the services to a specific profile
// if no profile is specified and profiles exist then use the first profile in the manifest
// if no profile is specified and profiles do not exist return an error
// if no profile is specified and no profiles exist in the manifest then does not trim (all services in the manifest are included)
func (m *Manifest) trim(profile string) (*Manifest, error) {
	var prof *Profile
	// if no profiles have been defined then perform no trimming
	if m.Profiles == nil || len(m.Profiles) == 0 {
		// if a specific profile was requested
		if len(profile) > 0 {
			return nil, fmt.Errorf("no profiles have been defined in application manifest\n")
		}
		// else return untrimmed
		return m, nil
	}
	// if no specific profile was requested, and profiles are defined, use the first one in the list
	if len(profile) == 0 {
		prof = &m.Profiles[0]
	} else {
		// try and get the requested profile
		prof = m.Profiles.Get(profile)
		// if the profile was not found
		if prof == nil {
			return nil, fmt.Errorf("profile '%s' was not found in the application manifest '%s'\n", profile, m.Name)
		}
	}
	// get a list of service names in the profile
	profServices := prof.servicesSlice()
	// deep clone the manifest
	appMan := new(Manifest)
	_ = m.deepCopy(appMan)
	// reset the services list
	appMan.Services = make(Services, 0)
	// a re-populate with the items in the requested profile
	for _, svc := range m.Services {
		for _, profSvc := range profServices {
			if profSvc == svc.Name {
				appMan.Services = append(appMan.Services, svc)
			}
		}
	}
	return appMan, nil
}

// explode adds service manifest information to the application manifest by querying remote sources
func (m *Manifest) explode() (*Manifest, error) {
	var err error
	// create a copy of the passed in light manifest to become the exploded version
	appMan := new(Manifest)
	_ = m.deepCopy(appMan)
	// validate the app manifest
	if err = m.validate(); err != nil {
		return nil, err
	}
	// loop through
	var svcMan *SvcManifest
	for i, svc := range m.Services {
		// image only
		if len(svc.Image) > 0 && len(svc.URI) == 0 {
			svcMan, err = loadSvcManFromImage(svc)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		} else if len(svc.Image) > 0 && len(svc.URI) > 0 {
			svcMan, err = loadSvcManFromURI(svc)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		}
		appMan.Services[i].Info = svcMan
		appMan.Services[i].Name = svcMan.Name
	}
	return appMan, nil
}

// wire evaluates all expressions in the service manifest (i.e. functions and bindings) and work out service dependencies
func (m *Manifest) wire() (*Manifest, error) {
	appMan := new(Manifest)
	_ = m.deepCopy(appMan)
	// do the wiring
	for six, service := range m.Services {
		for vix, v := range service.Info.Var {
			// if the variable is a function expression
			if strings.HasPrefix(strings.Replace(v.Value, " ", "", -1), "{{fx=") {
				content := v.Value[len("{{fx=") : len(v.Value)-2]
				parts := strings.Split(content, ":")
				// qualifies the name of the variable with the service name
				// variable name without ${...} wrapper
				vName := fmt.Sprintf("%s_%s", strings.ToUpper(strings.Replace(service.Name, "-", "_", -1)), v.Name)
				// variable name wrapped with ${...}
				vNameWrapped := fmt.Sprintf("${%s}", vName)
				switch strings.ToLower(parts[0]) {
				case "pwd":
					subParts := strings.Split(parts[1], ",")
					length, _ := strconv.Atoi(subParts[0])
					symbols, _ := strconv.ParseBool(subParts[1])
					appMan.Services[six].Info.Var[vix].Value = vNameWrapped
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        vName,
						Description: v.Description,
						Value:       RandomPwd(length, symbols),
						Secret:      true,
						Service:     strings.ToUpper(service.Name),
					})
				case "name":
					number, _ := strconv.Atoi(parts[1])
					appMan.Services[six].Info.Var[vix].Value = vNameWrapped
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        vName,
						Description: v.Description,
						Value:       RandomName(number),
						Secret:      false,
						Service:     strings.ToUpper(service.Name),
					})
				default:
					return nil, fmt.Errorf("invalid function %s='%s' in service '%s'\n", v.Name, v.Value, service.Name)
				}
			} else { // if the variable is a binding
				b := bindings(v.Value)
				for _, binding := range b {
					content := binding[len("{{bind=") : len(binding)-2]
					parts := strings.Split(content, ":")
					switch len(parts) {
					case 1:
						svcName := parts[0]
						// check the name exists
						found := false
						for _, s := range m.Services {
							if s.Name == svcName {
								found = true
								break
							}
						}
						if !found {
							return nil, fmt.Errorf("invalid service name '%s' => %s='%s' in service '%s'\n", svcName, v.Name, v.Value, service.Name)
						}
						appMan.Services[six].Info.Var[vix].Value = strings.Replace(appMan.Services[six].Info.Var[vix].Value, binding, svcName, 1)
						appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, svcName)
						ix := getServiceIx(*appMan, svcName)
						appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name)
					case 2:
						switch parts[1] {
						case "schema_uri":
							if uri := m.getSchemaURI(parts[0]); len(uri) > 0 {
								appMan.Services[six].Info.Var[vix].Value = uri
							} else {
								return nil, fmt.Errorf("variable %s='%s' in service '%s' request schema_ui from service '%s' but is missing\n", v.Name, v.Value, service.Name, parts[0])
							}
						case "port":
							if port := m.getSvcPort(parts[0]); len(port) > 0 {
								appMan.Services[six].Info.Var[vix].Value = strings.Replace(appMan.Services[six].Info.Var[vix].Value, binding, port, 1)
								appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, port)
								ix := getServiceIx(*appMan, parts[0])
								appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name)
							} else {
								return nil, fmt.Errorf("port not defined in service '%s', invoked from variable %s => '%s' in service %s\n", parts[0], v.Name, binding, service.Name)
							}
						default:
							return nil, fmt.Errorf("invalid binding %s='%s' in service '%s'\n", v.Name, binding, service.Name)
						}
					case 3:
						switch parts[1] {
						case "var":
							if m.varExists(parts[2]) {
								appMan.Services[six].Info.Var[vix].Value = strings.Replace(appMan.Services[six].Info.Var[vix].Value, binding, strings.ToUpper(fmt.Sprintf("${%s_%s}", parts[0], parts[2])), 1)
								appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, parts[0])
								ix := getServiceIx(*appMan, parts[0])
								appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name)
							} else {
								return nil, fmt.Errorf("cannot find variable %s='%s' in service '%s'\n", v.Name, v.Value, service.Name)
							}
						default:
							return nil, fmt.Errorf("invalid binding %s='%s' in service '%s'\n", v.Name, v.Value, service.Name)
						}
					default:
						return nil, fmt.Errorf("invalid binding %s='%s' in service '%s'\n", v.Name, v.Value, service.Name)
					}
				}
			}
		}
	}
	// sort the services by dependencies (most widely used first)
	sort.Slice(m.Services, func(i, j int) bool {
		return len(m.Services[i].UsedBy) > len(m.Services[j].UsedBy)
	})
	return appMan, nil
}

func addDependency(dependsOn []string, svc string) []string {
	result := make([]string, len(dependsOn))
	copy(result, dependsOn)
	exists := false
	for _, d := range result {
		if d == svc {
			exists = true
			break
		}
	}
	if !exists {
		result = append(result, svc)
	}
	return result
}

func (m *Manifest) getSchemaURI(svc string) string {
	for _, service := range m.Services {
		if service.Name == svc && len(service.SchemaURI) > 0 {
			return service.SchemaURI
		}
	}
	return ""
}

func getServiceIx(m Manifest, svcName string) int {
	for ix, service := range m.Services {
		if service.Name == svcName {
			return ix
		}
	}
	return -1
}

func (m *Manifest) varExists(varName string) bool {
	for _, service := range m.Services {
		for _, v := range service.Info.Var {
			if v.Name == varName {
				return true
			}
		}
	}
	return false
}

func (m *Manifest) validate() error {
	for _, svc := range m.Services {
		// case of manifest embedded in docker image then no URI is needed (image only)
		// case of manifest in git repo (uri + image required)
		// so cases to avoid is uri only
		if len(svc.Image) == 0 && len(svc.URI) > 0 {
			return fmt.Errorf("invalid entry for service '%s' manifest in application manifest: only one of Image or URI attributes must be specified\n", svc.Name)
		}
		// or uri & image not provided
		if len(svc.Image) == 0 && len(svc.URI) == 0 {
			return fmt.Errorf("invalid entry for service '%s' manifest in application manifest: either one of Image or URI attributes must be specified\n", svc.Name)
		}
	}
	return nil
}

func (m *Manifest) deepCopy(dst interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(m); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}

func (m *Manifest) getSvcPort(svcName string) string {
	for _, service := range m.Services {
		if service.Name == svcName && len(service.Info.Port) > 0 {
			return service.Info.Port
		}
	}
	return ""
}

func bindings(value string) []string {
	r, _ := regexp.Compile("{{bind=(?P<NAME>[^}]+)}}")
	return r.FindAllString(value, -1)
}
