//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/oxc"
	"net/http"
)

type Plan struct {
	Releases []Info `json:"releases"`
}

type Info struct {
	DbVersion  string `json:"dbVersion"`
	AppVersion string `json:"appVersion"`
	Path       string `json:"path"`
}

// get a JSON bytes reader for the Plan
func (plan *Plan) json() (*bytes.Reader, error) {
	jsonBytes, err := plan.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Plan
func (plan *Plan) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(plan)
	return &b, err
}

// get the Plan in the http Response
func (plan *Plan) decode(response *http.Response) (*Plan, error) {
	result := new(Plan)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// return the info for the getReleaseInfo of the specified app version
// Info: getReleaseInfo information
func (plan *Plan) info(appVersion string) *Info {
	for _, release := range plan.Releases {
		if release.AppVersion == appVersion {
			return &release
		}
	}
	return nil
}

// check if an upgrade path is available from the current to the target app version
func (plan *Plan) canUpgrade(currentAppVersion string, targetAppVersion string) (bool, string) {
	current := plan.info(currentAppVersion)
	if current == nil {
		return false, fmt.Sprintf("!!! I could not find information for current application version %s", currentAppVersion)
	}
	target := plan.info(targetAppVersion)
	if target == nil {
		return false, fmt.Sprintf("!!! I could not find information for target application version %s", targetAppVersion)
	}
	return true, ""
}

func (plan *Plan) getUpgradeWindow(currentAppVersion string, targetAppVersion string) (currentReleaseIndex int, targetReleaseIndex int) {
	var (
		currentIx, targetIx int
	)
	for ix, release := range plan.Releases {
		if release.AppVersion == currentAppVersion {
			currentIx = ix + 1
		}
		if release.AppVersion == targetAppVersion {
			targetIx = ix + 1
			break
		}
	}
	return currentIx, targetIx
}
