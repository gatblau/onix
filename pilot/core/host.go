/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"github.com/rs/zerolog/log"
)

type host struct {
}

func NewHost() (*host, error) {
	return &host{}, nil
}

func (h *host) Start() {
	info, err := NewHostInfo()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	log.Info().Msgf(info.String())
}
