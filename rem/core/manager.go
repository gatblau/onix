package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/oxc"
	"log"
	"strconv"
	"strings"
	"time"
)

// ReMan remote service manager API
type ReMan struct {
	conf *Conf
	db   *Db
	ox   *oxc.Client
}

func NewReMan() (*ReMan, error) {
	cfg := NewConf()
	db, err := NewDb(cfg.getDbHost(), cfg.getDbPort(), cfg.getDbName(), cfg.getDbUser(), cfg.getDbPwd())
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
	result, err := r.ox.PutItem(&oxc.Item{
		Key:         reg.MachineId,
		Name:        reg.Hostname,
		Description: "Pilot registered remote host",
		Status:      0,
		Type:        "U_HOST", // universal model host
		Tag:         nil,
		Meta:        nil,
		Txt:         "",
		Attribute: map[string]interface{}{
			"CPU":      reg.CPUs,
			"OS":       reg.OS,
			"MEMORY":   reg.TotalMemory,
			"PLATFORM": reg.Platform,
			"VIRTUAL":  reg.Virtual,
		},
	})
	// business error?
	if result != nil && result.Error {
		// return it
		return fmt.Errorf(result.Message)
	}
	// otherwise return technical error or nil if successful
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
			Id:        row[0],
			Customer:  "-",
			Region:    "-",
			Location:  "-",
			Connected: conn,
			Since:     row[2],
		})
	}
	return hosts, nil
}

func (r *ReMan) GetAdmissions() ([]Admission, error) {
	admissions := make([]Admission, 0)
	result, err := r.db.RunQuery("select * from rem_get_admissions(NULL)")
	if err != nil {
		return nil, fmt.Errorf("cannot get host status '%s'", err)
	}
	for _, row := range result.Rows {
		active, err2 := strconv.ParseBool(row[1])
		if err2 != nil {
			fmt.Printf("cannot parse 'active', value was '%s'", row[1])
		}
		admissions = append(admissions, Admission{
			Key:    row[0],
			Active: active,
			Tag:    row[2],
		})
	}
	return admissions, nil
}

func (r *ReMan) SetAdmission(admission *Admission) error {
	query := fmt.Sprintf("select rem_set_admission('%s', %s, '%s')", admission.Key, strconv.FormatBool(admission.Active), toTextArray(admission.Tag))
	_, err := r.db.RunCommand([]string{query})
	return err
}

// Authenticate authenticate a pilot based on its time stamp and machine Id admission status
func (r *ReMan) Authenticate(token string) bool {
	value, err := base64.StdEncoding.DecodeString(reverse(token))
	if err != nil {
		log.Printf("error decoding authentication token '%s': %s\n", token, err)
		return false
	}
	str := string(value)
	parts := strings.Split(str, "|")
	tokenTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Printf("error parsing authentication token: %s\naccess will be denied\n", err)
		return false
	}
	timeOk := (time.Now().Unix() - tokenTime) < (5 * 60)
	hostId := parts[0]
	if !timeOk {
		log.Printf("authentication failed for Machine Id='%s': token has expired\n", hostId)
		return false
	}
	result, err := r.db.RunQuery(fmt.Sprintf("select rem_is_admitted('%s')", hostId))
	if err != nil {
		fmt.Printf("authentication failed for Machine Id='%s': cannot query admission table: %s\n", hostId, err)
		return false
	}
	admitted, err := strconv.ParseBool(result.Rows[0][0])
	if err != nil {
		log.Printf("authentication failed for Machine Id='%s': cannot parse admission flag - %s\n", hostId, err)
		return false
	}
	if !admitted {
		// log an authentication error
		log.Printf("authentication failed for Machine Id='%s', host has not been admitted to service\n", hostId)
	}
	return admitted
}

func (r *ReMan) RecordConnStatus(interval int) error {
	_, err := r.db.RunCommand([]string{fmt.Sprintf("select rem_record_conn_status('%d secs')", interval)})
	return err
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
