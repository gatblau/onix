/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0

  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
	"strings"
)

// the marker to indicate the start or end of a front marker segment
const frontMatterMarker = "+++"

var (
	// toml front matter segment start marker
	tomlMarker = fmt.Sprintf("%stoml", frontMatterMarker)
	// json front matter segment start marker
	jsonMarker = fmt.Sprintf("%sjson", frontMatterMarker)
	// yaml front matter segment start marker
	yamlMarker = fmt.Sprintf("%syaml", frontMatterMarker)
)

// parses application configuration information stored in Onix and separates pilot configuration (front matter)
// from application configuration
type parser struct {
}

// a distinct piece of application configuration
type appCfg struct {
	// the application configuration
	config string
	// the application configuration metadata (front matter)
	meta *frontMatter
}

// different types of markers to facilitate parsing of application configuration
// and extraction of front matter
type tokenType int

// the various types of tokens
const (
	// the start of a toml front matter segment
	Toml tokenType = iota
	// the start of a json front matter segment
	Json
	// the start of a yaml front matter segment
	Yaml
	// the end of any front matter segment
	EOF
)

// a token used by the tokenizer
type token struct {
	// type of token
	Type tokenType
	// the line number the marker was found
	LineNo int
}

// retrieve a list of tokens marking where the front matter segments are
// used by the parser to extract front matter and configuration from the single source content passed in
func (p *parser) tokenize(content string) ([]token, error) {
	tokens := make([]token, 0)
	lines := strings.Split(content, "\n")
	for lineNo, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), tomlMarker) {
			token := &token{
				Type:   Toml,
				LineNo: lineNo,
			}
			tokens = append(tokens, *token)
		} else if strings.HasPrefix(strings.ToLower(line), jsonMarker) {
			token := &token{
				Type:   Json,
				LineNo: lineNo,
			}
			tokens = append(tokens, *token)
		} else if strings.HasPrefix(strings.ToLower(line), yamlMarker) {
			token := &token{
				Type:   Yaml,
				LineNo: lineNo,
			}
			tokens = append(tokens, *token)
		} else if strings.HasPrefix(strings.ToLower(line), frontMatterMarker) {
			// if an end of front matter marker is found, there must me a previous start marker
			if len(tokens) == 0 || (tokens[len(tokens)-1].Type != Toml && tokens[len(tokens)-1].Type != Json && tokens[len(tokens)-1].Type != Yaml) {
				return nil, errors.New(fmt.Sprintf("end token found but not corresponding start token exist: start token was %v", tokens[len(tokens)-1].Type))
			}
			token := &token{
				Type:   EOF,
				LineNo: lineNo,
			}
			tokens = append(tokens, *token)
		}
	}
	return tokens, nil
}

// parse the application configuration content stored in Onix
func (p *parser) parse(content string) ([]*appCfg, error) {
	tokens, err := p.tokenize(content)
	if err != nil {
		return nil, err
	}

	confList := make([]*appCfg, 0)

	lines := strings.Split(content, "\n")
	for i := 0; i < len(tokens); i = i + 2 {
		cfg := new(appCfg)
		// gets the front matter
		frontMatterStr := sliceToStr(lines[tokens[i].LineNo+1 : tokens[i+1].LineNo])
		fm, err := unmarshallFrontMatter(frontMatterStr, tokens[i].Type)
		if err != nil {
			return nil, err
		}
		if ok, err := fm.valid(); !ok {
			return nil, errors.New(fmt.Sprintf("%v; parsing after '+++' marker number %v", err, i+1))
		}
		cfg.meta = &fm

		// gets the configuration
		var eof int
		if i+2 == len(tokens) {
			// if this is the last iteration there is not further marker for the end of the configuration
			// therefore uses the last line number
			eof = len(lines)
		} else {
			// if this is not the last iteration it can use the location of the marker of the start of the next
			// front matter block to determine the end of the previous configuration
			eof = tokens[i+2].LineNo - 1
		}
		configStr := sliceToStr(lines[tokens[i+1].LineNo:eof])
		cfg.config = configStr

		confList = append(confList, cfg)
	}
	return confList, nil
}

// get a string by joining the elements in the passed in string slice
func sliceToStr(slice []string) string {
	var buf bytes.Buffer
	for _, s := range slice {
		buf.Write([]byte(s))
		buf.Write([]byte("\n"))
	}
	return buf.String()
}

// unmarshal the passed content string in the format defined by the tokenType
// into a frontMatter structure
func unmarshallFrontMatter(content string, t tokenType) (frontMatter, error) {
	var (
		err error
		fm  = frontMatter{}
	)
	// unmarshall the front matter
	switch t {
	case Toml:
		err = toml.Unmarshal([]byte(content), &fm)
	case Json:
		err = json.Unmarshal([]byte(content), &fm)
	case Yaml:
		err = yaml.Unmarshal([]byte(content), &fm)
	}
	return fm, err
}
