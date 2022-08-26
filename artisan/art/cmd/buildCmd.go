/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
)

// BuildCmd builds an artisan package
type BuildCmd struct {
	Cmd         *cobra.Command
	branch      string
	gitTag      string
	packageName string
	gitToken    string
	from        string
	fromPath    string
	profile     string
	copySource  bool
	interactive bool
	target      string
	image       string
}

func NewBuildCmd() *BuildCmd {
	c := &BuildCmd{
		Cmd: &cobra.Command{
			Use:   "build [flags] /build/file/path or https://build/file/git/uri or /path/to/folder (without build file)",
			Short: "builds a package",
			Long: `
=============================================
Build a Package
=============================================

The command use to create Artisan packages.

Packages are combinations of a zip file and a json file stored in a package registry:
- the zip file can contain one or more files or folders.
- the json file act as a digital seal that carries both package information (e.g. manifest) and the integrity check information
  to determine the integrity of the package (digest).

Executable Packages:
===================
Packages that contain a build file in its root, are classed as executable.
Executable packages export functions that can be run in runtimes.
Runtimes are containers that have the tool-chain required to run a specific function.

Package Interface:
=================
The declaration of these functions (e.g. their interface) are stored in the package manifest.

The package interface is a combination of:
  - the function 
  - the function parameters
  - the runtime that should be used to run the function
`,
			Example: `
To build a package you should have a clear identification scheme, the package name, and the location of the file(s)
to be packaged.

The package name is similar to container images repository/tag combinations.
For example, assuming a package registry located at my-registry.com and a repository called repository-group/repository-name the package
name is <my-registry.com/repository-group/repository-name>

Packages can also be tagged. A tag is a piece of text that is added to the package to facilitate referring to it.
For example, a tag could be any combination of letters and numbers, such as the day and time the package was created.

A package which such name can be built, tagged, pushed to and pulled from a registry, and opened in the file system.

To build a package with the name above and the 010121-v2 tag using the my-build-profile in the build file in the current 
folder ".", the following command should be issued:

$ art build -t my-registry.com/repository-group/repository-name:010121-v2 -p my-build-profile .

NOTE: in general, and similarly to building a container image with a Dockerfile, the build command requires a build file that defines
at least one build profile. Build profiles specify which files to package.

IMPORTANT: if the path used in the build command does not contain a build file, artisan creates a "content" package of type "files".
Such package cannot execute any functions but it is only destined to serve as a packaging mechanism for general files.
In order to create a content package do the following:
1. create a folder and add any files you would like to package (note that a build file is not needed in the folder)
2. run the build command with the --target flag as follows:
  $ art build -t my-registry.com/repository-group/repository-name:tag --target ./folder/to/package

Packaging Container Images:
It is possible to package any container imge in the local container registry of the host where the artisan is running.
Use the --image flag to specify which image to package as follows:

  $ art build --image docker.io/mongo:5 -t localhost:8080/images/mongo:5
`,
		},
	}
	c.Cmd.Run = c.Run
	c.Cmd.Flags().StringVarP(&c.gitToken, "token", "k", "", "the git access token to use to read a build file remotely stored in a protected git repository")
	c.Cmd.Flags().StringVar(&c.image, "image", "", "builds a package with an embedded container image specified by this name")
	c.Cmd.Flags().StringVarP(&c.packageName, "package-name", "t", "", "package name and optionally a tag in the 'name:tag' format")
	c.Cmd.Flags().StringVarP(&c.fromPath, "path", "f", "", "if a git repository is specified as the location to the build file, it defines the path within the git repository where the build file is")
	c.Cmd.Flags().StringVar(&c.target, "target", "", "if a explicit target folder is defined, then build the package without relying on a build file")
	c.Cmd.Flags().StringVarP(&c.profile, "profile", "p", "", "the build profile to use. if not provided, the default profile defined in the build file is used. if no default profile is found, then the first profile in the build file is used.")
	c.Cmd.Flags().BoolVarP(&c.interactive, "interactive", "i", false, "if true, it prompts the user for information if not provided")
	c.Cmd.Flags().BoolVarP(&c.copySource, "copy", "c", false, "indicates if a copy should be made of the project files before building the package. it is only applicable if the source is in the file system.")
	c.Cmd.MarkFlagRequired("package-name")
	return c
}

func (b *BuildCmd) Run(_ *cobra.Command, args []string) {
	// validate build path
	switch len(args) {
	case 0:
		b.from = "."
	case 1:
		b.from = args[0]
	default:
		core.RaiseErr("too many arguments")
	}
	builder := build.NewBuilder("")
	name, err := core.ParseName(b.packageName)
	i18n.Err("", err, i18n.ERR_INVALID_PACKAGE_NAME)
	// if an image option as been specified uses the image builder
	if len(b.image) > 0 {
		core.CheckErr(release.BuildImagePackage(b.image, b.packageName, "", "", ""), "cannot build image package")
	} else {
		core.CheckErr(builder.Build(b.from, b.fromPath, b.gitToken, name, b.profile, b.copySource, b.interactive, b.target), "cannot build package")
	}
}
