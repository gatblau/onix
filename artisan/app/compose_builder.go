/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"github.com/gatblau/onix/artisan/app/behaviour"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ComposeBuilder struct {
	manifest Manifest
}

// newComposeBuilder called internally by NewBuilder()
func newComposeBuilder(appMan Manifest) Builder {
	return &ComposeBuilder{manifest: appMan}
}

func (b *ComposeBuilder) Build() ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	composeProject, err := b.buildProject()
	if err != nil {
		return nil, err
	}
	rsx = append(rsx, *composeProject, b.buildEnv())
	files, err := b.buildFiles()
	if err != nil {
		return nil, err
	}
	svcScripts, err := b.buildInit()
	if err != nil {
		return nil, err
	}
	rsx = append(rsx, svcScripts...)
	rsx = append(rsx, files...)
	deployScript, err := b.buildDeploy()
	if err != nil {
		return nil, err
	}
	disposeScript, err := b.buildDispose()
	if err != nil {
		return nil, err
	}
	buildFile, err := b.buildFile()
	if err != nil {
		return nil, err
	}
	rsx = append(rsx, deployScript, disposeScript, buildFile)
	return rsx, nil
}

func (b *ComposeBuilder) buildProject() (*DeploymentRsx, error) {
	p := new(types.Project)
	p.Name = fmt.Sprintf("Docker Compose Project for %s", strings.ToUpper(b.manifest.Name))
	for _, svc := range b.manifest.Services {
		publishedPort, err := strconv.Atoi(svc.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid published port '%s'\n", svc.Port)
		}
		targetPort, err := strconv.Atoi(svc.Info.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid target port '%s'\n", svc.Port)
		}
		var ports []types.ServicePortConfig
		// if the public behaviour is defined then add a port mapping to the service
		if _, exists := svc.Is[behaviour.Public]; exists {
			ports = []types.ServicePortConfig{{Target: uint32(targetPort), Published: uint32(publishedPort)}}
		}
		if _, exists := svc.Is[behaviour.EncryptedInTransit]; exists {
			core.WarningLogger.Printf("service '%s' requested encryption of data in transit; it is currently not supported by the compose builder; skipping behaviour\n", svc.Name)
		}
		s := types.ServiceConfig{
			Name:          svc.Name,
			ContainerName: svc.Name,
			DependsOn:     getDeps(svc.DependsOn),
			Environment:   getEnv(svc.Info.Var),
			Image:         svc.Image,
			Ports:         ports,
			Restart:       "always",
			Volumes:       append(getSvcVols(svc.Info.Volume), getFileVols(svc.Info.File)...),
		}
		// if the load_balanced behaviour is defined then add replicated deployment mode to the service
		if replicas, exists := svc.Is[behaviour.LoadBalanced]; exists {
			rep, err2 := strconv.ParseUint(replicas, 10, 64)
			if err2 != nil {
				core.WarningLogger.Printf("failed to read load_balanced behaviour value '%s': %s\n", replicas, err2)
			} else {
				s.Deploy = &types.DeployConfig{
					Mode:     "replicated",
					Replicas: &rep,
				}
			}
		}
		p.Services = append(p.Services, s)
	}
	p.Volumes = getVols(b.manifest.Services)
	p.Networks = types.Networks{
		"default": types.NetworkConfig{
			Name: b.network(),
		},
	}
	composeProject, err := yaml.Marshal(p)
	if err != nil {
		return nil, err
	}
	composeProject = append([]byte("version: '3'\n\n"), composeProject...)
	return &DeploymentRsx{
		Name:    "docker-compose.yml",
		Content: composeProject,
		Type:    ComposeProject,
	}, nil
}

func (b *ComposeBuilder) network() string {
	return fmt.Sprintf("%s_network", strings.Replace(strings.ToLower(b.manifest.Name), " ", "_", -1))
}

