/*
  Onix Config Manager - Build Manager
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

// response from skopeo inspect docker://image-name-tag
type ImgInfo struct {
	Name          string            `json:"Name"`
	Digest        string            `json:"Digest"`
	RepoTags      []string          `json:"RepoTags"`
	Created       string            `json:"Created"`
	DockerVersion string            `json:"DockerVersion"`
	Labels        map[string]string `json:"Labels"`
	Architecture  string            `json:"Architecture"`
	Os            string            `json:"Os"`
	Layers        []string          `json:"Layers"`
	Env           []string          `json:"Env"`
}
