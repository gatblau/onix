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
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExportImage(imgName, packName, targetUri, creds, artHome string) error {
	pName, err := core.ParseName(packName)
	if err != nil {
		return fmt.Errorf("invalid package name: %s, ensure the container image name fully specify host and user/group if from docker.io", err)
	}
	// if a target has been specified
	if len(targetUri) > 0 {
		// if a final slash does not exist add it
		if targetUri[len(targetUri)-1] != '/' {
			targetUri = fmt.Sprintf("%s/", targetUri)
		}
		// automatically adds a tar filename to the URI based on the package name:tag
		targetUri = fmt.Sprintf("%s%s", targetUri, pkgFilename(*pName))
	} else {
		return fmt.Errorf("a destination URI must be specified to export the image")
	}
	// should we use docker or podman?
	containerCli, err := containerCmd()
	if err != nil {
		return fmt.Errorf("cannot create image archive: %s", err)
	}
	// create a build file to build the package containing the image tar
	pbf := data.BuildFile{
		Runtime: "ubi-min",
		Labels: map[string]string{
			"image": imgName,
		},
		Profiles: []*data.Profile{
			{
				Name:   "package-image",
				Target: "./build",
				Type:   "image",
			},
		},
	}
	pbfBytes, err := yaml.Marshal(pbf)
	if err != nil {
		return fmt.Errorf("cannot marshall packaging build file: %s", err)
	}
	// create a build file to import image tar in package
	export := true
	bf := data.BuildFile{
		Runtime: "ubi-min",
		Labels: map[string]string{
			"image": imgName,
		},
		Functions: []*data.Function{
			{
				Name:        "import",
				Description: "imports docker image in local docker registry",
				Export:      &export,
				Runtime:     "ubi-min",
				Run: []string{
					fmt.Sprintf("%s load -i %s.tar", containerCli, pkgName(imgName)),
				},
			},
		},
	}
	bfBytes, err := yaml.Marshal(bf)
	if err != nil {
		return fmt.Errorf("cannot marshall package build file: %s", err)
	}

	tmp, err := core.NewTempDir(artHome)
	if err != nil {
		return fmt.Errorf("cannot create temp folder for processing image archive: %s", err)
	}
	// create a target folder for the artisan package
	targetFolder := filepath.Join(tmp, "build")
	err = os.MkdirAll(targetFolder, 0755)
	// workout the docker save command
	cmd := fmt.Sprintf("%s save %s -o %s/%s.tar", containerCli, imgName, targetFolder, imgFilename(imgName))
	// execute the command synchronously
	core.InfoLogger.Printf("exporting image %s to tarball file", imgName)
	_, err = build.Exe(cmd, tmp, merge.NewEnVarFromSlice([]string{}), false)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot execute archive command: %s", err)
	}
	core.InfoLogger.Println("packaging image tarball file")
	err = os.WriteFile(filepath.Join(tmp, "build.yaml"), pbfBytes, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save packaging build file: %s", err)
	}
	err = os.WriteFile(filepath.Join(targetFolder, "build.yaml"), bfBytes, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save package build file: %s", err)
	}
	b := build.NewBuilder(artHome)
	b.Build(tmp, "", "", pName, "", false, false, "")
	r := registry.NewLocalRegistry(artHome)
	// export package
	core.InfoLogger.Printf("exporting image package to tarball file")
	err = r.ExportPackage([]core.PackageName{*pName}, "", targetUri, creds)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save package to destination: %s", err)
	}
	return nil
}

// return the command to run to launch a container
func containerCmd() (string, error) {
	if isCmdAvailable("docker") {
		return "docker", nil
	} else if isCmdAvailable("podman") {
		return "podman", nil
	}
	return "", fmt.Errorf("either podman or docker is required to launch a container")
}

// checks if a command is available
func isCmdAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func imgFilename(name string) string {
	n := strings.Replace(name, ".", "_", -1)
	n = strings.Replace(n, "/", "_", -1)
	n = strings.Replace(n, "-", "_", -1)
	return n
}

func pkgFilename(name core.PackageName) string {
	filename := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%s:%s", name.FullyQualifiedName(), name.Tag), "/", "_"), ".", "_"), "-", "_")
	return fmt.Sprintf("%s.tar", filename)
}
