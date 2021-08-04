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

func (r *ReMan) GetCommandValue(fxKey string) (*CmdValue, error) {
	item, err := r.ox.GetItem(&oxc.Item{Key: fxKey})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve function specification from Onix: %s", err)
	}
	input, err := getInputFromMap(item.Meta)
	if err != nil {
		return nil, fmt.Errorf("cannot get input from map: %s", err)
	}
	return &CmdValue{
		Function:      item.GetStringAttr("FX"),
		Package:       item.GetStringAttr("PACKAGE"),
		User:          item.GetStringAttr("USER"),
		Pwd:           item.GetStringAttr("PWD"),
		Verbose:       item.GetBoolAttr("VERBOSE"),
		Containerised: item.GetBoolAttr("CONTAINERISED"),
		Input:         input,
	}, nil
}

func (r *ReMan) Beat(machineId string) (jobId int64, fxKey string, fxVersion int64, err error) {
	rows, err := r.db.Query("select * from pilotctl_beat($1)", machineId)
	if err != nil {
		return -1, "", -1, err
	}
	for rows.Next() {
		err = rows.Scan(&jobId, &fxKey, &fxVersion)
		if err != nil {
			return 0, "", 0, err
		}
	}
	// returns the next job to execute for the machine Id or -1 if no job is available
	return jobId, fxKey, fxVersion, nil
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
		err := rows.Scan(&id, &connected, &since, &customer, &region, &location)
		if err != nil {
			return nil, err
		}
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

func (r *ReMan) SetAdmission(admission *Admission) error {
	if len(admission.MachineId) == 0 {
		return fmt.Errorf("machine Id is missing")
	}
	return r.db.RunCommand("select pilotctl_set_admission($1, $2, $3, $4, $5, $6)",
		admission.MachineId,
		admission.OrgGroup,
		admission.Org,
		admission.Area,
		admission.Location,
		admission.Tag)
}

// Authenticate a pilot based on its time stamp and machine Id admission status
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
	machineId := parts[0]
	if !timeOk {
		log.Printf("authentication failed for Machine Id='%s': token has expired\n", machineId)
		return false
	}
	rows, err := r.db.Query("select * from pilotctl_is_admitted($1)", machineId)
	if err != nil {
		fmt.Printf("authentication failed for Machine Id='%s': cannot query admission table: %s\n", machineId, err)
		return false
	}
	var admitted bool
	for rows.Next() {
		err = rows.Scan(&admitted)
		if err != nil {
			log.Printf("authentication failed for Machine Id='%s': %s\n", machineId, err)
			return false
		}
		break
	}
	if !admitted {
		// log an authentication error
		log.Printf("authentication failed for Machine Id='%s', host has not been admitted to service\n", machineId)
	}
	return admitted
}

