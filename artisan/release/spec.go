/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package release

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gatblau/onix/artisan/data"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/oxlib/resx"
	"gopkg.in/yaml.v2"
)

// Spec the specification for artisan artefacts to be exported
type Spec struct {
	// Name of the released Application
	Name string `yaml:"name,omitempty"`
	// Description of the Application
	Description string `yaml:"description,omitempty"`
	// Author of the release
	Author string `yaml:"author,omitempty"`
	// License associated to the release
	License string `yaml:"license,omitempty"`
	// Version of the release
	Version string `yaml:"version"`
	// Info general release information
	Info string `yaml:"info,omitempty"`
	// Images the container images in the release
	Images map[string]string `yaml:"images,omitempty"`
	// Packages the artisan packages in the release
	Packages map[string]string `yaml:"packages,omitempty"`
	// OsPackages operating system packages that are part of the release
	OsPackages map[string]map[string]string `yaml:"os_packages,omitempty"`
	// Run commands
	Run []Run

	content []byte
}

// Run defines one or more package/function to run in specific cases
type Run struct {
	// Package name to run
	Package string `yaml:"package"`
	// Function in the package to run
	Function string `yaml:"function"`
	// Var list of variables to be passed to the function to run
	Input *data.Input `yaml:"input,omitempty"`
	// Event the lifecycle event that is associated to the function to run
	Event string `yaml:"event"`
}

// LifecycleEvent the event in the lifecycle of the spec that
type LifecycleEvent string

const (
	ReleaseSetup  string = "SETUP"
	ReleaseDeploy        = "DEPLOY"
	ReleaseDecom         = "DECOMMISSION"
)

func (e LifecycleEvent) String() string {
	extensions := [...]string{"SETUP", "DEPLOY", "DECOMMISSION"}
	x := string(e)
	for _, v := range extensions {
		if v == x {
			return x
		}
	}
	return ""
}

func NewSpec(path, creds string) (*Spec, error) {
	var (
		content []byte
		err     error
	)
	// if the path does not contain a spec file
	if !strings.HasSuffix(path, "spec.yaml") {
		if strings.HasSuffix(path, "yaml") || strings.HasSuffix(path, "yml") || strings.HasSuffix(path, "txt") || strings.HasSuffix(path, "json") {
			return nil, fmt.Errorf("invalid spec file, it should be spec.yaml")
		}
		path = fmt.Sprintf("%s/spec.yaml", path)
	}
	// if path contains scheme it is remote
	if strings.Contains(path, "://") {
		content, err = resx.ReadFile(path, creds)
		if err != nil {
			return nil, fmt.Errorf("cannot read remote spec file '%s': %s", path, err)
		}
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path: %s", err)
		}
		content, err = resx.ReadFile(path, creds)
		if err != nil {
			return nil, fmt.Errorf("cannot read spec file %s: %s", path, err)
		}
	}
	spec := new(Spec)
	err = yaml.Unmarshal(content, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal spec file: %s", err)
	}
	spec.content = content

	return spec, spec.Valid()
}