func (b *ComposeBuilder) buildEnv() DeploymentRsx {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# %s application environment file for docker-compose.yml\n", strings.ToUpper(b.manifest.Name)))
	builder.WriteString(fmt.Sprintf("# auto-generated by Onix Artisan on %s\n\n", time.Now().UTC()))
	sort.Slice(b.manifest.Var.Items, func(i, j int) bool {
		return b.manifest.Var.Items[i].Service < b.manifest.Var.Items[j].Service
	})
	currentSvc := ""
	for _, v := range b.manifest.Var.Items {
		if v.Service != currentSvc {
			builder.WriteString("# -----------------------------------------------------------------\n")
			builder.WriteString(fmt.Sprintf("# %s service\n", strings.ToUpper(v.Service)))
			builder.WriteString("# -----------------------------------------------------------------\n\n")
			currentSvc = v.Service
		}
		builder.WriteString(fmt.Sprintf("# %s \n", v.Description))
		builder.WriteString(fmt.Sprintf("%s=%s\n\n", unwrap(v.Name), v.Value))
	}
	return DeploymentRsx{
		Name:    ".env",
		Content: []byte(builder.String()),
		Type:    EnvironmentFile,
	}
}

func unwrap(variable string) string {
	if strings.HasPrefix(variable, "${") && strings.HasSuffix(variable, "}") {
		return variable[2 : len(variable)-1]
	}
	return variable
}

func (b ComposeBuilder) buildFiles() ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	for _, svc := range b.manifest.Services {
		for _, f := range svc.Info.File {
			if len(f.Content) > 0 {
				rsx = append(rsx, DeploymentRsx{
					Name:    f.Path,
					Content: []byte(f.Content),
					Type:    ConfigurationFile,
				})
			} else {
				return nil, fmt.Errorf("definition of file '%s' in '%s' service manifest has no content\n", f.Path, svc.Name)
			}
		}
	}
	return rsx, nil
}

func (b ComposeBuilder) buildInit() ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	for _, svc := range b.manifest.Services {
		// if there is database schema configuration for the service
		if svc.Info.Db != nil {
			dbHeader := newHeaderBuilder("initialise database for %s service", svc.Name).String()
			rsx = append(rsx, DeploymentRsx{
				Name:    fmt.Sprintf("%s.sh", dbInitScriptName(svc)),
				Content: append([]byte(dbHeader), getDbScript(*svc.Info.Db)...),
				Type:    DbInitScript,
			})
		}
		// if there is specific initialisation logic for the service
		for _, init := range svc.Info.Init {
			if strings.ToLower(init.Builder) == "compose" {
				for _, script := range init.Scripts {
					i := svc.Info.ScriptIx(script)
					s := svc.Info.Script[i]
					scriptHeader := newHeaderBuilder("%s: %s", svc.Name, s.Description).String()
					rsx = append(rsx, DeploymentRsx{
						Name:    fmt.Sprintf("%s.sh", s.Name),
						Content: []byte(fmt.Sprintf("%s\n%s", scriptHeader, s.Content)),
						Type:    SvcInitScript,
					})
				}
			}
		}
	}
	return rsx, nil
}

func dbInitScriptName(svc SvcRef) string {
	return fmt.Sprintf("setup_db_%s", svc.Info.Db.Name)
}

func (b ComposeBuilder) buildDeploy() (DeploymentRsx, error) {
	header := newHeaderBuilder("application '%s' deploy script using docker-compose", b.manifest.Name)
	s := new(strings.Builder)
	s.WriteString(header.String())
	s.WriteString(fmt.Sprintf(`if ! command -v docker &> /dev/null; then
	echo "docker is required but not installed"
	exit
fi
if ! command -v docker-compose &> /dev/null; then
	echo "docker-compose is required but not installed"
	exit
fi
`))
	s.WriteString(fmt.Sprintf(`
# ensure attachable docker network is already created
if [[ $(docker network inspect %[1]s) == "[]" ]]; then
	echo Creating Docker network %[1]s ...
	docker network create %[1]s
fi
`, b.network()))
	s.WriteString("\n# create docker volumes\n")
	for _, service := range b.manifest.Services {
		for _, volume := range service.Info.Volume {
			s.WriteString(fmt.Sprintf("docker volume create %s\n", volume.Name))
		}
	}
	s.WriteString(fmt.Sprintf(`
# launch docker containers
docker-compose up -d
`))
	return DeploymentRsx{
		Name:    "deploy.sh",
		Content: []byte(s.String()),
		Type:    DeployScript,
	}, nil
}

