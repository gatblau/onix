package release

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/oxlib/resx"
	"gopkg.in/yaml.v2"
)

func ExportDebianPackage(pkgNames []string, opts ExportOptions) error {

	targetUri := opts.TargetUri
	creds := opts.TargetCreds
	artHome := opts.ArtHome

	if len(pkgNames) == 0 {
		return fmt.Errorf("debian package name is missing : %s", pkgNames)
	}

	pName, err := core.ParseName(aptPkgName())
	if err != nil {
		return fmt.Errorf("invalid artisan package name %s : %s, ", aptPkgName(), err)
	}

	// if a target has been specified
	if len(targetUri) > 0 {
		// if a final slash does not exist add it
		if targetUri[len(targetUri)-1] != '/' {
			targetUri = fmt.Sprintf("%s/", targetUri)
		}
		// automatically adds a tar filename to the URI based on the package name:tag
		targetUri = fmt.Sprintf("%s%s", targetUri, aptPkgTarFileName())
	} else {
		return fmt.Errorf("a destination URI must be specified to export the image")
	}

	// execution path
	tmp, err := core.NewTempDir(artHome)
	if err != nil {
		return fmt.Errorf("cannot create temp folder for processing image archive: %s", err)
	}
	core.DebugLogger.Printf("location of temporary folder %s", tmp)

	// create a target folder for the artisan package
	targetFolder := filepath.Join(tmp, "build")
	err = os.MkdirAll(targetFolder, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("failed to create build folder : %s", err)
	}
	//create a patch folder where packages to be used for patching will be downloaded
	patchFolder := filepath.Join(targetFolder, "patch")
	err = os.MkdirAll(patchFolder, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("failed to create patch folder : %s", err)
	}

	//get the debian packages either locally or from remote
	err = getPackages(pkgNames, patchFolder)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("failed to get package and its dependencies for packages %s \n error :- %s", pkgNames, err)
	}

	// get all the package name as string from folder "patch" which will be used to build
	// backup command
	bckupCmds, err := buildBackupCmds(patchFolder)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("failed to read .deb file from path %s, \n error is %s", patchFolder, err)
	}

	// generate build function build.yaml for this package
	bfBytes, err := generateBuildFunctions(bckupCmds)
	if err != nil {
		return fmt.Errorf("failed to marshall debian package build file: %s", err)
	}

	// create a build file to build the package containing the debian packages tar
	pbfBytes, err := generateArtBuild(aptPkgName())
	if err != nil {
		return fmt.Errorf("failed to marshall debian packaging build file: %s", err)
	}

	//save package build and function build file
	core.InfoLogger.Println("packaging debian packages tarball file")
	err = os.WriteFile(filepath.Join(tmp, "build.yaml"), pbfBytes, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save debian packaging build file: %s", err)
	}
	err = os.WriteFile(filepath.Join(targetFolder, "build.yaml"), bfBytes, 0755)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save debian build function file: %s", err)
	}

	b := build.NewBuilder(artHome)
	b.Build(tmp, "", "", pName, "", false, false, "")
	r := registry.NewLocalRegistry(artHome)
	// export package
	core.InfoLogger.Printf("exporting debian package to tarball file")
	err = r.ExportPackage([]core.PackageName{*pName}, "", targetUri, creds)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot save debian package to destination: %s", err)
	}

	//append package name to the spec file
	spec := new(Spec)
	err = yaml.Unmarshal(opts.Specification.content, spec)
	if err != nil {
		os.RemoveAll(tmp)
		return fmt.Errorf("failed unmarshal spec file's content part: %s", err)
	}

	m := spec.Packages
	if m == nil {
		m = make(map[string]string)
		spec.Packages = m
	}
	m["PACKAGE_APT"] = pName.String()
	contents, err := core.ToYamlBytes(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal the spec file's content part: %s", err)
	}
	opts.Specification.content = contents

	return nil
}