func (r *ReMan) RecordConnStatus(interval int) error {
	return r.db.RunCommand("select pilotctl_record_conn_status($1)", time.Duration(interval))
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

// PutCommand put the command in the Onix database
func (r *ReMan) PutCommand(cmd *Cmd) error {
	var meta map[string]interface{}
	inputBytes, err := json.Marshal(cmd.Input)
	if err != nil {
		return fmt.Errorf("cannot marshal command input: %s", err)
	}
	err = json.Unmarshal(inputBytes, &meta)
	if err != nil {
		return fmt.Errorf("cannot unmarshal input bytes: %s", err)
	}
	result, err := r.ox.PutItem(&oxc.Item{
		Key:         fmt.Sprintf("ART_FX_%s", strings.Replace(cmd.Key, " ", "", -1)),
		Name:        cmd.Key,
		Description: cmd.Description,
		Type:        "ART_FX",
		Meta:        meta,
		Attribute: map[string]interface{}{
			"PACKAGE":       cmd.Package,
			"FX":            cmd.Function,
			"USER":          cmd.User,
			"PWD":           cmd.Pwd,
			"VERBOSE":       cmd.Verbose,
			"CONTAINERISED": cmd.Containerised,
		},
	})
	if result != nil && result.Error {
		return fmt.Errorf("cannot set command in Onix: %s\n", result.Message)
	}
	if err != nil {
		return fmt.Errorf("cannot set command in Onix: %s\n", err)
	}
	return nil
}

func (r *ReMan) GetAllCommands() ([]Cmd, error) {
	items, err := r.ox.GetItemsByType("ART_FX")
	if err != nil {
		return nil, fmt.Errorf("cannot get commands from Onix: %s", err)
	}
	cmds := make([]Cmd, 0)
	for _, item := range items.Values {
		input, err := getInputFromMap(item.Meta)
		if err != nil {
			return nil, fmt.Errorf("cannot transform input map: %s", err)
		}
		cmds = append(cmds, Cmd{
			Key:           item.Key,
			Description:   item.Description,
			Package:       item.GetStringAttr("PACKAGE"),
			Function:      item.GetStringAttr("FX"),
			User:          item.GetStringAttr("USER"),
			Pwd:           item.GetStringAttr("PWD"),
			Verbose:       item.GetBoolAttr("VERBOSE"),
			Containerised: item.GetBoolAttr("CONTAINERISED"),
			Input:         input,
		})

	}
	return cmds, nil
}

func (r *ReMan) GetCommand(cmdName string) (*Cmd, error) {
	item, err := r.ox.GetItem(&oxc.Item{Key: cmdName})
	if err != nil {
		return nil, fmt.Errorf("cannot get command with key '%s' from Onix: %s", cmdName, err)
	}
	input, err := getInputFromMap(item.Meta)
	if err != nil {
		return nil, fmt.Errorf("cannot transform input map: %s", err)
	}
	return &Cmd{
		Key:           item.Key,
		Description:   item.Description,
		Package:       item.GetStringAttr("PACKAGE"),
		Function:      item.GetStringAttr("FX"),
		User:          item.GetStringAttr("USER"),
		Pwd:           item.GetStringAttr("PWD"),
		Verbose:       item.GetBoolAttr("VERBOSE"),
		Containerised: item.GetBoolAttr("CONTAINERISED"),
		Input:         input,
	}, nil
}

func (r *ReMan) CompleteJob(status *Result) error {
	return r.db.RunCommand("select pilotctl_complete_job($1, $2, $3)", status.JobId, status.Log, !status.Success)
}

func (r *ReMan) GetAreas(orgGroup string) ([]Area, error) {
	items, err := r.ox.GetChildrenByType(&oxc.Item{Key: orgGroup}, "U_AREA")
	if err != nil {
		return nil, err
	}
	var areas []Area
	for _, item := range items.Values {
		areas = append(areas, Area{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return areas, nil
}

func (r *ReMan) GetOrgs(orgGroup string) ([]Org, error) {
	items, err := r.ox.GetChildrenByType(&oxc.Item{Key: orgGroup}, "U_ORG")
	if err != nil {
		return nil, err
	}
	var orgs []Org
	for _, item := range items.Values {
		orgs = append(orgs, Org{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return orgs, nil
}

func (r *ReMan) GetLocations(area string) ([]Location, error) {
	items, err := r.ox.GetChildrenByType(&oxc.Item{Key: area}, "U_LOCATION")
	if err != nil {
		return nil, err
	}
	var orgs []Location
	for _, item := range items.Values {
		orgs = append(orgs, Location{
			Key:  item.Key,
			Name: item.Name,
		})
	}
	return orgs, nil
}

func (r *ReMan) GetOrgGroups() ([]Org, error) {
	items, err := r.ox.GetItemsByType("U_ORG_GROUP")
	if err != nil {
		return nil, err
	}
	var orgs []Org
	for _, item := range items.Values {
		orgs = append(orgs, Org{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return orgs, nil
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
