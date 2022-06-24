/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/doorman/db"
	"github.com/gatblau/onix/artisan/release"
	"github.com/gatblau/onix/oxlib/oxc"
	"github.com/gatblau/onix/oxlib/resx"
	"github.com/minio/minio-go/v7/pkg/notification"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/doorman/types"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	util "github.com/gatblau/onix/oxlib/httpserver"
)

const (
	DoormanLogging = "DOORMAN_LOGGING"
	ArtSpecFxType  = "ART_SPEC_FX"
	UCatalogueType = "U_CATALOGUE"
)

type Process struct {
	serviceId  string
	bucketName string
	folderName string
	tmp        string
	log        *bytes.Buffer
	db         db.Database
	reg        *registry.LocalRegistry
	spec       *release.Spec
	jobNo      string
	cmdLog     string
	pipe       *types.Pipeline
	ox         *oxc.Client
}

func NewProcess(serviceId, bucketPath, folderName, artHome string) (Processor, error) {
	p := new(Process)
	p.serviceId = serviceId
	p.bucketName = bucketPath
	p.folderName = folderName
	p.log = new(bytes.Buffer)
	p.db = *db.New()
	p.reg = registry.NewLocalRegistry(artHome)
	// if an Onix Web API is defined
	if len(os.Getenv(OxWapiUri)) > 0 {
		oxClient, err := newOxClient()
		if err != nil {
			return nil, err
		}
		p.ox = oxClient
	}
	return p, nil
}

func newOxClient() (*oxc.Client, error) {
	uri, err := GetOxWapiUri()
	if err != nil {
		return nil, err
	}
	user, err := GetOxWapiUser()
	if err != nil {
		return nil, err
	}
	pwd, err := GetOxWapiPwd()
	if err != nil {
		return nil, err
	}
	skip, err := GetOxWapiInsecureSkipVerify()
	if err != nil {
		return nil, err
	}
	oxcfg := &oxc.ClientConf{
		BaseURI:            uri,
		Username:           user,
		Password:           pwd,
		InsecureSkipVerify: skip,
	}
	oxcfg.SetAuthMode("basic")
	ox, err := oxc.NewClient(oxcfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create onix http client: %s", err)
	}
	return ox, err
}

// Info logger
func (p *Process) Info(format string, a ...interface{}) {
	// align format with ART logger
	format = fmt.Sprintf("%s ART INFO %s\n", time.Now().Format("2006/01/02 15:04:05.999999"), format)
	msg := fmt.Sprintf(format, a...)
	if len(os.Getenv(DoormanLogging)) > 0 {
		fmt.Println(msg)
	}
	_, err := p.log.WriteString(msg)
	if err != nil {
		fmt.Printf("cannot log INFO: %s\n", err)
	}
}

// Error logger
func (p *Process) Error(format string, a ...interface{}) error {
	// align format with ART logger
	format = fmt.Sprintf("%s ART ERROR %s", time.Now().Format("2006/01/02 15:04:05.999999"), format)
	msg := fmt.Sprintf(format, a...)
	if len(os.Getenv(DoormanLogging)) > 0 {
		fmt.Println(msg)
	}
	p.log.WriteString(fmt.Sprintf("%s\n", msg))
	return fmt.Errorf(msg)
}

// Warn logger
func (p *Process) Warn(format string, a ...interface{}) {
	// align format with ART logger
	format = fmt.Sprintf("%s ART WARN %s\n", time.Now().Format("2006/01/02 15:04:05.999999"), format)
	msg := fmt.Sprintf(format, a...)
	if len(os.Getenv(DoormanLogging)) > 0 {
		fmt.Println(msg)
	}
	p.log.WriteString(msg)
}

// Start processing a pipeline
func (p *Process) Start() {
	go p.run()
}

