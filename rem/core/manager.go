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
	"github.com/gatblau/oxc"
	"strconv"
)

// ReMan remote service manager API
type ReMan struct {
	conf *Conf
	db   *Db
	ox   *oxc.Client
}

func NewReMan() (*ReMan, error) {
	cfg := NewConf()
	db := NewDb(cfg.getDbHost(), cfg.getDbPort(), cfg.getDbName(), cfg.getDbUser(), cfg.getDbPwd())
	oxcfg := &oxc.ClientConf{
		BaseURI:            cfg.getOxWapiUrl(),
		Username:           cfg.getOxWapiUsername(),
		Password:           cfg.getOxWapiPassword(),
		InsecureSkipVerify: cfg.getOxWapiInsecureSkipVerify(),
	}
	oxcfg.SetAuthMode("basic")
	ox, err := oxc.NewClient(oxcfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create onix http client: %s", err)
	}
	return &ReMan{
		db:   db,
		conf: cfg,
		ox:   ox}, nil
}

func (r *ReMan) Register(reg *Registration) error {
	// registers the host with the cmdb
	_, err := r.ox.PutItem(&oxc.Item{
		Key:         reg.MachineId,
		Name:        reg.Hostname,
		Description: "Pilot registered remote host",
		Status:      0,
		Type:        "",
		Tag:         nil,
		Meta:        nil,
		Txt:         "",
		Attribute: map[string]interface{}{
			"CPU":          reg.CPUs,
			"OS":           reg.OS,
			"Total-Memory": reg.TotalMemory,
			"Platform":     reg.Platform,
			"Virtual":      reg.Virtual,
		},
	})
	return err
}

func (r *ReMan) Beat(host string) error {
	_, err := r.db.RunQuery(fmt.Sprintf("select rem_beat('%s')", host))
	if err != nil {
		return err
	}
	return nil
}

func (r *ReMan) GetHostStatus() ([]Host, error) {
	hosts := make([]Host, 0)
	result, err := r.db.RunQuery("select * from rem_get_conn_status()")
	if err != nil {
		return nil, fmt.Errorf("cannot get host status '%s'", err)
	}
	for _, row := range result.Rows {
		conn, err2 := strconv.ParseBool(row[1])
		if err2 != nil {
			fmt.Printf("cannot parse 'connected', value was '%s'", row[1])
		}
		hosts = append(hosts, Host{
			Name:      row[0],
			Customer:  "-",
			Region:    "-",
			Location:  "-",
			Connected: conn,
			LastSeen:  row[2],
		})
	}
	return hosts, nil
}
