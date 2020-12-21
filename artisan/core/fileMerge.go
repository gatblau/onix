/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func MergeFiles(filenames []string, envFilename string) {
	files := filenames
	regex, err := regexp.Compile("\\${(?P<NAME>[^}]*)}")
	if err != nil {
		log.Printf("cannot compile regex: %s\n", err)
		return
	}
	if len(files) == 0 {
		log.Printf("you must provide files to merge!\n")
		return
	}

	// load environment variables from file, if file not specified then try loading .env
	LoadEnvFromFile(envFilename)

	// loop through the specified configuration files
	for _, file := range files {
		merged := false
		if filepath.Ext(file) != ".tem" {
			RaiseErr("file '%s' does not have a .tem extension, cannot process it", file)
		}
		// read the file content
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("cannot read file %s: %s\n", file, err)
			return
		}
		content := string(bytes)
		// find all environment variable placeholders in the content
		vars := regex.FindAll(bytes, 1000)
		// loop though the found vars to merge
		for _, v := range vars {
			defValue := ""
			// removes placeholder marks: ${...}
			vname := strings.TrimSuffix(strings.TrimPrefix(string(v), "${"), "}")
			// is a default value defined?
			cut := strings.Index(vname, ":")
			// split default value and var name
			if cut > 0 {
				// get the default value
				defValue = vname[cut+1:]
				// get the name of the var without the default value
				vname = vname[0:cut]
			}
			// check the name of the env variable is not "PWD" as it can return the current directory in some OSs
			if vname == "PWD" {
				log.Printf("environment variable cannot be PWD, choose a different name\n")
				return
			}
			// fetch the env variable value
			ev := os.Getenv(vname)
			// if the variable is not defined in the environment
			if len(ev) == 0 {
				// if no default value has been defined
				if len(defValue) == 0 {
					log.Fatalf("environment variable '%s' required and not defined, cannot merge\n", vname)
				} else {
					// merge with the default value
					content = strings.Replace(content, string(v), defValue, 1000)
					merged = true
				}
			} else {
				// merge with the env variable value
				content = strings.Replace(content, string(v), ev, 1000)
				merged = true
			}
		}
		// if variables have been merged at all
		if merged {
			// override file with merged values
			err = writeToFile(file, content)
			if err != nil {
				log.Printf("cannot update config file: %s\n", err)
			}
		}
	}
}

func LoadEnvFromFile(envFilename string) {
	// attempt to load config variables from .env file if exists
	if len(envFilename) == 0 {
		// try to load .env file
		_, err := os.Stat(filepath.Join(WorkDir(), ".env"))
		if !os.IsNotExist(err) {
			godotenv.Load()
		}
	} else {
		var (
			err error
		)
		// load vars from specified file
		if !filepath.IsAbs(envFilename) {
			envFilename, err = filepath.Abs(envFilename)
			if err != nil {
				log.Fatalf("error converting environment file path to an absolute path: %v", err)
			}
		}
		godotenv.Load(envFilename)
	}
}

func writeToFile(filename string, data string) error {
	// create a file without the .tem extension
	file, err := os.Create(FilenameWithoutExtension(filename))
	if err != nil {
		return err
	}
	defer file.Close()
	// write the merged content into the file
	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	log.Printf("'%v' bytes written to file '%s'\n", len(data), FilenameWithoutExtension(filename))
	return file.Sync()
}