func getPackages(pkgNames []string, target string) error {
	if len(pkgNames) == 0 {
		return fmt.Errorf("package names to be exported is empty")
	}

	// make sure pkgNames slice contains either all with debian package file name with .deb extension
	// or all will be debian package name (with no .deb extension)
	ext := filepath.Ext(pkgNames[0])
	prePkg := pkgNames[0]
	for _, p := range pkgNames {
		ex := filepath.Ext(p)
		if strings.Compare(strings.TrimSpace(ext), strings.TrimSpace(ex)) != 0 {
			return fmt.Errorf("all package names are not of same format, some has extension some don't [ %s ] [ %s ]", prePkg, p)
		}
		prePkg = p
	}

	// if the pkgNames slice contains deb package file names then copy them all from source folder
	// to target folder patch
	var err error
	if len(ext) > 0 {
		err = copyPackagesFromLocal(pkgNames, target)
	} else {
		err = downloadPackagesFromRemote(pkgNames, target)
	}

	return err
}

func downloadPackagesFromRemote(pkgNames []string, target string) error {

	// convert pkgNames slice to a single string, each element of slice separated by space
	allPkgs := strings.Join(pkgNames, " ")

	// get name of dependent packages
	compList, err := getDependencies(allPkgs, target)
	if err != nil {
		return fmt.Errorf("failed to query dependencies for packages %s \n error :- %s", allPkgs, err)
	}

	// download all the packages incudling dependencies
	err = downloadPackages(compList, target)
	if err != nil {
		return fmt.Errorf("failed to download package [ %s ] and its dependencies \n error :- %s", compList, err)
	}

	return nil
}

func copyPackagesFromLocal(pkgNames []string, target string) error {
	for _, p := range pkgNames {
		ext := filepath.Ext(p)
		if ext != ".deb" {
			return fmt.Errorf("package names extension is not .deb %s", p)
		}
		b, e := resx.ReadFile(p, "")
		if e != nil {
			return nil
		}
		core.DebugLogger.Printf("< byte size [ %d ] \n file path [ %s ]\n >", len(b), filepath.Join(target, filepath.Base(p)))
		e = resx.WriteFile(b, filepath.Join(target, filepath.Base(p)), "")
		if e != nil {
			return nil
		}
	}

	return nil
}

func aptPkgTarFileName() string {
	r := strings.NewReplacer(
		"/", "_",
		".", "_",
	)
	fileWithExt := fmt.Sprintf("%s.%s", r.Replace(aptPkgName()), "tar")
	return fileWithExt
}

func aptPkgName() string {
	return fmt.Sprintf("127.0.0.1/os/packages/apt:%d", time.Now().Unix())
}

func generateBuildFunctions(bckupCmds []string) ([]byte, error) {
	export_yes := true
	export_no := false
	patchCmd := "(cd patch && echo \"\" && echo \"************ patching started ***************\"" +
		" && sudo dpkg --force-all -i * && echo \"************ patch completed ***************\" && echo \"\")"
	rollbackCmd := "(cd backup && echo \"\" && echo \"************ patching failed, rollback started ***************\"" +
		" && sudo dpkg --force-all -i * && echo \"************ rollback completed ***************\" && echo \"\")"
	combinedCmd := "bash -c '" + patchCmd + " || " + rollbackCmd + "'"
	bf := data.BuildFile{
		Runtime: "ubi-min",
		Functions: []*data.Function{
			{
				Name:        "apply",
				Description: "apply debian packages to the operating system. In case of failure, will automatically rollback",
				Export:      &export_yes,
				Run: []string{
					"bash -c 'mkdir -p backup'",
					"$(backup)",
					combinedCmd,
				},
			},
			{
				Name:        "backup",
				Description: "take backup of current version of packages",
				Export:      &export_no,
				Run:         bckupCmds,
			},
		},
	}
	return yaml.Marshal(bf)
}