func (b ComposeBuilder) buildDispose() (DeploymentRsx, error) {
	header := newHeaderBuilder("application '%s' dispose script using docker-compose", b.manifest.Name)
	s := new(strings.Builder)
	s.WriteString(header.String())
	s.WriteString(fmt.Sprintf(`
# bring down services
docker-compose down
`))
	s.WriteString("\n# remove docker volumes\n")
	for _, service := range b.manifest.Services {
		for _, volume := range service.Info.Volume {
			s.WriteString(fmt.Sprintf("docker volume rm %s\n", volume.Name))
		}
	}
	return DeploymentRsx{
		Name:    "dispose.sh",
		Content: []byte(s.String()),
		Type:    DeployScript,
	}, nil
}

func (b ComposeBuilder) buildFile() (DeploymentRsx, error) {
	buildFile := new(data.BuildFile)
	deploy := []string{"sh deploy.sh"}
	for svcix, service := range b.manifest.Services {
		if service.Info.Db != nil {
			buildFile.Functions = append(buildFile.Functions, &data.Function{
				Name:    dbInitScriptName(service),
				Run:     []string{fmt.Sprintf("sh %s.sh", dbInitScriptName(service))},
				Runtime: "dbman",
			})
			deploy = append(deploy, fmt.Sprintf("art runc -n %s %s", b.network(), dbInitScriptName(service)))
		}
		for _, init := range service.Info.Init {
			for _, script := range init.Scripts {
				i := service.Info.ScriptIx(script)
				runtime := b.manifest.Services[svcix].Info.Script[i].Runtime
				if len(runtime) > 0 {
					deploy = append(deploy, fmt.Sprintf("art runc -n %s %s", b.network(), script))
				} else {
					deploy = append(deploy, fmt.Sprintf("art run %s", script))
				}
				buildFile.Functions = append(buildFile.Functions, &data.Function{
					Name:    script,
					Run:     []string{fmt.Sprintf("sh %s.sh", script)},
					Runtime: runtime,
				})
			}
		}
	}
	export := true
	buildFile.Functions = append(buildFile.Functions, &data.Function{
		Name:        "deploy",
		Description: fmt.Sprintf("deploys the %s application using docker-compose", b.manifest.Name),
		Run:         deploy,
		Export:      &export,
	})
	buildFile.Functions = append(buildFile.Functions, &data.Function{
		Name:        "dispose",
		Description: fmt.Sprintf("disposes of all resources for the %s application", b.manifest.Name),
		Run:         []string{"sh dispose.sh"},
		Export:      &export,
	})
	content, err := yaml.Marshal(buildFile)
	return DeploymentRsx{
		Name:    "build.yaml",
		Content: content,
		Type:    BuildFile,
	}, err
}

