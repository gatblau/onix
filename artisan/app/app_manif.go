/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gatblau/onix/artisan/release"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/gatblau/onix/artisan/app/behaviour"
	"github.com/gatblau/onix/oxlib/resx"
	"gopkg.in/yaml.v2"
)

// Manifest the application manifest that is made up of one or more service manifests
type Manifest struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Version     string   `yaml:"version"`
	Profiles    Profiles `yaml:"profiles"`
	Services    Services `yaml:"services"`
	Var         Vars     `yaml:"var,omitempty"`
	// specification for images and packages
	Spec *release.Spec `yaml:"spec,omitempty"`

	// internal use
	// for git credentials if required
	credentials string
	artHome     string
	// the location of this manifest
	Root string
}

type Profile struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Services    []ProfileService `yaml:"services"`
}

type Profiles []Profile

type ProfileService struct {
	Name string                                `yaml:"name"`
	Is   map[behaviour.ServiceBehaviour]string `yaml:"is"`
}

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

		result = append(result, svc.Name)
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
	// the service port(s), if not specified, the application port (in the service manifest) is used
	// it can be a string value for a single port or a map for multiple ports
	Port interface{} `yaml:"port"`
	// the service manifest loaded from remote image
	Info *SvcManifest `yaml:"service,omitempty"`
	// the other services it depends on
	DependsOn []string `yaml:"depends_on,omitempty"`
	// the other services using it
	UsedBy []string `yaml:"used_by_count,omitempty"`
	// instructions to customise deployment
	Is map[behaviour.ServiceBehaviour]string `yaml:"is,omitempty"`
}

// PortMap return a parsed map of ports for the port attribute
func (s *SvcRef) PortMap() (map[string]int, error) {
	ports := map[string]int{}
	// if no port is specified at the service level in the app manifest then use the one in the service manifest
	if s.Port == nil {
		s.Port = s.Info.Port
	}
	if p, isString := s.Port.(string); isString {
		value, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		ports = map[string]int{
			"default": value,
		}
	} else if p, isMap := s.Port.(map[interface{}]interface{}); isMap {
		for key, value := range p {
			iv, err := strconv.Atoi(value.(string))
			if err != nil {
				return nil, err
			}
			ports[key.(string)] = iv
		}
	} else {
		return nil, fmt.Errorf("invalid port value: %s", s.Port)
	}
	return ports, nil
}

