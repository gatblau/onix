/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

//git events to be handled by host runner
type Push struct {
	Ref        string     `json:"ref"`
	Commits    []Commits  `json:"commits"`
	Repository Repository `json:"repository"`
}
type Commits struct {
	Id string `json:"id"`
}
type Repository struct {
	Name string `json:"name"`
	// json format from gitlab
	GitHttpUrl string `json:"git_http_url"`
	// json format from github
	GitUrl string `json:"clone_url"`
}

//NewPushEvent build new git push event
func NewPushEvent(flowJSONBytes []byte) (*Push, error) {
	p := new(Push)
	err := json.Unmarshal(flowJSONBytes, p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Push event definition %s", err)
	}
	return p, nil
}

//IsValidUri validate whether Push event has git url or not and whether git url
// available in the Push event matches to the one in flow spec
func (p *Push) IsValidUri(uri string) error {
	var r string
	if len(p.Repository.GitHttpUrl) > 0 {
		r = p.Repository.GitHttpUrl
	} else {
		r = p.Repository.GitUrl

	}

	if len(r) == 0 {
		return errors.New("invalid git push event, git uri is missing from push event")
	}

	u1, err := url.Parse(r)
	if err != nil {
		panic(err)
	}

	u2, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	if u1.Path != u2.Path {
		msg := fmt.Sprintf("git uri in flow spec is not same as git uri received in push event [ %s ] , [ %s ]", uri, r)
		return errors.New(msg)

	}

	return nil
}

//IsValidBranch validate Push event for branch name exist or not
func (p *Push) IsValidBranch(branch string) error {
	b := p.Ref[strings.LastIndex(p.Ref, "/")+1:]
	// check if push event received is for same branch as branch mentioned in flow spec.
	if len(branch) > 0 && branch != b {
		msg := fmt.Sprintf("git branch in flow spec is not same as git ref received in push event [ %s ] , [ %s ]", branch, b)
		return errors.New(msg)
	}
	return nil
}
