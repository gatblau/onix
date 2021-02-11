/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
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
	"github.com/ohler55/ojg/jp"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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
	fmt.Printf("error!\n* %s\n", fmt.Sprintf(msg, a...))
	os.Exit(1)
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
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
		RaiseErr("run command failed in build.yaml: %s", err)
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
						RaiseErr("the environment variable '%s' is not defined, are you missing a binding? you can always run the command in interactive mode to manually input its value", name)
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

// gets a random string of specified length
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ToAbsPath(flowPath string) string {
	if !path.IsAbs(flowPath) {
		abs, err := filepath.Abs(flowPath)
		CheckErr(err, "cannot convert '%s' to absolute path", flowPath)
		flowPath = abs
	}
	return flowPath
}

// encode strings to be used in tekton pipelines names
func Encode(value string) string {
	length := 30
	value = strings.ToLower(value)
	value = strings.Replace(value, " ", "-", -1)
	if len(value) > length {
		value = value[0:length]
	}
	return value
}

func Wait(uri, filter, token string, maxAttempts int) {
	var (
		filtered []interface{}
		attempts = 0
	)
	// executes the query
	filtered = httpGetFiltered(uri, token, filter)
	// if no result loop
	for len(filtered) == 0 {
		// wait for next attempt
		time.Sleep(500 * time.Millisecond)
		// executes query
		filtered = httpGetFiltered(uri, token, filter)
		// increments the number of attempts
		attempts++
		// exits if max attempts reached
		if attempts >= maxAttempts {
			RaiseErr("call to %s did not return expected value after %d attempts", uri, maxAttempts)
		}
	}
}

func httpGetFiltered(uri, token, filter string) []interface{} {
	result := httpGet(uri, token)
	var jason interface{}
	err := json.Unmarshal(result, &jason)
	CheckErr(err, "cannot unmarshal response")
	// filtered, err = jsonpath.Read(jason, filter)
	f, err := jp.ParseString(filter)
	CheckErr(err, "cannot apply filter")
	return f.Get(jason)
}

func httpGet(uri, token string) []byte {
	// create request
	req, err := http.NewRequest("GET", uri, nil)
	CheckErr(err, "cannot create new request")
	// add authorization header if there is a token defined
	if len(token) > 0 {
		req.Header.Set("Authorization", token)
	}
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	CheckErr(err, "cannot call URI %s", uri)
	if resp.StatusCode > 299 {
		RaiseErr("http request return error: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err, "cannot read response body")
	// if the result is not in JSON format
	if !isJSON(body) {
		RaiseErr("the http response body was not in json format, cannot apply JSON path filter")
	}
	return body
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

// defaults to quay.io/artisan root if not specified
func QualifyRuntime(runtime string) string {
	// container images must be in lower case
	runtime = strings.ToLower(runtime)
	// if no repository is specified then assume artisan library at quay.io/artisan
	if !strings.ContainsAny(runtime, "/") {
		return fmt.Sprintf("quay.io/artisan/%s", runtime)
	}
	return runtime
}