func ExportSpec(opts ExportOptions) error {
	if err := opts.Valid(); err != nil {
		return fmt.Errorf("invalid export options: %s\n", err)
	}
	var skipArtefact bool
	// save packages first
	l := registry.NewLocalRegistry(opts.ArtHome)
	for _, value := range opts.Specification.Packages {
		if skipArtefact, opts.Filter = skip(opts.Filter, value); skipArtefact {
			if len(opts.Filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", opts.Filter)
			}
			core.InfoLogger.Printf("skipping package %s\n", value)
			continue
		}
		name, err := core.ParseName(value)
		if err != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		uri := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(value))
		err = l.ExportPackage([]core.PackageName{*name}, opts.SourceCreds, uri, opts.TargetCreds)
		if err != nil {
			return fmt.Errorf("cannot save package %s: %s", value, err)
		}
	}
	// save images
	for _, value := range opts.Specification.Images {
		if skipArtefact, opts.Filter = skip(opts.Filter, value); skipArtefact {
			if len(opts.Filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", opts.Filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", value)
			continue
		}
		// note: the package is saved with a name exactly the same as the container image
		// to avoid the art package name parsing from failing, any images with no host or user/group in the name should be avoided
		// e.g. docker.io/mongo-express:latest will fail so use docker.io/library/mongo-express:latest instead
		err := ExportImage(value, value, opts.TargetUri, opts.TargetCreds, opts.ArtHome)
		if err != nil {
			return fmt.Errorf("cannot save image %s: %s", value, err)
		}
	}

	// download linux packages
	for key, value := range opts.Specification.OsPackages {
		if len(value) == 0 {
			return fmt.Errorf("missing package names in the spec file for type %s", value)
		}
		if strings.ToLower(key) == "apt" {
			cmd := "apt-get -v"
			res, err := build.Exe(cmd, opts.ArtHome, merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return fmt.Errorf("failed to get the apt-get version number %s", err)
			}
			re := regexp.MustCompile(`\d+\.(\d)*`)
			if len(res) == 0 || len(re.FindString(res)) == 0 {
				return fmt.Errorf("the host is not a debian distribution or apt-get package is missing")
			}

			core.InfoLogger.Printf("performing %s exporting \n", key)
			// collect all package names into slice
			var pkges []string
			for _, v := range value {
				if skipArtefact, opts.Filter = skip(opts.Filter, v); skipArtefact {
					if len(opts.Filter) == 0 {
						core.WarningLogger.Printf("invalid filter expression for Os_Packages '%s'\n", opts.Filter)
					}
					core.InfoLogger.Printf("skipping Os_Package filtering %s\n", value)
					continue
				}
				pkges = append(pkges, v)
			}
			//export all debian packages into a single artisan package
			err = ExportDebianPackage(pkges, opts)
			if err != nil {
				return fmt.Errorf("failed to export debian package %s: %s", value, err)
			}
		} else if strings.ToLower(key) == "rpm" {
			return fmt.Errorf("rpm packages are currently not support")
		} else {
			return fmt.Errorf("%s is invalid packaging option, valid values are apt, rpm", key)
		}
	}
	// finally, save the spec to the target location
	// note: this is done last so that a minio notification can be triggered based on this file
	// once all other artefacts have been exported
	uri := fmt.Sprintf("%s/spec.yaml", opts.TargetUri)
	err := resx.WriteFile(opts.Specification.content, uri, opts.TargetCreds)
	if err != nil {
		return fmt.Errorf("cannot save spec file: %s", err)
	}
	core.InfoLogger.Printf("writing spec.yaml to %s", opts.TargetUri)
	return nil
}

func ImportSpec(opts ImportOptions) (*Spec, error) {
	if err := opts.Valid(); err != nil {
		return nil, fmt.Errorf("invalid import options: %s\n", err)
	}
	var skipArtefact bool
	r := registry.NewLocalRegistry(opts.ArtHome)
	uri := fmt.Sprintf("%s/spec.yaml", opts.TargetUri)
	core.InfoLogger.Printf("retrieving %s\n", uri)
	specBytes, err := resx.ReadFile(uri, opts.TargetCreds)
	if err != nil {
		return nil, fmt.Errorf("cannot read spec.yaml: %s", err)
	}
	spec := new(Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal spec.yaml: %s", err)
	}
	// import packages
	for _, pkName := range spec.Packages {
		if skipArtefact, opts.Filter = skip(opts.Filter, pkName); skipArtefact {
			if len(opts.Filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", opts.Filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", pkName)
			continue
		}
		name := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(pkName))
		err2 := r.Import([]string{name}, opts.TargetCreds, opts.VProc)
		if err2 != nil {
			return spec, fmt.Errorf("cannot read %s.tar: %s", pkgName(pkName), err2)
		}
	}
	// import images
	for _, image := range spec.Images {
		if skipArtefact, opts.Filter = skip(opts.Filter, image); skipArtefact {
			if len(opts.Filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", opts.Filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", image)
			continue
		}
		name := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(image))
		err2 := r.Import([]string{name}, opts.TargetCreds, opts.VProc)
		if err2 != nil {
			return spec, fmt.Errorf("cannot read %s.tar: %s", pkgName(image), err)
		}
		core.InfoLogger.Printf("loading => %s\n", image)
		_, err2 = build.Exe(fmt.Sprintf("art exe %s import", image), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err2 != nil {
			return spec, fmt.Errorf("cannot import image %s: %s", image, err2)
		}
	}
	return spec, nil
}

func DownloadSpec(opts UpDownOptions) (*Spec, error) {
	if err := opts.Valid(); err != nil {
		return nil, fmt.Errorf("invalid download options: %s\n", err)
	}
	spec, err := NewSpec(opts.TargetUri, opts.TargetCreds)
	if err != nil {
		return nil, fmt.Errorf("cannot load specification: %s", err)
	}
	if err = checkPath(opts.LocalPath); err != nil {
		return nil, fmt.Errorf("cannot create local path: %s", err)
	}
	err = os.WriteFile(filepath.Join(opts.LocalPath, "spec.yaml"), spec.content, 0755)
	if err != nil {
		return nil, err
	}
	for _, pkg := range spec.Packages {
		pkgUri := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(pkg))
		core.InfoLogger.Printf("downloading => %s\n", pkgUri)
		tarBytes, err2 := resx.ReadFile(pkgUri, opts.TargetCreds)
		if err2 != nil {
			return nil, err2
		}
		pkgPath := filepath.Join(opts.LocalPath, filepath.Base(pkgUri))
		core.InfoLogger.Printf("writing => %s\n", pkgPath)
		err = os.WriteFile(pkgPath, tarBytes, 0755)
		if err != nil {
			return nil, err
		}
	}
	for _, image := range spec.Images {
		imageUri := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(image))
		core.InfoLogger.Printf("downloading => %s\n", imageUri)
		tarBytes, err2 := resx.ReadFile(imageUri, opts.TargetCreds)
		if err2 != nil {
			return nil, err2
		}
		targetFile := filepath.Join(opts.LocalPath, fmt.Sprintf("%s.tar", pkgName(image)))
		core.InfoLogger.Printf("writing => %s\n", targetFile)
		err = os.WriteFile(targetFile, tarBytes, 0755)
		if err != nil {
			return nil, err
		}
	}
	return spec, nil
}

