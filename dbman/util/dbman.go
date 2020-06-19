//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"strings"
	"time"
)

var DM *DbMan

type DbMan struct {
	// configuration
	Cfg *Config
	// scrips manager
	script *ScriptManager
	// db provider
	db DatabaseProvider
	// is it ready?
	ready bool
}

func NewDbMan() (*DbMan, error) {
	// create an instance of the current configuration set
	cfg := NewAppCfg("", "")
	// create an instance of the script http client
	scriptClient, err := oxc.NewClient(NewOxClientConf(cfg))
	if err != nil {
		return nil, err
	}
	// create an instance of the script manager
	scriptManager, err := NewScriptManager(cfg, scriptClient)
	if err != nil {
		return nil, err
	}
	db := NewDbProvider(cfg)

	return &DbMan{
		Cfg:    cfg,
		script: scriptManager,
		db:     db,
	}, nil
}

func (dm *DbMan) GetReleasePlan() (*Plan, error) {
	return dm.script.fetchPlan()
}

func (dm *DbMan) GetReleaseInfo(appVersion string) (*Manifest, error) {
	return dm.script.fetchManifest(appVersion)
}

func (dm *DbMan) SaveConfig() {
	dm.Cfg.save()
}

func (dm *DbMan) SetConfig(key string, value string) {
	dm.Cfg.set(key, value)
}

func (dm *DbMan) GetConfig(key string) {
	dm.Cfg.Get(key)
}

// toString the current configuration set to stdout
func (dm *DbMan) ConfigSetAsString() string {
	return dm.Cfg.toString()
}

// use the configuration set specified by name
// name: the name of the configuration set to use
// filepath: the path to the configuration set
func (dm *DbMan) UseConfigSet(filepath string, name string) {
	dm.Cfg.load(filepath, name)
}

// get the content of the current configuration set
func (dm *DbMan) GetConfigSet() string {
	return dm.Cfg.ConfigFileUsed()
}

// get the current configuration directory
func (dm *DbMan) GetConfigSetDir() string {
	return dm.Cfg.root.path()
}

// performs various connectivity checks using the information in the current configuration set
// returns a map containing entries with the type of check and the result
func (dm *DbMan) CheckConfigSet() map[string]string {
	results := make(map[string]string)
	_, err := dm.script.fetchPlan()
	if err != nil {
		fmt.Printf("!!! check failed: %v\n", err)
		results["scripts uri"] = err.Error()
	} else {
		results["scripts uri"] = "OK"
	}
	// try and connect to the database
	// create a dummy action with no scripts to test the connection
	testConnCmd := &Command{
		Name:          "test connection",
		Description:   "",
		Transactional: false,
		AsAdmin:       true,
		UseDb:         false,
		Scripts:       []Script{},
	}
	_, err = dm.db.RunCommand(testConnCmd)
	if err != nil {
		results["db connection"] = fmt.Sprintf("FAILED: %v", err)
	} else {
		results["db connection"] = "OK"
	}
	return results
}

func (dm *DbMan) Create() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	appVer := dm.Cfg.Get(AppVersion)
	// get database release version
	log.WriteString(fmt.Sprintf("? I am checking that the database '%s' does not already exist\n", dm.Cfg.Get(DbName)))
	appVersion, dbVersion, err := dm.db.GetVersion()
	if err == nil {
		// there is already a database and cannot continue
		return log, errors.New(fmt.Sprintf("!!! I have found an existing database version %v, which is for application version %v", dbVersion, appVersion)), time.Since(start)
	}
	// fetch the release manifest for appVersion
	log.WriteString(fmt.Sprintf("? I am retrieving the release manifest for application version '%v'\n", dm.Cfg.Get(AppVersion)))
	manifest, err := dm.script.fetchManifest(appVer)
	if err != nil {
		return log, err, time.Since(start)
	}
	// get the commands for the create action
	cmds := manifest.getCommands(manifest.Create.Commands)
	// run the commands on the database
	output, err := dm.runCommands(cmds, manifest)
	log.WriteString(output.String())
	// return
	return log, err, time.Since(start)
}

func (dm *DbMan) Deploy() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	appVer := dm.Cfg.Get(AppVersion)
	// get database release version
	appVersion, dbVersion, err := dm.db.GetVersion()
	if err == nil && len(appVersion) > 0 {
		// there is already a database with a pre-existing deployment so cannot continue
		return log, errors.New(fmt.Sprintf("!!! I have found an existing database version %v, which is for application version %v", dbVersion, appVersion)), time.Since(start)
	}
	// fetch the release manifest for appVersion
	manifest, err := dm.script.fetchManifest(appVer)
	if err != nil {
		return log, err, time.Since(start)
	}
	// get the commands for the deploy action
	cmds := manifest.getCommands(manifest.Deploy.Commands)
	// run the commands on the database
	output, err := dm.runCommands(cmds, manifest)
	log.WriteString(output.String())
	if err != nil {
		return log, err, time.Since(start)
	}
	// update release version
	err = dm.db.SetVersion(appVer, manifest.DbVersion, fmt.Sprintf("Database Release %v", manifest.DbVersion), dm.Cfg.Get(SchemaURI))
	// return
	return log, err, time.Since(start)
}

func (dm *DbMan) Upgrade() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	return log, nil, time.Since(start)
}

func (dm *DbMan) RunQuery(manifest *Manifest, query *Query, params []string) (Table, time.Duration, error) {
	start := time.Now()
	// fetch the query content
	query, err := dm.script.fetchQueryContent(dm.Cfg.Get(AppVersion), manifest.QueriesPath, *query, params)
	if err != nil {
		return Table{}, time.Since(start), errors.New(fmt.Sprintf("!!! I cannot fetch content for query: %v\n", query.Name))
	}
	// run the query
	result, err := dm.db.RunQuery(query, params)
	return result, time.Since(start), err
}

func (dm *DbMan) CheckReady() (bool, error) {
	// ready if check passes
	results := dm.CheckConfigSet()
	for check, result := range results {
		if !strings.Contains(strings.ToLower(result), "ok") {
			dm.ready = false
			return false, errors.New(fmt.Sprintf("%v: %v", check, result))
		}
	}
	dm.ready = true
	return true, nil
}

// launch DbMan as an http server
func (dm *DbMan) Serve() {
	server := NewServer(dm.Cfg)
	server.Serve()
}

func (dm *DbMan) runCommands(cmds []Command, manifest *Manifest) (log bytes.Buffer, err error) {
	log = bytes.Buffer{}
	// fetch the scripts for the commands
	var commands []*Command
	for _, cmd := range cmds {
		cmd, err := dm.script.fetchCommandContent(dm.Cfg.Get(AppVersion), manifest.CommandsPath, cmd)
		if err != nil {
			return log, err
		}
		commands = append(commands, cmd)
	}
	// execute the commands
	for _, c := range commands {
		log.WriteString(fmt.Sprintf("? I have started execution of the command '%s'\n", c.Name))
		output, err := dm.db.RunCommand(c)
		log.WriteString(output)
		if err != nil {
			log.WriteString(fmt.Sprintf("!!! the execution of the command '%s' has failed: %s\n", c.Name, err))
			return log, err
		}
		log.WriteString(fmt.Sprintf("? the execution of the command '%s' has succeeded\n", c.Name))
	}
	return log, err
}
