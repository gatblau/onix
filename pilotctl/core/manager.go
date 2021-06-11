package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/oxc"
	"github.com/jackc/pgtype"
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
	if err != nil {
		return nil, err
	}
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
	_, err := r.db.Query("select pilotctl_beat($1)", host)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReMan) GetHostStatus() ([]Host, error) {
	hosts := make([]Host, 0)
	rows, err := r.db.Query("select * from pilotctl_get_conn_status()")
	if err != nil {
		return nil, fmt.Errorf("cannot get host status '%s'", err)
	}
	var (
		id        string
		connected bool
		since     time.Time
		customer  sql.NullString
		region    sql.NullString
		location  sql.NullString
	)
	for rows.Next() {
		rows.Scan(&id, &connected, &since, &customer, &region, &location)
		hosts = append(hosts, Host{
			Id:        id,
			Customer:  customer.String,
			Region:    region.String,
			Location:  location.String,
			Connected: connected,
			Since:     toTime(since.UnixNano()),
		})
	}
	return hosts, rows.Err()
}

func (r *ReMan) GetAdmissions() ([]Admission, error) {
	admissions := make([]Admission, 0)
	rows, err := r.db.Query("select * from pilotctl_get_admissions($1)", nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get host status '%s'", err)
	}
	var (
		key    string
		active bool
		tag    []string
	)
	for rows.Next() {
		rows.Scan(&key, &active, &tag)
		admissions = append(admissions, Admission{
			Key:    key,
			Active: active,
			Tag:    tag,
		})
	}
	return admissions, rows.Err()
}

func (r *ReMan) SetAdmission(admission *Admission) error {
	return r.db.RunCommand("select pilotctl_set_admission($1, $2, $3)", admission.Key, admission.Active, admission.Tag)
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
	rows, err := r.db.Query("select pilotctl_is_admitted($1)", hostId)
	if err != nil {
		fmt.Printf("authentication failed for Machine Id='%s': cannot query admission table: %s\n", hostId, err)
		return false
	}
	var admitted bool
	err = rows.Scan(&admitted)
	if err != nil {
		log.Printf("authentication failed for Machine Id='%s': %s\n", hostId, err)
		return false
	}
	if !admitted {
		// log an authentication error
		log.Printf("authentication failed for Machine Id='%s', host has not been admitted to service\n", hostId)
	}
	return admitted
}

func (r *ReMan) RecordConnStatus(interval int) error {
	return r.db.RunCommand("select pilotctl_record_conn_status($1)", interval)
}

// GetPackages get a list of packages in the backing Artisan registry
func (r *ReMan) GetPackages() ([]string, error) {
	// the URI to connect to the Artisan registry
	uri := fmt.Sprintf("%s/repository", r.conf.getArtRegUri())
	bytes, err := makeRequest(uri, "GET", r.conf.getArtRegUser(), r.conf.getArtRegPwd(), nil)
	if err != nil {
		return nil, err
	}
	var repos []registry.Repository
	err = json.Unmarshal(bytes, &repos)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	// removes protocol prefix from URI
	regDomain := r.conf.getArtRegUri()[strings.Index(r.conf.getArtRegUri(), "//")+2:]
	for _, repo := range repos {
		for _, p := range repo.Packages {
			for _, tag := range p.Tags {
				// append constructed package name
				result = append(result, fmt.Sprintf("%s/%s:%s", regDomain, repo.Repository, tag))
			}
		}
	}
	return result, nil
}

func (r *ReMan) GetPackageAPI(name string) ([]*data.FxInfo, error) {
	n, err := core.ParseName(name)
	if err != nil {
		return nil, err
	}
	// the URI to connect to the Artisan registry
	uri := fmt.Sprintf("%s/package/manifest/%s/%s/%s", r.conf.getArtRegUri(), n.Group, n.Name, n.Tag)
	bytes, err := makeRequest(uri, "GET", r.conf.getArtRegUser(), r.conf.getArtRegPwd(), nil)
	if err != nil {
		return nil, err
	}
	var manif data.Manifest
	err = json.Unmarshal(bytes, &manif)
	if err != nil {
		return nil, err
	}
	if manif.Functions == nil {
		return make([]*data.FxInfo, 0), nil
	}
	return manif.Functions, nil
}

func (r *ReMan) SetCommand(cmd *Cmd) error {
	inputHS := toHStoreString(cmd.Input)
	return r.db.RunCommand("select pilotctl_set_command($1, $2, $3, $4, $5)", cmd.Name, cmd.Description, cmd.Package, cmd.Function, inputHS)
}

func (r *ReMan) GetAllCommands() ([]Cmd, error) {
	rows, err := r.db.Query("select id, name, description, package, fx, input from pilotctl_get_command($1)", nil)
	if err != nil {
		return nil, fmt.Errorf("cannot execute query: %s", err)
	}
	cmds := make([]Cmd, 0)
	var (
		id          int64
		name        string
		description sql.NullString
		pack        string
		fx          string
		input       pgtype.Hstore
	)
	for rows.Next() {
		rows.Scan(&id, &name, &description, &pack, &fx, &input)
		cmds = append(cmds, Cmd{
			Id:          id,
			Name:        name,
			Description: description.String,
			Package:     pack,
			Function:    fx,
			Input:       toMap(input),
		})
	}
	return cmds, rows.Err()
}

func (r *ReMan) GetCommand(cmdName string) (*Cmd, error) {
	rows, err := r.db.Query("select id, name, description, package, fx, input from pilotctl_get_command($1)", cmdName)
	if err != nil {
		return nil, fmt.Errorf("cannot execute query: %s", err)
	}
	var (
		id          int64
		name        string
		description sql.NullString
		pack        string
		fx          string
		input       pgtype.Hstore
	)
	return &Cmd{
		Id:          id,
		Name:        name,
		Description: description.String,
		Package:     pack,
		Function:    fx,
		Input:       toMap(input),
	}, rows.Err()
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func toMap(hs pgtype.Hstore) map[string]string {
	m := make(map[string]string)
	for k, v := range hs.Map {
		m[k] = v.String
	}
	return m
}
