/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/artisan/runner"
	"github.com/gatblau/onix/oxlib/httpserver"
	"os"
	"testing"
)

func TestExeC(t *testing.T) {
	packageName := "uri/recipe/java-quarkus"
	fxName := "setup"
	// create an instance of the runner
	run, err := runner.New()
	core.CheckErr(err, "cannot initialise runner")
	env, err := merge.NewEnVarFromFile(".env")
	if err != nil {
		fmt.Printf("cannot load env file: %s\n", err.Error())
		t.FailNow()
	}
	// launch a runtime to execute the function
	err = run.ExeC(packageName, fxName, "admin:sss", "", false, env)
	i18n.Err("", err, i18n.ERR_CANT_EXEC_FUNC_IN_PACKAGE, fxName, packageName)
}

func TestExe(t *testing.T) {
	packageName, err := core.ParseName("test")
	fxName := "t1"
	builder := build.NewBuilder(ArtHome)
	core.CheckErr(err, "cannot initialise builder")
	env, err := merge.NewEnVarFromFile(".env")
	if err != nil {
		fmt.Printf("cannot load env file: %s\n", err.Error())
		t.FailNow()
	}
	// launch a runtime to execute the function
	builder.Execute(packageName, fxName, "admin:sss", true, "", false, env)
}

func TestBuild(t *testing.T) {
	packageName, _ := core.ParseName("test")
	builder := build.NewBuilder(ArtHome)
	builder.Build(".", "", "", packageName, "", false, false, "")
}

func TestRunC(t *testing.T) {
	run, err := runner.NewFromPath(".", ArtHome)
	core.CheckErr(err, "cannot initialise runner")
	err = run.RunC("deploy", false, merge.NewEnVarFromSlice([]string{}), "")
}

func TestPush(t *testing.T) {
	reg := registry.NewLocalRegistry(ArtHome)
	name, err := core.ParseName("localhost:8080/lib/test1")
	if err != nil {
		t.FailNow()
	}
	err = reg.Push(name, "admin:admin", false)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
}

func TestPull(t *testing.T) {
	reg := registry.NewLocalRegistry(ArtHome)
	name, err := core.ParseName("localhost:8082/gatblau/tools/artisan")
	if err != nil {
		t.FailNow()
	}
	reg.Pull(name, "admin:admin", false)
}

func TestRLs(t *testing.T) {
	reg, _ := registry.NewRemoteRegistry("localhost:8080", "admin", "adm1n", ArtHome)
	reg.List(false)
}

func TestVars(t *testing.T) {
	env, _ := merge.NewEnVarFromFile(".env")
	builder := build.NewBuilder(ArtHome)
	builder.Run("test", ".", false, env)
}

// test the merging of .tem templates
func TestMergeTem(t *testing.T) {
	filename := "test/test.txt"
	tm, err := merge.NewTemplMerger()
	checkErr(err, t)
	err = tm.LoadTemplates([]string{filename + ".tem"})
	checkErr(err, t)
	err = tm.Merge(merge.NewEnVarFromSlice([]string{"VAR1=World"}))
	checkErr(err, t)
	tm.Save()
	_, err = os.Stat(filename)
	checkErr(err, t)
	_ = os.Remove(filename)
}

// test the merging of .art templates
func TestMergeArt(t *testing.T) {
	filename := "test/test.txt"
	tm, err := merge.NewTemplMerger()
	checkErr(err, t)
	err = tm.LoadTemplates([]string{filename + ".art"})
	checkErr(err, t)
	err = tm.Merge(merge.NewEnVarFromSlice([]string{"VAR1=World"}))
	checkErr(err, t)
	tm.Save()
	_, err = os.Stat(filename)
	checkErr(err, t)
	_ = os.Remove(filename)
}

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
}

func TestRun(t *testing.T) {
	builder := build.NewBuilder(ArtHome)
	// add the build file level environment variables
	env := merge.NewEnVarFromSlice(os.Environ())
	// execute the function
	builder.Run("release-bin", "../", false, env)
}

func TestCurl(t *testing.T) {
	core.Curl("http://localhost:8080/user/ONIX_PILOTCTL",
		"PUT",
		httpserver.BasicToken("admin", "0n1x"),
		[]int{200, 201},
		"{\n  \"email\":\"a@a.com\", \"name\":\"aa\", \"pwd\":\"aaAA88!=12222\", \"service\":\"false\", \"acl\":\"*:*:*\"\n}",
		"",
		5,
		5,
		5,
		[]string{"Content-Type: application/json"},
		"")
}

func TestSave(t *testing.T) {
	names, err := core.ValidateNames([]string{"test", "artisan"})
	if err != nil {
		t.Error(err)
	}
	r := registry.NewLocalRegistry(ArtHome)
	err = r.ExportPackage(names, "", "./export", "")
	if err != nil {
		t.Error(err)
	}
}

func TestImport(t *testing.T) {
	// create a local registry
	r := registry.NewLocalRegistry(ArtHome)
	// import the tar archive(s)
	err := r.Import([]string{"../archive.tar"}, "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemove(t *testing.T) {
	r := registry.NewLocalRegistry(ArtHome)
	p := r.AllPackages()
	for _, s := range p {
		fmt.Println(s)
	}
	err := r.Remove(p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveRemote(t *testing.T) {
	r, _ := registry.NewRemoteRegistry("localhost:8080", "admin", "admin", "")
	err := r.RemoveByNameOrId([]string{"cfe1761845c7", "fb7d78733eaf"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewSpecPush(t *testing.T) {
	// err := export.PushSpec()
}