func buildBackupCmds(patchFolder string) ([]string, error) {
	files, err := ioutil.ReadDir(patchFolder)
	if err != nil {
		return nil, err
	}
	core.DebugLogger.Printf("count of .deb files found at the path [ %s] is [ %d ] ", patchFolder, len(files))
	var bckupCmds []string
	core.InfoLogger.Printf("generating run command for .deb files")
	exp := regexp.MustCompile(`\r?\n`)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".deb" {
			//go to patch folder and get the debian package name
			cmd := fmt.Sprintf("bash -c 'cd %s && dpkg -f %s Package'", patchFolder, file.Name())
			pname, er := build.Exe(cmd, patchFolder, merge.NewEnVarFromSlice([]string{}), false)
			pname = exp.ReplaceAllString(pname, " ")
			if er != nil {
				return nil, er
			}
			bckupCmd := fmt.Sprintf("bash -c \"cd backup && sudo dpkg-query -W %s | awk '{print $1}' | sudo xargs dpkg-repack\"", pname)
			bckupCmds = append(bckupCmds, bckupCmd)
		}
	}

	return bckupCmds, nil
}

func downloadPackages(compList, patchFolder string) error {
	downloadPkg := "bash -c 'apt-get download " + strings.TrimSpace(compList) + "'"
	core.InfoLogger.Printf("downloading package %s and its dependencies", compList)
	core.DebugLogger.Printf("downloading debian package [%s] dependencies using command :\n %s", compList, downloadPkg)
	_, err := build.Exe(downloadPkg, patchFolder, merge.NewEnVarFromSlice([]string{}), false)

	return err
}

func generateArtBuild(pkgName string) ([]byte, error) {
	// create a build file to build the package containing the debian packages tar
	pbf := data.BuildFile{
		Labels: map[string]string{
			"package": pkgName,
		},
		Profiles: []*data.Profile{
			{
				Name:   "debian-package",
				Target: "./build",
				Type:   "content/apt",
			},
		},
	}
	return yaml.Marshal(pbf)
}

func getDependencies(pkgNames, executionPath string) (string, error) {

	// get all dependencies for current package
	pkgNames = strings.TrimSpace(pkgNames)
	qryDependencies := fmt.Sprintf("bash -c 'apt-rdepends %s | grep -v \"^ \" | sed 's/debconf-2.0/debconf/g'' | sed 's/time-daemon/systemd-timesyncd/g'", pkgNames)
	core.DebugLogger.Printf("querying debian package [%s] dependencies using command :\n %s\n ", pkgNames, qryDependencies)
	// execute the command synchronously
	core.InfoLogger.Printf("querying dependencies of package %s", pkgNames)
	//using async because when using sync occassionally it been observed that return list contains
	//status messages key word also like "processing" along with the dependent package names
	dep, err := build.ExeAsync(qryDependencies, executionPath, merge.NewEnVarFromSlice([]string{}), false)
	//note: the dep list will contain the parent package name also for which we looked dependencies
	if err != nil {
		return "", err
	}

	replacer := strings.NewReplacer("\n", " ", "  ", " ", "\t", "")
	dep = replacer.Replace(dep)
	return dedup(dep)
}

func dedup(dep string) (string, error) {
	core.DebugLogger.Printf("dedup input data is \n %s \n", dep)
	var b strings.Builder
	words := strings.Split(dep, " ")

	for _, word := range words {
		//find exact match for the word in the string builder
		w := strings.Replace(fmt.Sprintf("\b%s\b", word), "+", "\\+", -1)
		x, err := regexp.MatchString(w, b.String())
		if err != nil {
			return "", fmt.Errorf("failed dedup process while trying for word %s in new package list %s: %s", word, b.String(), err)
		}

		if x {
			continue
		}
		b.WriteString(word)
		b.WriteString(" ")
	}
	return strings.TrimSpace(b.String()), nil
}