func (p *Process) run() {
	defer func() {
		// remove the artisan home once the processing is complete
		if err := os.RemoveAll(p.artHome()); err != nil {
			core.WarningLogger.Printf("cannot cleanup artisan home @ '%s': %s", p.artHome(), err)
		}
	}()
	p.Info("processing release Id=%s → %s/%s", p.serviceId, p.bucketName, p.folderName)
	pipes, err := p.db.MatchPipelines(p.serviceId, p.bucketName)
	if err != nil {
		e := p.Error("cannot retrieve pipelines for bucket Id='%s' and bucket name='%s': %s\n", p.serviceId, p.bucketName, err)
		fmt.Println(e)
		p.notify(err)
	}
	if len(pipes) == 0 {
		e := p.Error("no pipeline configuration found for release Id=%s and bucket name='%s': %s\n", p.serviceId, p.bucketName)
		fmt.Println(e)
		p.notify(err)
	}
	for i, pipe := range pipes {
		// set the current pipeline
		p.pipe = &pipes[i]
		// record the start of a new job and obtains a new job number
		jobNo, startTime, jobErr := StartJob(&pipe, p)
		if jobErr != nil {
			p.notify(jobErr)
		}
		p.jobNo = jobNo
		// process the pipeline
		err = p.Pipeline(&pipe)
		// if there was an error
		if err != nil {
			// record the job as failed passing the logs
			jobErr = FailJob(startTime, &pipe, p)
			if jobErr != nil {
				p.notify(jobErr)
			}
			p.notify(err)
		}
		// if no error, record the job as completed passing the logs
		jobErr = CompleteJob(startTime, &pipe, p)
		if jobErr != nil {
			p.notify(jobErr)
		}
		// notify success
		p.notify(nil)
	}
}