func UploadSpec(opts UpDownOptions) error {
	if err := opts.Valid(); err != nil {
		return fmt.Errorf("invalid upload options: %s\n", err)
	}
	spec, err := NewSpec(opts.LocalPath, opts.TargetCreds)
	if err != nil {
		return fmt.Errorf("cannot load specification: %s", err)
	}
	if err = checkPath(opts.LocalPath); err != nil {
		return fmt.Errorf("cannot create local path: %s", err)
	}
	err = resx.WriteFile(spec.content, fmt.Sprintf("%s/spec.yaml", opts.TargetUri), opts.TargetCreds)
	if err != nil {
		return err
	}
	for _, pkg := range spec.Packages {
		localUri := fmt.Sprintf("%s/%s.tar", opts.LocalPath, pkgName(pkg))
		remoteUri := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(pkg))
		core.InfoLogger.Printf("uploading => %s\n", remoteUri)
		content, readErr := os.ReadFile(localUri)
		if readErr != nil {
			return readErr
		}
		err = resx.WriteFile(content, remoteUri, opts.TargetCreds)
		if err != nil {
			return err
		}
	}
	for _, image := range spec.Images {
		localUri := fmt.Sprintf("%s/%s.tar", opts.LocalPath, pkgName(image))
		remoteUri := fmt.Sprintf("%s/%s.tar", opts.TargetUri, pkgName(image))
		core.InfoLogger.Printf("uploading => %s\n", remoteUri)
		content, readErr := os.ReadFile(localUri)
		if readErr != nil {
			return readErr
		}
		err = resx.WriteFile(content, remoteUri, opts.TargetCreds)
		if err != nil {
			return err
		}
	}
	return nil
}

