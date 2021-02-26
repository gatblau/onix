package data

import "strings"

type Var struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required" json:"required"`
	Type        string `yaml:"type" json:"type"`
	Value       string `yaml:"value,omitempty" json:"value,omitempty"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty"`
}

type Vars []*Var

func (list Vars) Len() int { return len(list) }

func (list Vars) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

func (list Vars) Less(i, j int) bool {
	var si string = list[i].Name
	var sj string = list[j].Name
	var si_lower = strings.ToLower(si)
	var sj_lower = strings.ToLower(sj)
	if si_lower == sj_lower {
		return si < sj
	}
	return si_lower < sj_lower
}
