/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package export

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/oxlib/resx"
	"gopkg.in/yaml.v2"
)

// Spec the specification for artisan artefacts to be exported
type Spec struct {
	Version  string            `yaml:"version"`
	Info     string            `yaml:"info,omitempty"`
	Images   map[string]string `yaml:"images,omitempty"`
	Packages map[string]string `yaml:"packages,omitempty"`

	content []byte
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
			return nil, fmt.Errorf("cannot read remote spec file: %s", err)
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
	return spec, nil
}

func ExportSpec(s Spec, targetUri, sourceCreds, targetCreds, filter string) error {
	var skipArtefact bool
	// save packages first
	l := registry.NewLocalRegistry()
	for _, value := range s.Packages {
		if skipArtefact, filter = skip(filter, value); skipArtefact {
			if len(filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", filter)
			}
			core.InfoLogger.Printf("skipping package %s\n", value)
			continue
		}
		name, err := core.ParseName(value)
		if err != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		uri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(value))
		err = l.ExportPackage([]core.PackageName{*name}, sourceCreds, uri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save package %s: %s", value, err)
		}
	}
	// save images
	for _, value := range s.Images {
		if skipArtefact, filter = skip(filter, value); skipArtefact {
			if len(filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", value)
			continue
		}
		// note: the package is saved with a name exactly the same as the container image
		// to avoid the art package name parsing from failing, any images with no host or user/group in the name should be avoided
		// e.g. docker.io/mongo-express:latest will fail so use docker.io/library/mongo-express:latest instead
		err := ExportImage(value, value, targetUri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save image %s: %s", value, err)
		}
	}
	// finally, save the spec to the target location
	// note: this is done last so that a minio notification can be triggered based on this file
	// once all other artefacts have been exported
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	err := resx.WriteFile(s.content, uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot save spec file: %s", err)
	}
	core.InfoLogger.Printf("writing spec.yaml to %s", targetUri)
	return nil
}

func ImportSpec(targetUri, targetCreds, filter, pubKeyPath string, ignoreSignature bool) error {
	// if it is not ignoring the package signature, then a public key path must be provided
	if !ignoreSignature && len(pubKeyPath) == 0 {
		return fmt.Errorf("the path to a public key must be provided to verify the package author, otherwise ignore signature")
	}
	var skipArtefact bool
	r := registry.NewLocalRegistry()
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	core.InfoLogger.Printf("retrieving %s\n", uri)
	specBytes, err := resx.ReadFile(uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot read spec.yaml: %s", err)
	}
	spec := new(Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		return fmt.Errorf("cannot unmarshal spec.yaml: %s", err)
	}
	// import packages
	for _, pkName := range spec.Packages {
		if skipArtefact, filter = skip(filter, pkName); skipArtefact {
			if len(filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", pkName)
			continue
		}
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkName))
		err2 := r.Import([]string{name}, targetCreds, pubKeyPath, ignoreSignature)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(pkName), err2)
		}
	}
	// import images
	for _, image := range spec.Images {
		if skipArtefact, filter = skip(filter, image); skipArtefact {
			if len(filter) == 0 {
				core.WarningLogger.Printf("invalid filter expression '%s'\n", filter)
			}
			core.InfoLogger.Printf("skipping image %s\n", image)
			continue
		}
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		err2 := r.Import([]string{name}, targetCreds, pubKeyPath, ignoreSignature)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(image), err)
		}
		core.InfoLogger.Printf("loading => %s\n", image)
		ignoreSigFlag := ""
		if ignoreSignature {
			ignoreSigFlag = "-s"
		}
		_, err2 = build.Exe(fmt.Sprintf("art exe %s import %s", image, ignoreSigFlag), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err2 != nil {
			return fmt.Errorf("cannot import image %s: %s", image, err2)
		}
	}
	return nil
}