func PullSpec(opts PullOptions) error {
	if err := opts.Valid(); err != nil {
		return fmt.Errorf("invalid pull options: %s\n", err)
	}
	cli, cmdErr := containerCmd()
	if cmdErr != nil {
		return cmdErr
	}
	local := registry.NewLocalRegistry(opts.ArtHome)
	spec, err := NewSpec(opts.TargetUri, opts.TargetCreds)
	if err != nil {
		return fmt.Errorf("cannot load specification: %s", err)
	}
	for _, pkg := range spec.Packages {
		p, parseErr := core.ParseName(pkg)
		if parseErr != nil {
			return parseErr
		}
		core.InfoLogger.Printf("pulling => %s\n", pkg)
		local.Pull(p, opts.SourceCreds, false)
	}
	for _, image := range spec.Images {
		core.InfoLogger.Printf("pulling => %s\n", image)
		_, err = build.Exe(fmt.Sprintf("%s pull %s", cli, image), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func PushSpec(opts PushOptions) error {
	if err := opts.Valid(); err != nil {
		return fmt.Errorf("invalid push options: %s\n", err)
	}
	var (
		cli, usr, pwd string
		err           error
	)
	// if it needs to work with container images
	if opts.Image {
		// obtain docker command name
		cli, err = containerCmd()
		if err != nil {
			return err
		}
		// if credentials have been provided to connect to registry
		if len(opts.User) > 0 {
			core.InfoLogger.Printf("logging to docker registry")
			// executes docker login --username=right-username --password=""
			usr, pwd = core.UserPwd(opts.User)
			out, eErr := build.Exe(fmt.Sprintf("%s login %s --username=%s --password=%s", cli, opts.Host, usr, pwd), ".", merge.NewEnVarFromSlice([]string{}), false)
			if eErr != nil {
				// do not return original error as it can contain sensitive info (e.g. password)
				return fmt.Errorf("docker login failed")
			}
			core.InfoLogger.Printf("%s\n", out)
		}
	}
	local := registry.NewLocalRegistry(opts.ArtHome)
	spec, err := NewSpec(opts.SpecPath, opts.Creds)
	if err != nil {
		return fmt.Errorf("cannot load spec.yaml: %s", err)
	}
	if !opts.Image {
		for _, pac := range spec.Packages {
			tgtNameStr, tgtName, tgtNameErr := targetName(pac, opts.Group, opts.Host)
			if tgtNameErr != nil {
				return fmt.Errorf("cannot work out target name: %s", tgtNameErr)
			}
			// tag the package with the target registry name
			core.InfoLogger.Printf("tagging => '%s' to '%s'\n", pac, tgtNameStr)
			err = local.Tag(pac, tgtNameStr)
			if err != nil {
				return err
			}
			// push to remote
			core.InfoLogger.Printf("pushing => '%s'\n", tgtNameStr)
			err = local.Push(tgtName, opts.User, false)
			if err != nil {
				return err
			}
			// if cleaning has been specified
			if opts.Clean {
				// remove the package from the local package registry
				core.InfoLogger.Printf("removing => '%s'\n", tgtNameStr)
				err = local.Remove([]string{tgtNameStr})
				if err != nil {
					return err
				}
			}
		}
	} else {
		for _, img := range spec.Images {
			tgtNameStr, _, tgtNameErr := targetName(img, opts.Group, opts.Host)
			if tgtNameErr != nil {
				return fmt.Errorf("cannot work out target name: %s", tgtNameErr)
			}
			// tag the package with the target registry name
			core.InfoLogger.Printf("tagging => '%s' to '%s'\n", img, tgtNameStr)
			// docker tag image
			_, err = build.Exe(fmt.Sprintf("%s tag %s %s", cli, img, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return err
			}
			// push to remote
			core.InfoLogger.Printf("pushing => '%s'\n", tgtNameStr)
			_, err = build.Exe(fmt.Sprintf("%s push %s", cli, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return err
			}
			// if cleaning has been specified
			if opts.Clean {
				// remove the image from the local container registry
				core.InfoLogger.Printf("removing => '%s'\n", tgtNameStr)
				_, err = build.Exe(fmt.Sprintf("%s rmi --force %s", cli, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
				if err != nil {
					return err
				}
			}
		}
		// if a logout was requested
		if opts.Image && len(opts.User) > 0 && opts.Logout {
			_, err = build.Exe(fmt.Sprintf("%s logout %s", cli, opts.Host), ".", merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func targetName(name string, group string, host string) (string, *core.PackageName, error) {
	srcName, parseErr := core.ParseName(name)
	if parseErr != nil {
		return "", nil, parseErr
	}
	var tName string
	// if a target group has been defined
	if len(group) > 0 {
		// build a target name using host and group but keeping source package name and tag
		tName = fmt.Sprintf("%s/%s/%s:%s", host, group, srcName.Name, srcName.Tag)
	} else {
		// if a target group has not been specified, then
		// build a target name using host but keeping source package group, name and tag
		tName = fmt.Sprintf("%s/%s/%s:%s", host, srcName.Group, srcName.Name, srcName.Tag)
	}
	tgtName, parseTargetErr := core.ParseName(tName)
	if parseTargetErr != nil {
		return "", nil, parseTargetErr
	}
	return tName, tgtName, nil
}

func (s *Spec) ContainsImage(name string) bool {
	for key, _ := range s.Images {
		if name == key {
			return true
		}
	}
	return false
}

func pkgName(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name, "/", "_"), ".", "_"), "-", "_")
}

func checkPath(path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// if the path does not exist
	if _, err = os.Stat(p); os.IsNotExist(err) {
		// creates it
		err = os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func skip(filter, value string) (bool, string) {
	// if there is a filter defined
	if len(filter) > 0 {
		matched, err := regexp.MatchString(filter, value)
		if err != nil {
			return false, ""
		}
		return !matched, filter
	}
	return false, filter
}

func (s *Spec) Valid() error {

	if len(s.Name) == 0 || len(s.Description) == 0 || len(s.Author) == 0 ||
		len(s.License) == 0 || len(s.Version) == 0 || len(s.Info) == 0 {

		return fmt.Errorf(" spec file must have 'name', 'description', 'author', 'license', 'version', 'info'")
	}

	if s.Packages == nil && s.Images == nil && s.OsPackages == nil {
		return fmt.Errorf(" spec file has no details of  'packages' or 'images' or 'os_packages'")
	}

	if len(s.Images) > 0 {
		invalidImgs := []string{}
		for _, v := range s.Images {
			if !(len(v) > 0 && len(strings.Split(v, "/")) > 2) {
				invalidImgs = append(invalidImgs, v)
			}
		}
		if len(invalidImgs) > 0 {
			return fmt.Errorf(" invalid format of container image name for following images [ %s ]"+
				"\n valid format is <host|domain>/<group>/<image-name> ", strings.Join(invalidImgs, ", "))
		}
	}

	return nil
}