// Pipeline process a pipeline
func (p *Process) Pipeline(pipe *types.Pipeline) error {
	defer p.cleanup()
	// create a new temp folder for processing, uses the artisan home specified in the registry
	// in order to be able to parallelize processes, the artisan home must be different for each instance of Process
	tmp, err := core.NewTempDir(p.artHome())
	if err != nil {
		return p.Error("cannot create temporary folder: %s\n", err)
	}
	p.tmp = tmp
	// find inbound route
	for _, inRoute := range pipe.InboundRoutes {
		// find the inbound route matching the bucket Id
		if inRoute.ServiceId == p.serviceId {
			// process the inbound route and run any specified commands
			err = p.InboundRoute(pipe, inRoute)
			if err != nil {
				return err
			}
			// process the outbound route(s)
			for _, outRoute := range pipe.OutboundRoutes {
				err = p.OutboundRoute(outRoute)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return p.BeforeComplete(pipe)
}

// InboundRoute process an inbound route
func (p *Process) InboundRoute(pipe *types.Pipeline, route types.InRoute) error {
	// download spec
	p.Info("downloading specification: started")
	spec, err := release.DownloadSpec(
		release.UpDownOptions{
			TargetUri:   fmt.Sprintf("%s/%s", route.ServiceHost, p.bucketPath()),
			TargetCreds: route.Creds(),
			LocalPath:   p.tmp,
		})
	if err != nil {
		return p.Error("cannot download specification: %s\n", err)
	}
	p.spec = spec
	p.Info("downloading specification: complete")
	// execute commands
	if len(pipe.Commands) > 0 {
		p.Info("running commands: started")
	}
	for _, command := range pipe.Commands {
		if err = p.Command(command); err != nil {
			return err
		}
	}
	if len(pipe.Commands) > 0 {
		p.Info("running commands: complete")
	}
	err = p.PreImport(route, err)
	if err != nil {
		return err
	}
	return p.ImportFiles()
}

func (p *Process) PreImport(route types.InRoute, err error) error {
	return nil
}

func (p *Process) Command(command types.Command) error {
	c := strings.ReplaceAll(command.Value, "${path}", p.tmp)
	p.Info("executing verification task: %s", c)
	out, exeErr := build.ExeAsync(c, ".", merge.NewEnVarFromSlice([]string{}), false)
	if exeErr != nil {
		return p.Error("execution failed: %s", exeErr)
	}
	// use the regex in the command definition to decide if the command execution failed based on the content of the output
	core.Debug("command: %s, error regex: %s", command.Name, command.ErrorRegex)
	core.Debug("Command Output:")
	core.Debug(out)
	matched, regexErr := regexp.MatchString(command.ErrorRegex, out)
	if regexErr != nil {
		return p.Error("invalid regex %s: %s", command.ErrorRegex, regexErr)
	}
	core.Debug("Regex matched: %v", matched)
	// if the regex matched return error and content of command output
	if matched {
		cmdErr := fmt.Sprintf("command %s failed:\n%s", command.Name, out)
		p.cmdLog = cmdErr
		// and should stop on error
		core.Debug("stop on error: %v", command.StopOnError)
		if command.StopOnError {
			// stops and return
			return p.Error(cmdErr)
		} else { // otherwise does not exit and add a warning to the log
			p.cmdLog += fmt.Sprintf("WARNING: the process is set to continue after the error...\n")
		}
	}
	return nil
}

// OutboundRoute process an outbound route
func (p *Process) OutboundRoute(outRoute types.OutRoute) error {
	p.Info("processing outbound route %s: started", outRoute.Name)
	if outRoute.S3Store != nil {
		if err := p.ExportFiles(outRoute.S3Store); err != nil {
			return err
		}
	}
	if outRoute.PackageRegistry != nil {
		if err := p.PushPackages(outRoute.PackageRegistry); err != nil {
			return err
		}
	}
	if outRoute.ImageRegistry != nil {
		if err := p.PushImages(outRoute.ImageRegistry); err != nil {
			return err
		}
	}
	p.Info("processing outbound route %s: completed", outRoute.Name)
	return nil
}

// PushImages to a target container registry
func (p *Process) PushImages(imageRegistry *types.ImageRegistry) error {
	// tagging images & pushing
	p.Info("tagging and pushing images to docker registry: started")
	userPwd := fmt.Sprintf("%s:%s", imageRegistry.User, imageRegistry.Pwd)
	err := release.PushSpec(
		release.PushOptions{
			SpecPath: p.tmp,
			Host:     imageRegistry.Domain,
			Group:    imageRegistry.Group,
			User:     userPwd,
			Image:    true,
			Clean:    true,
			Logout:   true,
			ArtHome:  p.artHome(),
		})
	if err != nil {
		return p.Error("cannot push spec artefacts to image registry: %s", err)
	}
	p.Info("tagging and pushing artefacts to image registry: completed")
	return nil
}

// PushPackages packages to an Artisan registry
func (p *Process) PushPackages(pkgRegistry *types.PackageRegistry) error {
	// tagging artefacts & pushing
	p.Info("tagging and pushing artefacts to package registry: started")
	userPwd := fmt.Sprintf("%s:%s", pkgRegistry.User, pkgRegistry.Pwd)
	err := release.PushSpec(
		release.PushOptions{
			SpecPath: p.tmp,
			Host:     pkgRegistry.Domain,
			Group:    pkgRegistry.Group,
			User:     userPwd,
			Clean:    true,
			Logout:   true,
			ArtHome:  p.artHome(),
		})
	if err != nil {
		return p.Error("cannot push spec artefacts to package registry: %s", err)
	}
	p.Info("tagging and pushing artefacts to package registry: completed")
	return nil
}

// ExportFiles a spec to S3
func (p *Process) ExportFiles(s3Store *types.S3Store) error {
	// export packages
	p.Info("exporting packages: started")
	spec, err := release.NewSpec(p.tmp, "")
	if err != nil {
		return p.Error("cannot load spec.yaml from working folder: %s", err)
	}
	targetURI := fmt.Sprintf("%s/%s", s3Store.BucketURI, p.folderName)
	// ensure bucket exists and a notification is set up if bucket is being created
	_, err = resx.EnsureBucketNotification(targetURI, s3Store.Creds(), "spec.yaml", getARN(s3Store))
	if err != nil {
		return p.Error("cannot ensure bucket existence for %s: %s", targetURI, err)
	}
	err = release.ExportSpec(
		release.ExportOptions{
			Specification: spec,
			TargetUri:     targetURI,
			TargetCreds:   s3Store.Creds(),
			ArtHome:       p.artHome(),
		})
	if err != nil {
		return p.Error("cannot export spec to %s: %s", targetURI, err)
	}
	p.Info("exporting packages: completed")
	return nil
}

// ImportFiles from a specification
func (p Process) ImportFiles() error {
	// import artefacts
	p.Info("importing specification files: started")
	// remove spec specific artefacts from local registry
	// NOTE: do not use prune() to avoid removing tmp folder!
	err := p.cleanSpec()
	if err != nil {
		return p.Error("cannot cleanup local registry: %s", err)
	}
	// import spec in tmp folder
	_, err = release.ImportSpec(
		release.ImportOptions{
			TargetUri: p.tmp,
			ArtHome:   p.artHome(),
		})
	if err != nil {
		return p.Error("cannot import spec: %s", err)
	}
	// reload the registry changes
	p.reg.Load()
	p.Info("importing specification files: complete")
	return nil
}

// SendNotification send a notification
func (p *Process) SendNotification(nType db.NotificationType) error {
	// pipe must have a value
	if p.pipe == nil {
		return fmt.Errorf("cannot send notification, pipeline is not set")
	}
	var n *types.PipeNotification
	switch nType {
	case db.SuccessNotification:
		n = p.pipe.SuccessNotification
	case db.CmdFailedNotification:
		n = p.pipe.CmdFailedNotification
	case db.ErrorNotification:
		n = p.pipe.ErrorNotification
	default:
		return fmt.Errorf("notification type %s is not supported", nType)
	}
	// merges release-name
	subject := strings.ReplaceAll(n.Subject, "<<release-name>>", fmt.Sprintf("%s:%s", p.bucketName, p.folderName))
	// merges release-artefacts
	buf := bytes.Buffer{}
	count := 0
	if p.spec != nil {
		for _, pac := range p.spec.Packages {
			if count == 0 {
				buf.WriteString(fmt.Sprintf("packages:\n"))
			}
			buf.WriteString(fmt.Sprintf("%s\n", pac))
			count++
		}
		count = 0
		for _, img := range p.spec.Images {
			if count == 0 {
				buf.WriteString(fmt.Sprintf("images:\n"))
			}
			buf.WriteString(fmt.Sprintf("%s\n", img))
			count++
		}
	} else {
		buf.WriteString(fmt.Sprintf("Spec file not available\n"))
	}
	content := n.Content
	content = strings.ReplaceAll(content, "<<release-artefacts>>", buf.String())
	content = strings.ReplaceAll(content, "<<issue-log>>", p.issueLog())
	content = strings.ReplaceAll(content, "<<command-log>>", p.cmdLog)
	return postNotification(NotificationMsg{
		Recipient: n.Recipient,
		Type:      n.Type,
		Subject:   subject,
		Content:   content,
	})
}

// BeforeComplete run any additional tasks before completing the processing of the pipeline
func (p *Process) BeforeComplete(pipe *types.Pipeline) error {
	// if doorman is configured to connect to the cmdb
	if len(os.Getenv(OxWapiUri)) > 0 {
		// if catalogue publication is enabled
		if p.pipe.CMDB != nil && p.pipe.CMDB.Catalogue {
			if err := p.submitSpec(pipe.CMDB); err != nil {
				return fmt.Errorf("cannot submit spec '%s' version '%s' to the cmdb: %s", p.spec.Name, p.spec.Version, err)
			}
		}
		// if the spec contains functions
		if p.spec != nil && p.spec.Run != nil {
			// if the pipeline has cmdb configuration
			if p.pipe.CMDB != nil {
				if err := p.submitSpecFx(pipe.CMDB); err != nil {
					return fmt.Errorf("cannot submit function for spec '%s' version '%s' to the cmdb: %s", p.spec.Name, p.spec.Version, err)
				}
			} else {
				p.Warn("spec contains functions but the pipeline is not configured to send them to the CMDB")
			}
		}
	}
	return nil
}

type NotificationMsg struct {
	// Recipient of the notification if type is email
	Recipient string `yaml:"recipient" json:"recipient" example:"info@email.com"`
	// Type of the notification (e.g. email, snow, etc.)
	Type string `yaml:"type" json:"type" example:"email"`
	// Subject of the notification
	Subject string `yaml:"subject" json:"subject" example:"New Notification"`
	// Content of the template
	Content string `yaml:"content" json:"content" example:"A new event has been received."`
}

func (m NotificationMsg) Valid() error {
	if len(m.Subject) == 0 {
		return fmt.Errorf("subject is required")
	}
	if len(m.Content) == 0 {
		return fmt.Errorf("content is required")
	}
	if len(m.Recipient) == 0 {
		return fmt.Errorf("recipient is required")
	}
	return nil
}

// utilities

func (p *Process) issueLog() string {
	var issue []string
	log := strings.Split(p.log.String(), "\n")
	for _, line := range log {
		if strings.Contains(line, "ERROR") {
			issue = append(issue, line)
		}
	}
	return strings.Join(issue, "\n")
}

func (p *Process) logs() []string {
	l := p.log.String()
	lines := strings.Split(l, "\n")
	return lines[:len(lines)-1]
}

func (p *Process) verifyKeyFile() string {
	return filepath.Join(p.tmp, "verify_key.pgp")
}

func (p Process) bucketPath() string {
	return fmt.Sprintf("%s/%s", p.bucketName, p.folderName)
}

func (p *Process) cleanSpec() error {
	var names []string
	for _, name := range p.spec.Packages {
		names = append(names, name)
	}
	for _, name := range p.spec.Images {
		names = append(names, name)
	}
	return p.reg.Remove(names)
}

func (p *Process) cleanup() {
	p.Info("cleaning up path %s", p.tmp)
	os.RemoveAll(p.tmp)
}

func postNotification(n NotificationMsg) error {
	content, err := json.Marshal(n)
	if err != nil {
		return err
	}
	uri, err := GetNotificationURI()
	if err != nil {
		return err
	}
	requestURI := fmt.Sprintf("%s/notify", uri)
	req, err := http.NewRequest("POST", requestURI, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("cannot create http request: %s", err)
	}
	user, err := GetNotificationUser()
	if err != nil {
		return fmt.Errorf("missing configuration")
	}
	pwd, err := GetNotificationPwd()
	if err != nil {
		return fmt.Errorf("missing configuration")
	}
	req.Header.Add("Authorization", util.BasicToken(user, pwd))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// do we have a nil response?
	if resp == nil {
		return fmt.Errorf("response was empty for resource: %s", requestURI)
	}
	// check error status codes
	if resp.StatusCode > 201 {
		return fmt.Errorf("response returned status: %s; resource: %s", resp.Status, requestURI)
	}
	return nil
}

// notify of an error
func (p *Process) notify(err error) {
	// if there is an error
	if err != nil {
		// if there is no command related log recorded
		if len(p.cmdLog) == 0 {
			// send an error notification
			if err = p.SendNotification(db.ErrorNotification); err != nil {
				fmt.Printf("cannot send error notification: %s\n", err)
			}
		} else {
			// otherwise send a command failed notification
			if err = p.SendNotification(db.CmdFailedNotification); err != nil {
				fmt.Printf("cannot send command failed notification: %s\n", err)
			}
			p.cmdLog = ""
		}
	} else { // if there is not an error
		// sends a success notification
		if err = p.SendNotification(db.SuccessNotification); err != nil {
			fmt.Printf("cannot send success notification: %s\n", err)
		}
	}
}

func (p *Process) artHome() string {
	return p.reg.ArtHome
}

// submitSpec to cmdb
func (p *Process) submitSpec(cmdb *types.CMDB) error {
	p.Info("submit spec to cmdb: started")
	defer p.Info("submit spec to cmdb: complete")
	// prepares spec attributes
	attrs := make(map[string]interface{})
	attrs["VERSION"] = p.spec.Version
	if len(p.spec.Author) > 0 {
		attrs["AUTHOR"] = p.spec.Author
	}
	if len(p.spec.License) > 0 {
		attrs["LICENSE"] = p.spec.License
	}
	// prepares the spec
	var specMap map[string]interface{}
	jsonBytes, err := json.Marshal(p.spec)
	if err != nil {
		p.Warn("cannot marshal spec to json: %s", err)
	} else {
		err = json.Unmarshal(jsonBytes, &specMap)
		if err != nil {
			p.Warn("cannot unmarshal spec to map: %s", err)
		}
	}
	// send spec to cmdb
	result, oxErr := p.ox.PutItem(&oxc.Item{
		Key:         catalogueName(p.spec),
		Name:        p.spec.Name,
		Description: p.spec.Description,
		Type:        UCatalogueType,
		Tag:         toTags(cmdb.Tag),
		Meta:        specMap,
		Attribute:   attrs,
	})
	if oxErr != nil {
		return p.Error("cannot put spec to cmdb: %s, %s", result.Message, oxErr)
	}
	p.Info("spec submitted to cmdb, operation was %s, changed: %t", result.Operation, result.Changed)
	return nil
}

func (p *Process) submitSpecFx(cmdb *types.CMDB) error {
	if cmdb.Events != nil && len(cmdb.Events) > 0 {
		p.Info("submit spec functions to cmdb: started")
		defer p.Info("submit spec functions to cmdb: complete")
		artRegUser, err := GetArRegUser()
		if err != nil {
			return p.Error("missing %s variable: cannot submit spec function to cmdb", ArtRegUser)
		}
		artRegPwd, err := GetArRegPwd()
		if err != nil {
			return p.Error("missing %s variable: cannot submit spec function to cmdb", ArtRegPwd)
		}
		// if there is a mandate to record an event
		for _, event := range cmdb.Events {
			// check if  the event is in the spec
			found := false
			for _, run := range p.spec.Run {
				// if the event is in the spec
				if strings.EqualFold(run.Event, event) {
					found = true
					meta, metaErr := toMetaInput(run.Input)
					if metaErr != nil {
						return p.Error("cannot construct spec function input: cannot submit spec function to cmdb", metaErr)
					}
					result, oxErr := p.ox.PutItem(&oxc.Item{
						Key:         runFxName(run),
						Name:        run.Function,
						Description: fmt.Sprintf("%s - %s - %s", run.Package, run.Function, run.Event),
						Type:        ArtSpecFxType,
						Tag:         toTags(cmdb.Tag),
						Attribute: map[string]interface{}{
							"PACKAGE": run.Package,
							"FX":      run.Function,
							"USER":    artRegUser,
							"PWD":     artRegPwd,
						},
						Meta: meta,
					})
					if oxErr != nil {
						return p.Error("cannot put spec function %s (%s) to cmdb: %s, %s", run.Function, run.Event, result.Message, oxErr)
					}
					p.Info("spec function %s (%s) submitted to cmdb, operation was %s, changed: %t", run.Function, run.Event, result.Operation, result.Changed)
				}
			}
			if !found {
				p.Warn("event '%s' was defined in the pipeline configuration but was not found in the specification, cannot put to cmdb", event)
			}
		}
	}
	return nil
}

func toMetaInput(input *data.Input) (map[string]interface{}, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal spec function input to json: %s", err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal spec function input from json: %s", err)
	}
	return m, nil
}

func catalogueName(spec *release.Spec) string {
	return fmt.Sprintf("CA:%s:%s",
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToUpper(spec.Name), " ", ""), "/", "_"),
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToUpper(spec.Version), " ", ""), "/", "_"),
	)
}

func runFxName(run release.Run) string {
	return fmt.Sprintf("SF:%s:%s",
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToUpper(run.Package), " ", ""), "/", "_"),
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToUpper(run.Function), " ", ""), "/", "_"),
	)
}

func toTags(m []string) []interface{} {
	var tag []interface{}
	tag = make([]interface{}, len(m))
	for i, v := range m {
		tag[i] = v
	}
	return tag
}

func getARN(s3Store *types.S3Store) *notification.Arn {
	// work out default values for ARN
	partition := s3Store.Partition
	if len(partition) == 0 {
		partition = "minio"
	}
	service := s3Store.Service
	if len(service) == 0 {
		service = "sqs"
	}
	region := s3Store.Region
	accountId := s3Store.AccountID
	if len(accountId) == 0 {
		accountId = "_"
	}
	resource := s3Store.Resource
	if len(resource) == 0 {
		resource = "webhook"
	}
	return &notification.Arn{
		Partition: partition,
		Service:   service,
		Region:    region,
		AccountID: accountId,
		Resource:  resource,
	}
}
