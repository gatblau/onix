package registry

import (
	"encoding/base64"
	"encoding/json"
)

type Repository struct {
	// the artefact repository (name without without tag)
	Repository string `json:"repository"`
	// the reference name of the artefact corresponding to different builds
	Artefacts []*Artefact `json:"artefacts"`
}

func (r *Repository) FindArtefact(id string) *Artefact {
	for _, artefact := range r.Artefacts {
		if artefact.Id == id {
			return artefact
		}
	}
	return nil
}

// updates the specified artefact
func (r *Repository) UpdateArtefact(a *Artefact) bool {
	position := -1
	for ix, artefact := range r.Artefacts {
		if artefact.Id == a.Id {
			position = ix
			break
		}
	}
	if position != -1 {
		r.Artefacts[position] = a
		return true
	}
	return false
}

type Artefact struct {
	// a unique identifier for the artefact calculated as the checksum of the complete seal
	Id string `json:"id"`
	// the type of application in the artefact
	Type string `json:"type"`
	// the artefact actual file name
	FileRef string `json:"file_ref"`
	// the list of Tags associated with the artefact
	Tags []string `json:"tags"`
	// the size
	Size string `json:"size"`
	// the creation time
	Created string `json:"created"`
}

func (a Artefact) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (a Artefact) ToJson() (string, error) {
	bs, err := json.Marshal(a)
	return base64.StdEncoding.EncodeToString(bs), err
}
