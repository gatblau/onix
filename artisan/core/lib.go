/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/hashicorp/go-uuid"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// converts the path to absolute path
func ToAbs(path string) string {
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		CheckErr(err, "cannot return an absolute representation of path")
		path = abs
	}
	return path
}

// convert the passed in parameter to a Json Byte Array
func ToJsonBytes(s interface{}) []byte {
	// serialise the seal to json
	source, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	err = json.Indent(dest, source, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return dest.Bytes()
}

func ToJsonFile(obj interface{}) (*os.File, error) {
	// create an UUId
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	// generate an internal random and transient name based on the UUId
	name := path.Join(TmpPath(), fmt.Sprintf("%s.json", uuid))
	// marshals the object into Json bytes
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// create a transient temp file
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	// write the bytes into it
	_, err = file.Write(b)
	if err != nil {
		return nil, err
	}
	// closes the file
	file.Close()
	// open the created file
	file = openFile(name)
	if err != nil {
		return nil, err
	}
	// remove the file from the file system
	err = os.Remove(name)
	// return the File object
	return file, err
}

func openFile(path string) *os.File {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return r
}

// remove an element in a slice
func RemoveElement(a []string, value string) []string {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == value {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = ""   // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

func CheckErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Printf("error: %s - %s\n", fmt.Sprintf(msg, a...), err)
		os.Exit(1)
	}
}

func RaiseErr(msg string, a ...interface{}) {
	fmt.Printf("error: %s\n", fmt.Sprintf(msg, a...))
	os.Exit(1)
}

func Msg(msg string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Printf("info: %s\n", fmt.Sprintf(msg, a...))
	} else {
		fmt.Printf("info: %s\n", msg)
	}
}

// gets the checksum of the passed string
func StringCheckSum(value string) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader([]byte(value))); err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func UserPwd(creds string) (user, pwd string) {
	// if credentials not provided then no user / pwd
	if len(creds) == 0 {
		return "", ""
	}
	parts := strings.Split(creds, ":")
	if len(parts) != 2 {
		log.Fatal(errors.New("credentials in incorrect format, they should be USER:PWD"))
	}
	return parts[0], parts[1]
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

// return a valid absolute path that exists
func AbsPath(filePath string) (string, error) {
	var p = filePath
	if !filepath.IsAbs(filePath) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return "", err
		}
		p = absPath
	}
	_, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	return p, nil
}

func HandleCtrlC(err error) {
	if err == terminal.InterruptErr {
		fmt.Println("\ncommand interrupted")
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}
}

// merges environment variables in the arguments
// returns the merged command list and the updated environment variables map if interactive mode is used
func MergeEnvironmentVars(args []string, env map[string]string, interactive bool) ([]string, map[string]string) {
	var result = make([]string, len(args))
	// the updated environment if interactive mode is used
	var updatedEnv = env
	// env variable regex
	evExpression := regexp.MustCompile("\\${(.*?)}")
	// check if the args have env variables and if so merge them
	for ix, arg := range args {
		result[ix] = arg
		// find all environment variables in the argument
		matches := evExpression.FindAllString(arg, -1)
		// if we have matches
		if matches != nil {
			for _, match := range matches {
				// get the name of the environment variable i.e. the name part in "${name}"
				name := match[2 : len(match)-1]
				// get the value of the variable
				value := env[name]
				// if not value exists and is not an embedded process variable
				if len(value) == 0 && !strings.HasPrefix(name, "ARTISAN_") {
					// if running in interactive mode
					if interactive {
						// prompt for the value
						prompt := &survey.Input{
							Message: fmt.Sprintf("%s:", name),
						}
						HandleCtrlC(survey.AskOne(prompt, &value, survey.WithValidator(survey.Required)))
						// add the variable to the updated environment map
						updatedEnv[name] = value
					} else {
						// if non-interactive then raise an error
						RaiseErr("environment variable '%s' is not defined", name)
					}
				}
				// merges the variable
				result[ix] = strings.Replace(result[ix], match, value, -1)
			}
		}
	}
	return result, updatedEnv
}

func HasFunction(value string) (bool, string) {
	matches := regexp.MustCompile("\\$\\((.*?)\\)").FindAllString(value, 1)
	if matches != nil {
		return true, matches[0][2 : len(matches[0])-1]
	}
	return false, ""
}

func HasShell(value string) (bool, string, string) {
	// pattern before escaping = "\$\(\((.*?)\)\)"
	matches := regexp.MustCompile("\\$\\(\\((.*?)\\)\\)").FindAllString(value, 1)
	if matches != nil {
		return true, matches[0], matches[0][3 : len(matches[0])-2]
	}
	return false, "", ""
}
