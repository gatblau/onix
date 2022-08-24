/*
  Onix Config Manager - Onix file exporter for OpenTelemetry
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package fileexporter

import (
	"errors"
	"fmt"
	"strconv"

	"go.opentelemetry.io/collector/config"
)

const (
	maxfilesize = int64(100) // 100kb
)

// Config defines configuration for file exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	// Path of the file to write to. Path is relative to current directory.
	Path       string `mapstructure:"path"`
	FileSizeKb string `mapstructure:"filesizekb"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {

	if cfg.Path == "" {
		return errors.New("path must be non-empty")
	}
	if cfg.FileSizeKb == "" {
		cfg.FileSizeKb = strconv.FormatInt(maxfilesize, 10)
	} else {
		s, err := strconv.ParseInt(cfg.FileSizeKb, 10, 64)
		if err != nil {
			return err
		}
		if s > maxfilesize {
			return fmt.Errorf(" file size %d defined in telem.yaml file is greater than max size %d allowed for the file", s, maxfilesize)
		}
	}

	return nil
}
