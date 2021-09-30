package core

/*
  Onix Pilot Host Control
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
	oxc "github.com/gatblau/onix/client"
	. "github.com/gatblau/onix/pilotctl/types"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// API backend services API
type API struct {
	conf *Conf
	db   *Db
	ox   *oxc.Client
	// host information
	hostUUID string
	hostname string
	hostIP   string
	// event publisher
	pub *EventPublisher
	// user access controls
	ACL *oxc.Controls
}

func NewAPI(cfg *Conf) (*API, error) {
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
	return &API{
		db:   db,
		conf: cfg,
		ox:   ox,
		pub:  NewEventPublisher(),
	}, nil
}

func (r *API) PingInterval() time.Duration {
	return r.conf.PingIntervalSecs()
}

func (r *API) PublishEvents(events *Events) {
	r.pub.Publish(events)
}

func (r *API) Register(reg *RegistrationRequest) (*RegistrationResponse, error) {
	// registers the host with the cmdb
	result, err := r.ox.PutItem(&oxc.Item{
		Key:         strings.ToUpper(fmt.Sprintf("HOST:%s:%s", reg.MachineId, reg.Hostname)),
		Name:        strings.ToUpper(fmt.Sprintf("HOST_%s", strings.Replace(reg.Hostname, " ", "_", -1))),
		Description: "Pilot registered remote host",
		Status:      0,
		Type:        "U_HOST", // universal model host
		Tag:         nil,
		Meta:        nil,
		Txt:         "",
		// TODO: add mac address attribute
		Attribute: map[string]interface{}{
			"CPU":        reg.CPUs,
			"OS":         reg.OS,
			"MEMORY":     reg.TotalMemory,
			"PLATFORM":   reg.Platform,
			"VIRTUAL":    reg.Virtual,
			"IP":         reg.HostIP,
			"MACHINE-ID": reg.MachineId,
			"HOSTNAME":   reg.Hostname,
		},
	})
	// business error?
	if result != nil && result.Error {
		// return it
		return nil, fmt.Errorf(result.Message)
	}
	// otherwise, return technical error or nil if successful
	return &RegistrationResponse{
		Operation: result.Operation,
	}, err
}

func (r *API) GetCommandValue(fxKey string) (*CmdInfo, error) {
	item, err := r.ox.GetItem(&oxc.Item{Key: fxKey})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve function specification from Onix: %s", err)
	}
	input, err := getInputFromMap(item.Meta)
	if err != nil {
		return nil, fmt.Errorf("cannot get input from map: %s", err)
	}
	return &CmdInfo{
		Function:      item.GetStringAttr("FX"),
		Package:       item.GetStringAttr("PACKAGE"),
		User:          item.GetStringAttr("USER"),
		Pwd:           item.GetStringAttr("PWD"),
		Verbose:       item.GetBoolAttr("VERBOSE"),
		Containerised: item.GetBoolAttr("CONTAINERISED"),
		Input:         input,
	}, nil
}

func (r *API) Ping() (jobId int64, fxKey string, fxVersion int64, err error) {
	rows, err := r.db.Query("select * from pilotctl_beat($1)", r.hostUUID)
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

// GetHosts get a list of hosts filtered by
// oGroup: organisation group key
// or: organisation key
// ar: area key
// loc: location key
func (r *API) GetHosts(oGroup, or, ar, loc string, label []string) ([]Host, error) {
	hosts := make([]Host, 0)
	pingInterval := r.PingInterval().Seconds()
	rows, err := r.db.Query("select * from pilotctl_get_host($1, $2, $3, $4, $5, $6)", fmt.Sprintf("%.0f secs", pingInterval+5), oGroup, or, ar, loc, label)
	if err != nil {
		return nil, fmt.Errorf("cannot get hosts: %s\n", err)
	}
	var (
		uuId      string
		connected bool
		lastSeen  sql.NullTime
		orgGroup  sql.NullString
		org       sql.NullString
		area      sql.NullString
		location  sql.NullString
		inService bool
		labels    []string
	)
	for rows.Next() {
		err = rows.Scan(&uuId, &connected, &lastSeen, &orgGroup, &org, &area, &location, &inService, &labels)
		if err != nil {
			return nil, err
		}
		var (
			tt        int64
			since     int
			sinceType string
		)
		if lastSeen.Valid {
			tt = lastSeen.Time.UnixNano()
			since, sinceType, err = toElapsedValues(lastSeen.Time.Format(time.RFC850))
			if err != nil {
				return nil, err
			}
		}
		hosts = append(hosts, Host{
			HostUUID:  uuId,
			OrgGroup:  orgGroup.String,
			Org:       org.String,
			Area:      area.String,
			Location:  location.String,
			Connected: connected,
			LastSeen:  tt,
			Since:     since,
			SinceType: sinceType,
			Label:     labels,
		})
	}
	return hosts, rows.Err()
}

func (r *API) SetAdmission(admission Admission) error {
	if len(admission.HostUUID) == 0 {
		return fmt.Errorf("host UUID is missing")
	}
	return r.db.RunCommand("select pilotctl_set_admission($1, $2, $3, $4, $5, $6)",
		admission.HostUUID,
		admission.OrgGroup,
		admission.Org,
		admission.Area,
		admission.Location,
		admission.Label)
}

// AuthenticatePilot a pilot based on its time stamp and machine Id admission status
func (r *API) AuthenticatePilot(token string) bool {
	if len(token) == 0 {
		log.Println("authentication token is required and not provided")
		return false
	}
	value, err := base64.StdEncoding.DecodeString(reverse(token))
	if err != nil {
		log.Printf("error decoding authentication token '%s': %s\n", token, err)
		return false
	}
	str := string(value)
	// token is: hostUUID (0) | hostIP (1) | hostName (2) | timestamp (3)
	parts := strings.Split(str, "|")
	tokenTime, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Printf("error parsing authentication token: %s\naccess will be denied\n", err)
		return false
	}
	timeOk := (time.Now().Unix() - tokenTime) < (5 * 60)
	hostUUId := parts[0]
	if !timeOk {
		log.Printf("authentication failed for Host UUID='%s': token has expired\n", hostUUId)
		return false
	}
	rows, err := r.db.Query("select * from pilotctl_is_admitted($1)", hostUUId)
	var hostname, hostIP string
	if err != nil {
		hostIP = parts[1]
		hostname = parts[2]
		fmt.Printf("authentication failed for Host UUID='%s': cannot query admission table: %s\n"+
			"additional info: host IP = '%s', hostname = '%s'\n", hostUUId, err, hostIP, hostname)
		return false
	}
	var admitted bool
	for rows.Next() {
		err = rows.Scan(&admitted)
		if err != nil {
			log.Printf("authentication failed for Host UUID='%s': %s\n", hostUUId, err)
			return false
		}
		break
	}

	if !admitted {
		// log an authentication error
		log.Printf("authentication failed for Host UUID='%s', host has not been admitted to service\n", hostUUId)
	}

	// captures token information for remote host
	r.hostUUID = hostUUId
	r.hostname = hostname
	r.hostIP = hostIP

	// returns authentication status
	return admitted
}

func (r *API) AuthenticateUser(request http.Request) bool {
	// get the credentials from the request header
	user, pwd := getCredentials(request)
	// validate the credentials and retrieve user access controls
	acl, err := r.ox.Login(&oxc.Login{
		Email:    user,
		Password: pwd,
	})
	if err != nil {
		fmt.Printf("WARNING: user authentication failed, %s\n", err)
		return false
	}
	// update the control list
	r.ACL = acl
	return true
}

// GetPackages get a list of packages in the backing Artisan registry
func (r *API) GetPackages() ([]PackageInfo, error) {
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
	// removes protocol prefix from URI
	regDomain := r.conf.getArtRegUri()[strings.Index(r.conf.getArtRegUri(), "//")+2:]
	// removes forward slash at the end if it finds one
	if strings.HasSuffix(regDomain, "/") {
		regDomain = regDomain[0 : len(regDomain)-1]
	}
	result := make([]PackageInfo, 0)
	for _, repo := range repos {
		pack := PackageInfo{
			Name: fmt.Sprintf("%s/%s", regDomain, repo.Repository),
			Tags: []TagInfo{},
		}
		for _, p := range repo.Packages {
			for _, tag := range p.Tags {
				pack.Tags = append(pack.Tags, TagInfo{
					Id:      p.Id,
					Name:    tag,
					Ref:     p.FileRef,
					Created: p.Created,
					Type:    p.Type,
					Size:    p.Size,
				})
			}
		}
		result = append(result, pack)
	}
	return result, nil
}

func (r *API) GetPackageAPI(name string) ([]*data.FxInfo, error) {
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
func (r *API) PutCommand(cmd *Cmd) error {
	var meta map[string]interface{}
	m := make(map[string]interface{}, 0)
	m["input"] = cmd.Input
	inputBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("cannot marshal command input: %s", err)
	}
	err = json.Unmarshal(inputBytes, &meta)
	if err != nil {
		return fmt.Errorf("cannot unmarshal input bytes: %s", err)
	}
	result, err := r.ox.PutItem(&oxc.Item{
		Key:         strings.Replace(cmd.Key, " ", "", -1),
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

func (r *API) GetAllCommands() ([]Cmd, error) {
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

func (r *API) GetCommand(cmdName string) (*Cmd, error) {
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

func (r *API) DeleteCommand(cmdName string) (string, error) {
	result, err := r.ox.DeleteItem(&oxc.Item{Key: cmdName})
	if err != nil {
		return "", fmt.Errorf("cannot delete command with key '%s' from Onix: %s", cmdName, err)
	}
	return result.Operation, nil
}

func (r *API) CompleteJob(status *JobResult) error {
	logMsg := status.Log
	// if there was a failure, and we have an error message, add it to the log
	if !status.Success && len(status.Err) > 0 {
		logMsg = fmt.Sprintf("%s !!! ERROR: %s\n", logMsg, status.Err)
	}
	return r.db.RunCommand("select pilotctl_complete_job($1, $2, $3)", status.JobId, logMsg, !status.Success)
}

func (r *API) GetAreas(orgGroup string) ([]Area, error) {
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

func (r *API) GetOrgs(orgGroup string) ([]Org, error) {
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

func (r *API) GetLocations(area string) ([]Location, error) {
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

func (r *API) GetOrgGroups() ([]Org, error) {
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

func (r *API) CreateJobBatch(info JobBatchInfo) (int64, error) {
	if len(info.HostUUID) == 0 {
		return -1, fmt.Errorf("host UUID is missing\n")
	}
	if len(info.FxKey) == 0 {
		return -1, fmt.Errorf("fx is missing\n")
	}
	// create a job batch identifier
	rows, err := r.db.Query("select * from pilotctl_create_job_batch($1, $2, $3, $4)", info.Name, info.Notes, "???", info.Label)
	if err != nil {
		return -1, fmt.Errorf("cannot create job batch: %s\n", err)
	}
	var batchId int64 = -1
	for rows.Next() {
		rows.Scan(&batchId)
	}
	if batchId == -1 {
		return -1, fmt.Errorf("cannot retrieve job batch Id\n")
	}
	// add jobs to the batch using the batch ID
	var returnError error
	for _, uuid := range info.HostUUID {
		err = r.db.RunCommand("select pilotctl_create_job($1, $2, $3, $4)", batchId, uuid, info.FxKey, info.FxVersion)
		// if there is an error creating the job
		if err != nil {
			if returnError == nil {
				returnError = err
			}
			// accumulates the error and continue with the next job
			returnError = fmt.Errorf("can't create job: %s\n", returnError)
		}
	}
	// return any job creation error
	return batchId, returnError
}

func (r *API) GetJobs(oGroup, or, ar, loc string, batchId *int64) ([]Job, error) {
	jobs := make([]Job, 0)
	rows, err := r.db.Query("select * from pilotctl_get_jobs($1, $2, $3, $4, $5)", oGroup, or, ar, loc, batchId)
	if err != nil {
		return nil, fmt.Errorf("cannot get jobs: %s\n", err)
	}
	var (
		id         int64
		hostUUID   string
		jobBatchId int64
		fxKey      string
		fxVersion  int64
		created    sql.NullTime
		started    sql.NullTime
		completed  sql.NullTime
		log        sql.NullString
		e          sql.NullBool
		orgGroup   sql.NullString
		org        sql.NullString
		area       sql.NullString
		location   sql.NullString
		tag        []string
	)
	for rows.Next() {
		err = rows.Scan(&id, &hostUUID, &jobBatchId, &fxKey, &fxVersion, &created, &started, &completed, &log, &e, &orgGroup, &org, &area, &location, &tag)
		if err != nil {
			return nil, fmt.Errorf("cannot scan job row: %e\n", err)
		}
		jobs = append(jobs, Job{
			Id:         id,
			HostUUID:   hostUUID,
			JobBatchId: jobBatchId,
			FxKey:      fxKey,
			FxVersion:  fxVersion,
			Created:    timeF(created),
			Started:    timeF(started),
			Completed:  timeF(completed),
			Log:        stringF(log),
			Error:      boolF(e),
			OrgGroup:   orgGroup.String,
			Org:        org.String,
			Area:       area.String,
			Location:   location.String,
			Tag:        tag,
		})
	}
	return jobs, rows.Err()
}

func (r *API) GetHost(uuid string) (*Host, error) {
	rows, err := r.db.Query("select * from pilotctl_get_host($1)", uuid)
	if err != nil {
		return nil, fmt.Errorf("cannot get host data: %s\n", err)
	}
	var (
		orgGroup string
		org      string
		area     string
		location string
		labels   []string
	)
	for rows.Next() {
		err = rows.Scan(&orgGroup, &org, &area, &location, &labels)
		if err != nil {
			return nil, fmt.Errorf("cannot scan host row: %e\n", err)
		}
		return &Host{HostUUID: uuid, OrgGroup: orgGroup, Org: org, Area: area, Location: location, Label: labels}, nil
	}
	return nil, fmt.Errorf("host uuid '%s' cannot be found in data source\n", uuid)
}

func (r *API) GetJobBatches(name, owner *string, from, to *time.Time, label *[]string) ([]JobBatch, error) {
	batches := make([]JobBatch, 0)
	rows, err := r.db.Query("select * from pilotctl_get_job_batches($1, $2, $3, $4, $5)", name, from, to, label, owner)
	if err != nil {
		return nil, fmt.Errorf("cannot get job batches: %s\n", err)
	}
	var (
		id      int64
		name2   string
		notes   string
		created sql.NullTime
		owner2  string
		labels  []string
		jobs    int
	)
	for rows.Next() {
		err = rows.Scan(&id, &name2, &notes, &labels, &created, &owner2, &jobs)
		if err != nil {
			return nil, fmt.Errorf("cannot scan job batch row: %e\n", err)
		}
		batches = append(batches, JobBatch{
			BatchId: id,
			Name:    name2,
			Notes:   notes,
			Label:   labels,
			Owner:   owner2,
			Jobs:    jobs,
			Created: created.Time,
		})
	}
	return batches, rows.Err()
}

func (r *API) Augment(events *Events) (*Events, error) {
	// as all events come from the same host get the host UUID from the first event
	hostUUID := events.Events[0].HostUUID
	// retrieve host information
	host, err := r.GetHost(hostUUID)
	if err != nil {
		return events, err
	}
	// add info to events
	var result []Event
	for _, event := range events.Events {
		ev := Event{
			EventID:           event.EventID,
			Client:            event.Client,
			Hostname:          event.Hostname,
			HostUUID:          event.HostUUID,
			MachineId:         event.MachineId,
			HostAddress:       event.HostAddress,
			Organisation:      host.Org,
			OrganisationGroup: host.OrgGroup,
			Area:              host.Area,
			Location:          host.Location,
			Facility:          event.Facility,
			Priority:          event.Priority,
			Severity:          event.Severity,
			Time:              event.Time,
			TLSPeer:           event.TLSPeer,
			BootTime:          event.BootTime,
			Content:           event.Content,
			Tag:               event.Tag,
			MacAddress:        event.MacAddress,
			HostLabel:         host.Label,
		}
		result = append(result, ev)
	}
	return &Events{Events: result}, nil
}

// Login if the user is authenticated returns a list of access controls otherwise an error is returned
func (r *API) Login(username string) ([]oxc.AccessControl, error) {
	user, err := r.ox.GetUser(&oxc.User{Key: username})
	if err != nil {
		return nil, fmt.Errorf("login failed for user '%s': %s\n", username, err)
	}
	// filter the controls using the realm
	var (
		controls = user.Controls()
		result   []oxc.AccessControl
	)
	for _, control := range controls {
		if control.Realm == "pilotcl" {
			result = append(result, control)
		}
	}
	return result, nil
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func timeF(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format(time.RFC822Z)
	}
	return ""
}

func stringF(t sql.NullString) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func boolF(t sql.NullBool) bool {
	if t.Valid {
		return t.Bool
	}
	return false
}

// getUser retrieve the username from the basic authentication token in the http request
func getCredentials(r http.Request) (user, pwd string) {
	// get the token from the authorization header
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", ""
	}
	decoded, err := base64.StdEncoding.DecodeString(token[len("Basic "):])
	if err != nil {
		log.Printf("WARNING: failed to decode Authorization header: %s, cannot retrieve username\n", err)
		return "", ""
	}
	parts := strings.Split(string(decoded[:]), ":")
	if len(parts) != 2 {
		log.Printf("WARNING: failed to parse Authorization header: invalid format '%s' assuming a Basic Authentication Token\n", string(decoded[:]))
		return "", ""
	}
	// retrieve the username part (i.e. #0: username:password => 0:1)
	return parts[0], parts[1]
}
