package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const digestsFilename = "digests"

type digestCache struct {
	Values   map[string]string
	filename string
	dirty    bool
}

func NewDigests() (*digestCache, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get the working directory: %s", err)
	}
	return &digestCache{
		Values:   make(map[string]string),
		filename: filepath.Join(dir, digestsFilename),
	}, nil
}

func (d *digestCache) load() error {
	_, err := os.Stat(d.filename)
	if os.IsNotExist(err) {
		// creates an empty file
		return ioutil.WriteFile(d.filename, d.bytes(), os.ModePerm)
	} else {
		b, err := ioutil.ReadFile(d.filename)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %s", d.filename, err)
		}
		err = json.Unmarshal(b, d)
		if err != nil {
			return fmt.Errorf("cannot unmarshal file %s: %s", d.filename, err)
		}
		return nil
	}
}

func (d *digestCache) bytes() []byte {
	b, _ := json.Marshal(d)
	return b
}

func (d *digestCache) set(key, value string) {
	if d.Values[key] != value {
		d.dirty = true
		d.Values[key] = value
	}
}

// check if the digest has changed
func (d *digestCache) changed(base string, digest string) bool {
	for key, value := range d.Values {
		if key == base {
			return digest != value
		}
	}
	// if we got here then the digest was not found in the cache
	// so add the digest to the cache
	d.Values[base] = digest
	// by returning false, the trigger will not fire the first time
	// the policy is assessed
	return false
}

// save the digestCache cache
func (d *digestCache) save() error {
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("cannot marshal digest cache: %s", err)
	}
	err = ioutil.WriteFile(d.filename, b, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save digest cache: %s", err)
	}
	return nil
}
