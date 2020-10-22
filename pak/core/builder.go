package core

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Builder struct {
	zipWriter            *zip.Writer
	workingDir           string
	pakFilename          string
	cmds                 []string
	labels               map[string]string
	fileOrPathToBeZipped string
}

func NewBuilder() *Builder {
	return &Builder{
		cmds:   []string{},
		labels: make(map[string]string),
	}
}

// the directory where the git cloned files are stored
// from there they will be zipped
// it is named "pak" so that when it is zipped a pak folder will always exists after unzipping
func (p *Builder) sourceDir() string {
	return fmt.Sprintf("%s/pak", p.workingDir)
}

func (p *Builder) Build(repoUrl string) {
	// creates a temporary working directory
	p.newWorkingDir()
	// clone the remote repository
	repo, err := git.PlainClone(p.sourceDir(), false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		_ = os.RemoveAll(p.workingDir)
		log.Fatal(err)
	}
	// load the pakfile
	p.loadPakfile()
	// run commands
	p.run()
	// set the package name
	p.pakName(repo)
	// defines the destination (i.e. the *.pak file) within the working directory
	dest := fmt.Sprintf("%s/%s.pak", p.workingDir, p.pakFilename)
	// defines the source for zipping as specified in the Pakfile within the source directory
	source := fmt.Sprintf("%s/%s", p.sourceDir(), p.fileOrPathToBeZipped)
	// create the zip package
	zipSource(source, dest)
	// cleanup all relevant folders and move package to target location
	p.cleanUp()
}

// cleanup all relevant folders and move package to target location
func (p *Builder) cleanUp() {
	// remove the zip folder
	p.removeFromWD("pak")
	// check the home dir exists
	p.checkHomeDir()
	// move the package to the user home
	p.moveToHome(fmt.Sprintf("%s.pak", p.pakFilename))
	// remove the working directory
	err := os.RemoveAll(p.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	// set the directory to empty
	p.workingDir = ""
}

// move the specified filename from the working directory to the home directory (~/.pak/)
func (p *Builder) moveToHome(filename string) {
	err := os.Rename(fmt.Sprintf("%s/%s", p.workingDir, filename), fmt.Sprintf("%s/.pak/%s", p.homeDir(), filename))
	if err != nil {
		log.Fatal(err)
	}
}

// check the home directory exists and if not creates it
func (p *Builder) checkHomeDir() {
	// check the home directory exists
	_, err := os.Stat(fmt.Sprintf("%s/.pak", p.homeDir()))
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("%s/.pak", p.homeDir()), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// create a new working directory and return its path
func (p *Builder) newWorkingDir() {
	basePath, _ := os.Getwd()
	uid := uuid.New()
	folder := strings.Replace(uid.String(), "-", "", -1)
	workingDirPath := fmt.Sprintf("%s/.%s", basePath, folder)
	// creates a temporary working directory
	err := os.Mkdir(workingDirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	p.workingDir = workingDirPath
	// create a sub-folder to zip
	err = os.Mkdir(p.sourceDir(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// zips a directory
// func (p *Builder) pakFolder(folder string) error {
// 	// construct the fqn for the zip file
// 	pakFilename := fmt.Sprintf("%s/%s.pak", p.workingDir, p.pakFilename)
// 	// create the zip file
// 	file, err := os.Create(pakFilename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	p.zipWriter = zip.NewWriter(file)
// 	defer p.zipWriter.Close()
// 	walker := p.walk
// 	return filepath.Walk(folder, walker)
// }

// zip a file or a folder
func zipSource(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

// tree walker function for zip
func (p *Builder) walk(path string, info os.FileInfo, err error) error {
	fmt.Printf("compressing: %#v\n", path)
	if err != nil {
		return err
	}
	// if it is a directory return
	if info.IsDir() {
		return nil
	}
	// it is a file so open it
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	// ensure the file IO is closed eventually
	defer file.Close()
	// must use a relative path for it to work
	relativePath := strings.TrimPrefix(path, fmt.Sprintf("%s/", p.workingDir))
	// adds the file to the zip file
	f, err := p.zipWriter.Create(relativePath)
	if err != nil {
		return err
	}
	// copy the zip content to file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

// construct a unique name for the package using the short HEAD commit hash and current time
func (p *Builder) pakName(repo *git.Repository) {
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	// get the current time
	t := time.Now()
	timeStamp := fmt.Sprintf("%d%d%s%d%d%d%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[:2], t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:3])
	p.pakFilename = fmt.Sprintf("%s-%s", timeStamp, ref.Hash().String()[:7])
}

// gets the user home directory
func (p *Builder) homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// remove from working directory
func (p *Builder) removeFromWD(path string) {
	err := os.RemoveAll(fmt.Sprintf("%s/%s", p.workingDir, path))
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Builder) removeIgnored() {
	// retrieve .pakignore
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/.pakignore", p.sourceDir()))
	if err != nil {
		// assume no .pakignore exists, do nothing
		log.Printf(".packignore not found")
		return
	}
	// get the lines in the ignore file
	lines := strings.Split(string(bytes), "\n")
	// loop and remove the included files or folders
	for _, line := range lines {
		path := fmt.Sprintf("%s/%s", p.sourceDir(), line)
		err := os.RemoveAll(path)
		if err != nil {
			log.Printf("failed to ignore file %s", path)
		}
	}
}

// parse the Pakfile build instructions
func (p *Builder) loadPakfile() {
	// retrieve Pakfile
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/Pakfile", p.sourceDir()))
	if err != nil {
		log.Fatal(err)
	}
	// get the lines in the Pakfile
	lines := strings.Split(string(bytes), "\n")
	// loop load data
	for _, line := range lines {
		// add labels
		if strings.HasPrefix(line, "LABEL ") {
			value := line[6:]
			parts := strings.Split(value, "=")
			p.labels[strings.Trim(parts[0], " ")] = strings.Trim(strings.Trim(parts[1], " "), "\"")
		}
		// add commands
		if strings.HasPrefix(line, "RUN ") {
			value := line[4:]
			p.cmds = append(p.cmds, value)
		}
		// add the output path
		if strings.HasPrefix(line, "PATH ") {
			p.fileOrPathToBeZipped = line[5:]
		}
	}
}

func (p *Builder) run() {
	for _, cmd := range p.cmds {
		execute(cmd, p.sourceDir())
	}
	p.waitForFileExist(p.fileOrPathToBeZipped, 5*time.Second)
}

// executes a command
func execute(cmd string, dir string) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		//nolint:gosec
		c = exec.Command(strArr[0])
	} else {
		//nolint:gosec
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	c.Dir = dir
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	log.Printf("executing: %s\n", strings.Join(c.Args, " "))
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	err := c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// check if a relative path in the zip directory is a directory
func (p *Builder) isDir(path string) bool {
	info, err := os.Stat(path)
	// if it does not
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
	return info.IsDir()
}

// wait a time duration for a file to be created on the path
func (p *Builder) waitForFileExist(path string, d time.Duration) {
	elapsed := 0
	for {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) || elapsed > 20 {
			break
		}
		elapsed++
		time.Sleep(500 * time.Millisecond)
	}
}
