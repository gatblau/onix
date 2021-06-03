package core

import "fmt"

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

func NewUpdateConnStatusJob() (*UpdateConnStatusJob, error) {
	conf := NewConf()
	db, err := NewDb(conf.getDbHost(), conf.getDbPort(), conf.getDbName(), conf.getDbUser(), conf.getDbPwd())
	if err != nil {
		return nil, err
	}
	return &UpdateConnStatusJob{
		db:           db,
		pingInterval: conf.GetPingInterval(),
	}, nil
}

// UpdateConnStatusJob updates the connection status based on ping age
type UpdateConnStatusJob struct {
	db           *Db
	pingInterval int
}

func (c *UpdateConnStatusJob) Execute() {
	_, err := c.db.RunCommand([]string{fmt.Sprintf("select rem_record_conn_status('%d secs')", c.pingInterval)})
	if err != nil {
		fmt.Printf("ERROR: cannot check for disconnected events, %s\n", err)
	}
}

func (c *UpdateConnStatusJob) Description() string {
	return "updates the disconnected status of hosts based on last seen time"
}

func (c *UpdateConnStatusJob) Key() int {
	return hashCode(c.Description())
}