func getDbScript(db Db) []byte {
	s := new(strings.Builder)
	s.WriteString(fmt.Sprintf("# configure '%s' database release information\n", db.Name))
	s.WriteString(fmt.Sprintf("dbman config use -n %s-config\n", db.Name))
	s.WriteString(fmt.Sprintf("dbman config set repo.uri %s\n", db.SchemaURI))
	s.WriteString(fmt.Sprintf("dbman config set db.provider %s\n", db.Provider))
	s.WriteString(fmt.Sprintf("dbman config set db.host %s\n", db.Host))
	s.WriteString(fmt.Sprintf("dbman config set db.port %d\n", db.Port))
	s.WriteString(fmt.Sprintf("dbman config set db.name %s\n", db.Name))
	s.WriteString(fmt.Sprintf("dbman config set db.username %s\n", db.User))
	s.WriteString(fmt.Sprintf("dbman config set db.password %s\n", db.Pwd))
	s.WriteString(fmt.Sprintf("dbman config set db.adminusername %s\n", db.AdminUser))
	s.WriteString(fmt.Sprintf("dbman config set db.adminpassword %s\n", db.AdminPwd))
	s.WriteString(fmt.Sprintf("dbman config set appversion %s\n\n", db.AppVersion))
	s.WriteString(fmt.Sprintf("# create '%s' database\n", db.Name))
	s.WriteString(fmt.Sprintf("dbman db create\n\n"))
	s.WriteString(fmt.Sprintf("# deploy '%s' database schema\n", db.Name))
	s.WriteString(fmt.Sprintf("dbman db deploy\n\n"))
	return []byte(s.String())
}

func getSvcVols(volume []Volume) []types.ServiceVolumeConfig {
	vo := make([]types.ServiceVolumeConfig, 0)
	// does any explicit volumes
	for _, v := range volume {
		vo = append(vo, types.ServiceVolumeConfig{
			Source: v.Name,
			Target: v.Path,
			Type:   "volume",
		})
	}
	return vo
}

// gets a list of volumes required by the specified files
func getFileVols(files []File) []types.ServiceVolumeConfig {
	vo := make([]types.ServiceVolumeConfig, 0)
	// does any explicit volumes
	for _, f := range files {
		relD := relDir(f.Path)
		found := false
		for _, x := range vo {
			if x.Source == relD {
				found = true
			}
		}
		if !found {
			vo = append(vo, types.ServiceVolumeConfig{
				Source: relD,
				Target: absDir(f.Path),
				Type:   "bind",
			})
		}
	}
	return vo
}

func relDir(path string) string {
	// if the path is absolute
	if path[0] == '/' {
		// returns a relative form
		return fmt.Sprintf("./%s", filepath.Dir(path[1:]))
	}
	// if the path is not absolute but does not start with ./ add it
	if path[0:1] != "./" {
		return fmt.Sprintf("./%s", filepath.Dir(path[1:]))
	}
	// otherwise, return as is
	return filepath.Dir(path)
}

func absDir(path string) string {
	if path[0] == '/' {
		return filepath.Dir(path)
	}
	return filepath.Dir(fmt.Sprintf("/%s", filepath.Dir(path)))
}

func getDeps(dependencies []string) types.DependsOnConfig {
	d := types.DependsOnConfig{}
	for _, dependency := range dependencies {
		d[dependency] = types.ServiceDependency{Condition: types.ServiceConditionStarted}
	}
	return d
}

func newHeaderBuilder(label string, a ...interface{}) *strings.Builder {
	mergedLabel := fmt.Sprintf(label, a...)
	script := &strings.Builder{}
	script.WriteString("#!/bin/bash\n")
	script.WriteString(fmt.Sprintf("# %s\n# auto-generated by Artisan on %s\n", mergedLabel, time.Now().UTC()))
	lines := strings.Split(script.String(), "\n")
	script = &strings.Builder{}
	for _, line := range lines {
		if len(line) > 0 && line[0] != '#' {
			script.WriteString("# ")
		}
		script.WriteString(line)
		script.WriteString("\n")
	}
	return script
}

func getEnv(vars []Var) types.MappingWithEquals {
	var values []string
	for _, v := range vars {
		values = append(values, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	return types.NewMappingWithEquals(values)
}

func getVols(svc []SvcRef) types.Volumes {
	vo := types.Volumes{}
	for _, s := range svc {
		for _, v := range s.Info.Volume {
			vo[v.Name] = types.VolumeConfig{
				External: types.External{External: true},
			}
		}
	}
	return vo
}