// NewAppMan creates a new application manifest from a URI (supported schemes are http(s)://, s3(s):// and file
func NewAppMan(uri, profile, credentials, artHome string) (man *Manifest, err error) {
	// validate credentials
	if len(credentials) > 0 {
		if len(strings.Split(credentials, ":")) != 2 {
			err = fmt.Errorf("invalid crdentials format: requires 'user:password'\n")
		}
	}
	// generically load the manifest
	manBytes, err := resx.ReadFile(uri, credentials)
	if err != nil {
		return
	}
	man = new(Manifest)
	err = yaml.Unmarshal(manBytes, man)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal app manifest: %s\n", err)
	}
	if man == nil {
		return nil, fmt.Errorf("invalid URI value '%s': should start with either file://, http:// or https://\n", uri)
	}
	// set the location of the manifest for later file resolution
	man.Root = filepath.Dir(uri)
	// set any credentials if provided
	man.credentials = credentials
	// set spec
	spec, err := man.loadSpec(uri)
	if err != nil {
		return
	}
	man.Spec = spec
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
	// deep clone the manifest
	appMan := new(Manifest)
	_ = m.deepCopy(appMan)
	// reset the services list
	appMan.Services = make(Services, 0)
	// a re-populate with the items in the requested profile
	for _, svc := range m.Services {
		for _, profSvc := range prof.Services {
			if profSvc.Name == svc.Name {
				appMan.Services = append(appMan.Services, svc)
				// if the profile service has behaviours then override the ones defined in the service
				if profSvc.Is != nil {
					// if the service does not define behaviours
					if svc.Is == nil {
						// use the ones in the profile
						svc.Is = profSvc.Is
					} else {
						// override only the behaviours defined in the profile
						for behaviour, value := range profSvc.Is {
							svc.Is[behaviour] = value
						}
					}
				}
				continue
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
			svcMan, err = loadSvcManFromImage(svc, m.artHome)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		} else if len(svc.URI) > 0 {
			// work out the location of the
			// if the URI is for the file system and is a relative path
			if resx.IsFile(svc.URI) && !strings.HasPrefix(svc.URI, "/") {
				svcUri, _ := filepath.Abs(path.Join(m.Root, svc.URI))
				svc.URI = svcUri
			}
			svcMan, err = loadSvcManFromURI(&svc, m.credentials)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		}
		appMan.Services[i].Info = svcMan
		// if profiles are defined then service name sin app manifest should match the ones in the service manifests
		if appMan.Profiles != nil && appMan.Services[i].Name != svcMan.Name {
			return nil, fmt.Errorf("service name mismatch: app manifest => %s; svc manifest => %s\n", appMan.Services[i].Name, svcMan.Name)
		}
		// if no profiles set, and service name in app manifest not set, then set it with the value in the service manifest
		if appMan.Profiles == nil && len(appMan.Services[i].Name) == 0 {
			appMan.Services[i].Name = svcMan.Name
		}
	}
	for _, service := range appMan.Services {
		// binding in a service cannot point to another binding but a value or a function expression
		if err = validateBindings(*appMan, service); err != nil {
			return nil, err
		}
		// check port is defined in service manifest and the service declaration in the app manifest is not a job (image-less service)
		if service.Info.Port == nil && len(service.Image) > 0 {
			return nil, fmt.Errorf("port not defined in service '%s' manifest\n", service.Info.Name)
		}
		// in the case the port in the app manifest not defined
		if service.Port == nil {
			// set it to the port in the svc manifest (through mapping)
			service.Port = service.Info.Port
		} else {
			// check that the port attributes in both the application manifest service declaration and in the service
			// manifests are either a single port string value or a multi-port map value
			_, publishedIsString := service.Port.(string)
			_, publishedIsMap := service.Port.(map[interface{}]interface{})
			_, targetIsString := service.Info.Port.(string)
			_, targetIsMap := service.Info.Port.(map[interface{}]interface{})
			if !(publishedIsString && targetIsString || publishedIsMap && targetIsMap) {
				return nil, fmt.Errorf("port type mistatch between application and service manifests in '%s': check that the port attributes in both manifests are set to either string or map", service.Name)
			}
		}
	}
	return appMan, nil
}