func DownloadSpec(targetUri, targetCreds, localPath string) (*Spec, error) {
	spec, err := NewSpec(targetUri, targetCreds)
	if err != nil {
		return nil, fmt.Errorf("cannot load specification: %s", err)
	}
	if err = checkPath(localPath); err != nil {
		return nil, fmt.Errorf("cannot create local path: %s", err)
	}
	err = os.WriteFile(filepath.Join(localPath, "spec.yaml"), spec.content, 0755)
	if err != nil {
		return nil, err
	}
	for _, pkg := range spec.Packages {
		pkgUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkg))
		core.InfoLogger.Printf("downloading => %s\n", pkgUri)
		tarBytes, err2 := resx.ReadFile(pkgUri, targetCreds)
		if err2 != nil {
			return nil, err2
		}
		pkgPath := filepath.Join(localPath, filepath.Base(pkgUri))
		core.InfoLogger.Printf("writing => %s\n", pkgPath)
		err = os.WriteFile(pkgPath, tarBytes, 0755)
		if err != nil {
			return nil, err
		}
	}
	for _, image := range spec.Images {
		imageUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		core.InfoLogger.Printf("downloading => %s\n", imageUri)
		tarBytes, err2 := resx.ReadFile(imageUri, targetCreds)
		if err2 != nil {
			return nil, err2
		}
		targetFile := filepath.Join(localPath, fmt.Sprintf("%s.tar", pkgName(image)))
		core.InfoLogger.Printf("writing => %s\n", targetFile)
		err = os.WriteFile(targetFile, tarBytes, 0755)
		if err != nil {
			return nil, err
		}
	}
	return spec, nil
}

func UploadSpec(targetUri, targetCreds, localPath string) error {
	spec, err := NewSpec(localPath, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot load specification: %s", err)
	}
	if err = checkPath(localPath); err != nil {
		return fmt.Errorf("cannot create local path: %s", err)
	}
	err = resx.WriteFile(spec.content, fmt.Sprintf("%s/spec.yaml", targetUri), targetCreds)
	if err != nil {
		return err
	}
	for _, pkg := range spec.Packages {
		localUri := fmt.Sprintf("%s/%s.tar", localPath, pkgName(pkg))
		remoteUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkg))
		core.InfoLogger.Printf("uploading => %s\n", remoteUri)
		content, readErr := os.ReadFile(localUri)
		if readErr != nil {
			return readErr
		}
		err = resx.WriteFile(content, remoteUri, targetCreds)
		if err != nil {
			return err
		}
	}
	for _, image := range spec.Images {
		localUri := fmt.Sprintf("%s/%s.tar", localPath, pkgName(image))
		remoteUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		core.InfoLogger.Printf("uploading => %s\n", remoteUri)
		content, readErr := os.ReadFile(localUri)
		if readErr != nil {
			return readErr
		}
		err = resx.WriteFile(content, remoteUri, targetCreds)
		if err != nil {
			return err
		}
	}
	return nil
}

func PullSpec(targetUri, targetCreds, sourceCreds string) error {
	cli, cmdErr := containerCmd()
	if cmdErr != nil {
		return cmdErr
	}
	local := registry.NewLocalRegistry()
	spec, err := NewSpec(targetUri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot load specification: %s", err)
	}
	for _, pkg := range spec.Packages {
		p, parseErr := core.ParseName(pkg)
		if parseErr != nil {
			return parseErr
		}
		core.InfoLogger.Printf("pulling => %s\n", pkg)
		local.Pull(p, sourceCreds)
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

func PushSpec(specPath, host, group, user, creds string, image, clean, logout bool) error {
	var (
		cli, usr, pwd string
		err           error
	)
	// if it needs to work with container images
	if image {
		// obtain docker command name
		cli, err = containerCmd()
		if err != nil {
			return err
		}
		// if credentials have been provided to connect to registry
		if len(user) > 0 {
			core.InfoLogger.Printf("logging to docker registry")
			// executes docker login --username=right-username --password=""
			usr, pwd = core.UserPwd(user)
			out, eErr := build.Exe(fmt.Sprintf("%s login %s --username=%s --password=%s", cli, host, usr, pwd), ".", merge.NewEnVarFromSlice([]string{}), false)
			if eErr != nil {
				return eErr
			}
			core.InfoLogger.Printf("%s\n", out)
		}
	}
	local := registry.NewLocalRegistry()
	spec, err := NewSpec(specPath, creds)
	if err != nil {
		return fmt.Errorf("cannot load spec.yaml: %s", err)
	}
	if !image {
		for _, pac := range spec.Packages {
			tgtNameStr, tgtName, tgtNameErr := targetName(pac, group, host)
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
			err = local.Push(tgtName, user)
			if err != nil {
				return err
			}
			// if cleaning has been specified
			if clean {
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
			tgtNameStr, _, tgtNameErr := targetName(img, group, host)
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
			if clean {
				// remove the image from the local container registry
				core.InfoLogger.Printf("removing => '%s'\n", tgtNameStr)
				_, err = build.Exe(fmt.Sprintf("%s rmi --force %s", cli, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
				if err != nil {
					return err
				}
			}
		}
		// if a logout was requested
		if image && len(user) > 0 && logout {
			_, err = build.Exe(fmt.Sprintf("%s logout %s", cli, host), ".", merge.NewEnVarFromSlice([]string{}), false)
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
