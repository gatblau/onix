package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
)

// ReMan remote service manager API
type ReMan struct {
	conf *Conf
	db   *Db
}

func NewReMan() *ReMan {
	cfg := NewConf()
	db := NewDb(cfg.getDbHost(), cfg.getDbPort(), cfg.getDbName(), cfg.getDbUser(), cfg.getDbPwd())
	return &ReMan{db: db, conf: cfg}
}

func (r *ReMan) Register(registration *Registration) error {
	_, err := r.db.RunQuery(fmt.Sprintf("select rem_beat('%s')", registration.MachineId))
	if err != nil {
		return err
	}
	return nil
}