// wire evaluates all expressions in the service manifest (i.e. functions and serviceBindings) and work out service dependencies
func (m *Manifest) wire() (*Manifest, error) {
	appMan := new(Manifest)
	err := m.deepCopy(appMan)
	if err != nil {
		return nil, err
	}
	// do the wiring of expressions in the service manifests
	for six, service := range m.Services {
		// wire expressions in variables
		for vix, v := range service.Info.Var {
			// if the variable is a function expression
			if strings.Contains(strings.Replace(v.Value, " ", "", -1), "${fx=") {
				startAt := strings.Index(v.Value, "${fx=") + len("${fx=")
				endsAt := strings.Index(v.Value, "}")
				content := v.Value[startAt:endsAt]
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
					// add a manifest variable
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        vName,
						Description: v.Description,
						// preserves prefix and suffix for a fx to generate password
						Value:   v.Value[0:startAt-len("${fx=")] + RandomPwd(length, symbols) + v.Value[endsAt+1:],
						Secret:  true,
						Service: strings.ToUpper(service.Name),
					})
				case "name":
					number, _ := strconv.Atoi(parts[1])
					appMan.Services[six].Info.Var[vix].Value = vNameWrapped
					// add a manifest variable
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
				b := serviceBindings(v.Value)
				// if variable is a value add it to the list of manifest variables so that it can be loaded using the .env file
				if len(b) == 0 {
					// add a manifest variable
					vName := fmt.Sprintf("%s_%s", strings.ToUpper(strings.Replace(service.Name, "-", "_", -1)), v.Name)
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        vName,
						Description: v.Description,
						Value:       v.Value,
						Secret:      v.Secret,
						Service:     strings.ToUpper(service.Name),
					})
				} else {
					for _, binding := range b {
						content := binding[len("${bind=") : len(binding)-1]
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
							appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, svcName, appMan.Services[six])
							ix := getServiceIx(*appMan, svcName)
							appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
							// add variable to manifest list
							vName := fmt.Sprintf("%s_%s", strings.ToUpper(strings.Replace(service.Name, "-", "_", -1)), v.Name)
							appMan.Var.Items = append(appMan.Var.Items, AppVar{
								Name:        vName,
								Description: v.Description,
								Value:       appMan.Services[six].Info.Var[vix].Value,
								Secret:      v.Secret,
								Service:     strings.ToUpper(service.Name),
							})
						case 2:
							if strings.HasPrefix(parts[1], "schema_uri") {
								if uri := m.getSchemaURI(parts[0]); len(uri) > 0 {
									appMan.Services[six].Info.Var[vix].Value = uri
								} else {
									return nil, fmt.Errorf("variable %s='%s' in service '%s' request schema_ui from service '%s' but is missing\n", v.Name, v.Value, service.Name, parts[0])
								}
							} else if strings.HasPrefix(parts[1], "port") {
								port, err2 := m.getSvcTargetPort(parts[0], portKey(parts[1]), binding)
								if err2 != nil {
									return nil, err2
								}
								if len(port) > 0 {
									appMan.Services[six].Info.Var[vix].Value = strings.Replace(appMan.Services[six].Info.Var[vix].Value, binding, port, 1)
									vName := fmt.Sprintf("%s_%s", strings.ToUpper(strings.Replace(service.Name, "-", "_", -1)), v.Name)
									found := false
									for ix, item := range appMan.Var.Items {
										if item.Name == vName {
											appMan.Var.Items[ix].Value = appMan.Services[six].Info.Var[vix].Value
											found = true
										}
									}
									if !found {
										appMan.Var.Items = append(appMan.Var.Items, AppVar{
											Name:        vName,
											Description: v.Description,
											Value:       appMan.Services[six].Info.Var[vix].Value,
											Secret:      v.Secret,
											Service:     strings.ToUpper(service.Name),
										})
									}
								} else {
									return nil, fmt.Errorf("missing port in application manifest: service '%s', binding %s\n", service.Name, binding)
								}
							} else {
								return nil, fmt.Errorf("invalid binding %s='%s' in service '%s'\n", v.Name, binding, service.Name)
							}
						case 3:
							switch parts[1] {
							case "var":
								if m.varExists(parts[2]) {
									varKey := strings.ToUpper(fmt.Sprintf("${%s_%s}", strings.Replace(parts[0], "-", "_", -1), parts[2]))
									appMan.Services[six].Info.Var[vix].Value = strings.Replace(appMan.Services[six].Info.Var[vix].Value, binding, varKey, 1)
									appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, parts[0], appMan.Services[six])
									ix := getServiceIx(*appMan, parts[0])
									appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
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
	}
	// do the wiring of expressions in the configuration files
	// note: this is done outside the previous wiring as a first pass is required to collect
	// variable data that can be used to do merging of configuration file templates
	for six, service := range m.Services {
		// wire bind expressions in file templates
		for fix, f := range service.Info.File {
			// if the file has a template
			if len(f.Template) > 0 {
				// merges the file template
				merged, err := appMan.eval(appMan.Services[six].Info.File[fix].Template)
				if err != nil {
					return nil, err
				}
				// set the content to the merged result
				appMan.Services[six].Info.File[fix].Content = merged
				// extract any bindings
				b := serviceBindings(merged)
				for _, binding := range b {
					content := binding[len("${bind=") : len(binding)-1]
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
							return nil, fmt.Errorf("invalid service name '%s' => binding '%s' in service '%s'\n", svcName, binding, service.Name)
						}
						appMan.Services[six].Info.File[fix].Content = strings.Replace(appMan.Services[six].Info.File[fix].Content, binding, svcName, 1)
						appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, svcName, appMan.Services[six])
						ix := getServiceIx(*appMan, svcName)
						appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
					case 2:
						if strings.HasPrefix(parts[1], "schema_uri") {
							if uri := m.getSchemaURI(parts[0]); len(uri) > 0 {
								appMan.Services[six].Info.File[fix].Content = uri
							} else {
								return nil, fmt.Errorf("binding '%s' in service '%s' request schema_ui from service '%s' but is missing\n", binding, service.Name, parts[0])
							}
						} else if strings.HasPrefix(parts[1], "port") {
							port, err := m.getSvcTargetPort(parts[0], portKey(parts[1]), binding)
							if err != nil {
								return nil, err
							}
							if len(port) > 0 {
								appMan.Services[six].Info.File[fix].Content = strings.Replace(appMan.Services[six].Info.File[fix].Content, binding, port, 1)
							} else {
								return nil, fmt.Errorf("port not defined for service '%s' in application manifest, invoked from binding '%s' in service %s\n", parts[0], binding, service.Name)
							}
						} else {
							return nil, fmt.Errorf("invalid binding '%s' in service '%s'\n", binding, service.Name)
						}
					case 3:
						switch parts[1] {
						case "var":
							if appMan.varExists(parts[2]) {
								// get the variable value
								varKey := strings.ToUpper(fmt.Sprintf("%s_%s", strings.Replace(parts[0], "-", "_", -1), parts[2]))
								found := false
								for _, v := range appMan.Var.Items {
									if v.Name == varKey {
										appMan.Services[six].Info.File[fix].Content = strings.Replace(appMan.Services[six].Info.File[fix].Content, binding, v.Value, 1)
										found = true
										break
									}
								}
								if !found {
									return nil, fmt.Errorf("binding in service '%s' points to variable '%s' in service '%s', but it is not defined;\npossible causes:\n - the binding points to the incorrect service;\n - the variable name has been mispelled;\n - the variable needs adding to the service manifest;\n", service.Name, parts[2], parts[0])
								}
								appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, parts[0], appMan.Services[six])
								ix := getServiceIx(*appMan, parts[0])
								appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
							} else {
								return nil, fmt.Errorf("cannot find variable '%s' in service '%s'\n", parts[2], service.Name)
							}
						default:
							return nil, fmt.Errorf("invalid binding '%s' in service '%s'\n", binding, service.Name)
						}
					default:
						return nil, fmt.Errorf("invalid binding '%s' in service '%s'\n", binding, service.Name)
					}
				}
			}
		}
		for _, init := range service.Info.Init {
			for _, script := range init.Scripts {
				i := service.Info.ScriptIx(script)
				// extract any bindings
				b := serviceBindings(appMan.Services[six].Info.Script[i].Content)
				for _, binding := range b {
					content := binding[len("${bind=") : len(binding)-1]
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
							return nil, fmt.Errorf("invalid service name '%s' => binding '%s' in service '%s'\n", svcName, binding, service.Name)
						}
						appMan.Services[six].Info.Script[i].Content = strings.Replace(appMan.Services[six].Info.Script[i].Content, binding, svcName, 1)
						appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, svcName, appMan.Services[six])
						ix := getServiceIx(*appMan, svcName)
						appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
					case 2:
						if strings.HasPrefix(parts[1], "port") {
							port, err := m.getSvcTargetPort(parts[0], portKey(parts[1]), binding)
							if err != nil {
								return nil, err
							}
							if len(port) > 0 {
								appMan.Services[six].Info.Script[i].Content = strings.Replace(appMan.Services[six].Info.Script[i].Content, binding, port, 1)
							} else {
								return nil, fmt.Errorf("port not defined for service '%s' in application manifest, invoked from binding '%s' in service %s\n", parts[0], binding, service.Name)
							}
						} else {
							return nil, fmt.Errorf("invalid binding '%s' in init script for service '%s'\n", binding, service.Name)
						}
					case 3:
						switch parts[1] {
						case "var":
							if appMan.varExists(parts[2]) {
								// get the variable value
								varKey := strings.ToUpper(fmt.Sprintf("%s_%s", strings.Replace(parts[0], "-", "_", -1), parts[2]))
								found := false
								for _, v := range appMan.Var.Items {
									if v.Name == varKey {
										appMan.Services[six].Info.Script[i].Content = strings.Replace(appMan.Services[six].Info.Script[i].Content, binding, fmt.Sprintf("${%s}", varKey), 1)
										found = true
										break
									}
								}
								if !found {
									return nil, fmt.Errorf("binding in init script for service '%s' points to variable '%s' in service '%s', but it is not defined;\npossible causes:\n - the binding points to the incorrect service;\n - the variable name has been mispelled;\n - the variable needs adding to the service manifest;\n", service.Name, parts[2], parts[0])
								}
								appMan.Services[six].DependsOn = addDependency(appMan.Services[six].DependsOn, parts[0], appMan.Services[six])
								ix := getServiceIx(*appMan, parts[0])
								appMan.Services[ix].UsedBy = addDependency(appMan.Services[ix].UsedBy, service.Name, appMan.Services[six])
							} else {
								return nil, fmt.Errorf("cannot find variable '%s' in init script for service '%s'\n", parts[2], service.Name)
							}
						default:
							return nil, fmt.Errorf("invalid binding '%s' in init script for service '%s'\n", binding, service.Name)
						}
					default:
						return nil, fmt.Errorf("invalid binding '%s' in init script for service '%s'\n", binding, service.Name)
					}
				}
			}
		}
		// evaluate expressions in db section, if one has been defined
		if service.Info.Db != nil {
			if strings.HasPrefix(strings.Replace(service.Info.Db.Host, " ", "", -1), "${bind=") {
				content := service.Info.Db.Host[len("${bind=") : len(service.Info.Db.Host)-1]
				parts := strings.Split(content, ":")
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
					return nil, fmt.Errorf("invalid service name '%s' => '%s' in service '%s'\n", svcName, content, service.Name)
				}
				appMan.Services[six].Info.Db.Host = svcName
			}
			// schema uri binding
			if strings.HasPrefix(strings.Replace(service.Info.Db.SchemaURI, " ", "", -1), "${bind=") {
				content := service.Info.Db.SchemaURI[len("${bind=") : len(service.Info.Db.SchemaURI)-1]
				parts := strings.Split(content, ":")
				if uri := m.getSchemaURI(parts[0]); len(uri) > 0 {
					appMan.Services[six].Info.Db.SchemaURI = uri
				} else {
					return nil, fmt.Errorf("schema_uri not defined in app '%s' manifest\n", parts[0])
				}
			}
			// db user name
			if strings.HasPrefix(strings.Replace(service.Info.Db.User, " ", "", -1), "${bind=") {
				content := service.Info.Db.User[len("${bind=") : len(service.Info.Db.User)-1]
				parts := strings.Split(content, ":")
				varKey := strings.ToUpper(fmt.Sprintf("${%s_%s}", strings.Replace(parts[0], "-", "_", -1), parts[2]))
				appMan.Services[six].Info.Db.User = varKey
			}
			// db user pwd
			if strings.HasPrefix(strings.Replace(service.Info.Db.Pwd, " ", "", -1), "${bind=") {
				content := service.Info.Db.Pwd[len("${bind=") : len(service.Info.Db.Pwd)-1]
				parts := strings.Split(content, ":")
				varKey := strings.ToUpper(fmt.Sprintf("${%s_%s}", strings.Replace(parts[0], "-", "_", -1), parts[2]))
				appMan.Services[six].Info.Db.Pwd = varKey
			} else if strings.HasPrefix(strings.Replace(service.Info.Db.Pwd, " ", "", -1), "${fx=") {
				content := service.Info.Db.Pwd[len("${fx=") : len(service.Info.Db.Pwd)-1]
				parts := strings.Split(content, ":")
				if strings.ToLower(parts[0]) == "pwd" {
					subParts := strings.Split(parts[1], ",")
					length, _ := strconv.Atoi(subParts[0])
					symbols, _ := strconv.ParseBool(subParts[1])
					varKey := strings.ToUpper(fmt.Sprintf("${%s_%s}", strings.Replace(appMan.Services[six].Name, "-", "_", -1), "DB_ADMIN_PWD"))
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        varKey,
						Description: fmt.Sprintf("The administrator password to connect to database host '%s'", appMan.Services[six].Info.Db.Host),
						Value:       RandomPwd(length, symbols),
						Secret:      true,
						Service:     strings.ToUpper(appMan.Services[six].Name),
					})
					appMan.Services[six].Info.Db.Pwd = varKey
				}
			}
			// db admin pwd
			if strings.HasPrefix(strings.Replace(service.Info.Db.AdminPwd, " ", "", -1), "${bind=") {
				content := service.Info.Db.AdminPwd[len("${bind=") : len(service.Info.Db.AdminPwd)-1]
				parts := strings.Split(content, ":")
				varKey := strings.ToUpper(fmt.Sprintf("${%s_%s}", strings.Replace(parts[0], "-", "_", -1), parts[2]))
				appMan.Services[six].Info.Db.AdminPwd = varKey
			} else if strings.HasPrefix(strings.Replace(service.Info.Db.AdminPwd, " ", "", -1), "${fx=") {
				content := service.Info.Db.AdminPwd[len("${fx=") : len(service.Info.Db.AdminPwd)-1]
				parts := strings.Split(content, ":")
				if strings.ToLower(parts[0]) == "pwd" {
					subParts := strings.Split(parts[1], ",")
					length, _ := strconv.Atoi(subParts[0])
					symbols, _ := strconv.ParseBool(subParts[1])
					varKey := strings.ToUpper(fmt.Sprintf("%s_%s", strings.Replace(appMan.Services[six].Name, "-", "_", -1), "DB_ADMIN_PWD"))
					appMan.Var.Items = append(appMan.Var.Items, AppVar{
						Name:        varKey,
						Description: fmt.Sprintf("The administrator password to connect to database host '%s'\n", appMan.Services[six].Info.Db.Host),
						Value:       RandomPwd(length, symbols),
						Secret:      true,
						Service:     strings.ToUpper(appMan.Services[six].Name),
					})
					appMan.Services[six].Info.Db.Pwd = varKey
				}
			}
		}
	}
	// merges the spec variables
	for six, service := range m.Services {
		// if the service image has been defined (not image-less service)
		if len(service.Image) > 0 {
			// if the value of Image in the service is not included in the spec
			if !m.Spec.ContainsImage(service.Image) {
				// returns an error
				return nil, fmt.Errorf("service %s should specify an image defined in the spec.yaml; the value found was %s", service.Name, service.Image)
			}
			// wraps it as an env variable
			appMan.Services[six].Image = fmt.Sprintf("${%s}", service.Image)
			// adds it to the manifest variable list
			appMan.Var.Items = append(appMan.Var.Items, AppVar{
				Name:        service.Image,
				Description: fmt.Sprintf("the container image name required by service %s", appMan.Services[six].Name),
				Value:       m.Spec.Images[service.Image],
				Secret:      false,
				Service:     strings.ToUpper(appMan.Services[six].Name),
			})
		}
		// spec bindings in scripts
		for _, init := range service.Info.Init {
			for _, script := range init.Scripts {
				i := service.Info.ScriptIx(script)
				// extract any bindings
				b := specBindings(appMan.Services[six].Info.Script[i].Content)
				for _, binding := range b {
					content := binding[len("${spec=") : len(binding)-1]
					parts := strings.Split(content, ":")
					switch len(parts) {
					case 2:
						if strings.HasPrefix(parts[0], "pkg") {
							pkg, exists := m.Spec.Packages[parts[1]]
							if !exists {
								return nil, fmt.Errorf("package key %s not found in spec version %s", parts[1], m.Spec.Version)
							}
							appMan.Services[six].Info.Script[i].Content = strings.Replace(appMan.Services[six].Info.Script[i].Content, binding, fmt.Sprintf("${%s}", parts[1]), 1)
							appMan.Var.Append(AppVar{
								Name:        parts[1],
								Description: "artisan package as defined in spec.yaml",
								Value:       pkg,
								Secret:      false,
								Service:     service.Name,
							})
						} else if strings.HasPrefix(parts[0], "img") {
							img, exists := m.Spec.Images[parts[1]]
							if !exists {
								return nil, fmt.Errorf("image key %s not found in spec version %s", parts[1], m.Spec.Version)
							}
							appMan.Services[six].Info.Script[i].Content = strings.Replace(appMan.Services[six].Info.Script[i].Content, binding, fmt.Sprintf("${%s}", parts[1]), 1)
							appMan.Var.Append(AppVar{
								Name:        parts[1],
								Description: "container image as defined in spec.yaml",
								Value:       img,
								Secret:      false,
								Service:     service.Name,
							})
						} else {
							return nil, fmt.Errorf("invalid spec spec binding '%s' in init script for service '%s'\n", binding, service.Name)
						}
					default:
						return nil, fmt.Errorf("invalid binding '%s' in init script for service '%s'\n", binding, service.Name)
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

// portKey extracts the port key from an expression like "port[http]"
// if the expression does not contain a map key then it returns the "default" key
func portKey(expression string) string {
	portkey := "default"
	if strings.HasPrefix(expression, "port[") {
		portkey = expression[len("port[") : len(expression)-1]
	}
	return portkey
}

// ensure one binding does not point to another so that the process of wiring variables is easier
func validateBindings(m Manifest, svc SvcRef) error {
	for _, v := range svc.Info.Var {
		// if the variable contains a binding expression
		if strings.Contains(v.Value, "${bind=") {
			// checks the target is not another binding
			parts := parseBinding(v.Value)
			if len(parts) == 3 && strings.ToLower(parts[1]) == "var" {
				svcName := parts[0]
				varName := parts[2]
				// find the target
				for _, service := range m.Services {
					if service.Name == svcName {
						for _, target := range service.Info.Var {
							if target.Name == varName && strings.Contains(target.Value, "${bind=") {
								return fmt.Errorf("a variable binding cannot point to another binding: in service %[1]s, "+
									"variable %[2]s=%[3]s points to service %s, variable %[4]s=%[5]s, which is a binding expression; "+
									"ensure the variable in service %[1]s points to a value, empty variable or a function expression", svc.Name, v.Name, v.Value, svcName, target.Name, target.Value)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func parseBinding(binding string) []string {
	value := binding[len("${bind=") : len(binding)-1]
	return strings.Split(value, ":")
}

func addDependency(dependsOn []string, svc string, s SvcRef) []string {
	result := make([]string, len(dependsOn))
	copy(result, dependsOn)
	exists := false
	for _, d := range result {
		if d == svc {
			exists = true
			break
		}
	}
	if !exists && s.Name != svc {
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
		// case of manifest in git repo defining scripts to run with artisan runtimes (hence no image)
		// so case to avoid is uri & image not provided
		if len(svc.Image) == 0 && len(svc.URI) == 0 {
			return fmt.Errorf("invalid entry for service '%s' manifest in application manifest: either one of Image or URI attributes must be specified\n", svc.Name)
		}
	}
	return nil
}

func (m *Manifest) deepCopy(dst interface{}) error {
	var buffer bytes.Buffer
	gob.Register(map[interface{}]interface{}{})
	if err := gob.NewEncoder(&buffer).Encode(m); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}

func (m *Manifest) getSvcTargetPort(svcName, portKey, binding string) (string, error) {
	for _, service := range m.Services {
		if service.Name == svcName {
			p, err := service.Info.PortMap()
			if err != nil {
				return "", err
			}
			targetPort, exists := p[portKey]
			if !exists {
				if portKey == "default" {
					return "", fmt.Errorf("port binding key not defined in service '%s': '%s'", svcName, binding)
				} else {
					return "", fmt.Errorf("port key '%s' not found in service '%s'", portKey, svcName)
				}
			}
			return fmt.Sprintf("%d", targetPort), nil
		}
	}
	return "", nil
}

func serviceBindings(value string) []string {
	r, _ := regexp.Compile("\\${bind=(?P<NAME>[^}]+)}")
	return r.FindAllString(value, -1)
}

func specBindings(value string) []string {
	r, _ := regexp.Compile("\\${spec=(?P<NAME>[^}]+)}")
	return r.FindAllString(value, -1)
}

// add credentials to http(s) URI
func addCredentialsToURI(uri string, creds string) (string, error) {
	// if there are no credentials or the uri is a file path
	if len(creds) == 0 || strings.HasPrefix(uri, "http") {
		// skip and return as is
		return uri, nil
	}
	parts := strings.Split(uri, "/")
	if !strings.HasPrefix(parts[0], "http") {
		return uri, fmt.Errorf("invalid URI scheme, http(s) expected when specifying credentials\n")
	}
	parts[2] = fmt.Sprintf("%s@%s", creds, parts[2])
	return strings.Join(parts, "/"), nil
}

func (m *Manifest) eval(t string) (string, error) {
	ctx := fileTempCtx{m: *m}
	tt, err := template.New("svc_file").Funcs(template.FuncMap{
		"service": ctx.serviceExists,
		"plus":    ctx.plus,
	}).Parse(t)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tt.Execute(&tpl, ctx)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (m *Manifest) loadSpec(uri string) (*release.Spec, error) {
	if !strings.Contains(uri, "/") {
		uri = fmt.Sprintf("./%s", uri)
	}
	root := uri[:strings.LastIndex(uri, "/")]
	specBytes, err := resx.ReadFile(fmt.Sprintf("%s/%s/%s", root, m.Version, "spec.yaml"), m.credentials)
	if err != nil {
		return nil, err
	}
	spec := new(release.Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		return nil, err
	}
	// check the spec version matches the version in the app manifest
	if spec.Version != m.Version {
		return nil, fmt.Errorf("version mismatch between app manifest and spec files: %s vs %s", m.Version, spec.Version)
	}
	return spec, nil
}

type fileTempCtx struct {
	m Manifest
}

// serviceExists checks if a service exists withing a file template
func (c *fileTempCtx) serviceExists(svcName reflect.Value) bool {
	for _, svc := range c.m.Services {
		if svc.Name == svcName.String() {
			return true
		}
	}
	return false
}

// plus adds two numbers
func (c *fileTempCtx) plus(number1, number2 reflect.Value) int64 {
	return number1.Int() + number2.Int()
}
