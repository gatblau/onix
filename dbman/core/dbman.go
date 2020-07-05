//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
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
	db *DatabaseProviderManager
	// is it ready?
	ready bool
}

func NewDbMan() (*DbMan, error) {
	// create an instance of the current configuration set
	cfg := NewConfig("", "")
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
	// load the database provider
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}
	if db.Provider() == nil {
		return nil, errors.New("!!! database Provider plugin loading failed")
	}
	// pass in DbMan's configuration to the database provider
	result := NewParameterFromJSON(db.Provider().Setup(cfg.All()))
	// if an error message came back
	if result.HasError() {
		// return the error
		return nil, result.Error()
	}
	// output the setup log
	result.PrintLog()
	// otherwise returns a DbMan instance
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
	dm.Cfg.Save()
}

func (dm *DbMan) SetConfig(key string, value string) {
	dm.Cfg.Set(context.Background(), key, value)
}

// toString the current configuration set to stdout
func (dm *DbMan) ConfigSetAsString() string {
	return dm.Cfg.ToString()
}

// use the configuration set specified by name
// name: the name of the configuration set to use
// filepath: the path to the configuration set
func (dm *DbMan) UseConfigSet(filepath string, name string) {
	dm.Cfg.Load(filepath, name)
}

// get the content of the current configuration set
func (dm *DbMan) GetConfigSet() string {
	return dm.Cfg.ConfigFileUsed()
}

// get the current configuration directory
func (dm *DbMan) GetConfigSetDir() string {
	return dm.Cfg.Cache.Path()
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
	r := dm.DbPlugin().RunCommand(testConnCmd.ToString())
	result := NewParameterFromJSON(r)
	if result.HasError() {
		results["db connection"] = fmt.Sprintf("FAILED: %v", result.Error())
	} else {
		results["db connection"] = "OK"
	}
	return results
}

func (dm *DbMan) Create() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	appVer := dm.get(AppVersion)
	// get database release version
	log.WriteString(fmt.Sprintf("? I am checking that the database '%s' does not already exist\n", dm.get(DbName)))
	r := dm.DbPlugin().GetVersion()
	result := NewParameterFromJSON(r)
	// if no error then
	if !result.HasError() {
		v := result.GetVersion()
		// if there is a version
		if v != nil {
			// there is already a database and cannot continue
			return log, errors.New(fmt.Sprintf("!!! I have found an existing database version %v, which is for application version %v", v.DbVersion, v.AppVersion)), time.Since(start)
		}
	}
	// fetch the release manifest for appVersion
	log.WriteString(fmt.Sprintf("? I am retrieving the release manifest for application version '%v'\n", dm.get(AppVersion)))
	manifest, err := dm.script.fetchManifest(appVer)
	if err != nil {
		return log, err, time.Since(start)
	}
	// get the commands for the create action
	cmds := manifest.GetCommands(manifest.Create.Commands)
	// run the commands on the database
	output, err := dm.runCommands(cmds, manifest)
	log.WriteString(output.String())
	// return
	return log, err, time.Since(start)
}

func (dm *DbMan) Deploy() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	appVer := dm.get(AppVersion)
	// get database release version
	r := dm.DbPlugin().GetVersion()
	result := NewParameterFromJSON(r)
	v := result.GetVersion()
	if !result.HasError() && len(v.AppVersion) > 0 {
		// there is already a database with a pre-existing deployment so cannot continue
		return log, errors.New(fmt.Sprintf("!!! I have found an existing database version %v, which is for application version %v",
			result.GetString("dbVersion"),
			result.GetString("appVersion"))), time.Since(start)
	}
	// fetch the release manifest for appVersion
	manifest, err := dm.script.fetchManifest(appVer)
	if err != nil {
		return log, err, time.Since(start)
	}
	// get the commands for the deploy action
	cmds := manifest.GetCommands(manifest.Deploy.Commands)
	// run the commands on the database
	output, err := dm.runCommands(cmds, manifest)
	log.WriteString(output.String())
	if err != nil {
		return log, err, time.Since(start)
	}
	// update release version
	input := &Version{
		AppVersion:  appVer,
		DbVersion:   manifest.DbVersion,
		Description: fmt.Sprintf("Database Release %v", manifest.DbVersion),
		Source:      dm.get(SchemaURI),
	}
	setVerResult := NewParameterFromJSON(dm.DbPlugin().SetVersion(input.ToString()))
	if setVerResult.HasError() {
		err = setVerResult.Error()
	}
	// return
	return log, err, time.Since(start)
}

func (dm *DbMan) Upgrade() (log bytes.Buffer, err error, elapsed time.Duration) {
	start := time.Now()
	log = bytes.Buffer{}
	return log, nil, time.Since(start)
}

func (dm *DbMan) RunQuery(manifest *Manifest, query *Query, params []string) (*Table, time.Duration, error) {
	start := time.Now()
	// fetch the query content
	query, err := dm.script.fetchQueryContent(dm.get(AppVersion), manifest.QueriesPath, *query, params)
	if err != nil {
		return nil, time.Since(start), errors.New(fmt.Sprintf("!!! I cannot fetch content for query: %v\n", query.Name))
	}
	// run the query on the plugin
	r := dm.DbPlugin().RunQuery(query.ToString())
	// recreate plugin response into parameter
	result := NewParameterFromJSON(r)
	// return table and error
	return result.GetTable(), time.Since(start), result.Error()
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
		cmd, err := dm.script.fetchCommandContent(dm.get(AppVersion), manifest.CommandsPath, cmd)
		if err != nil {
			return log, err
		}
		commands = append(commands, cmd)
	}
	// execute the commands
	for _, c := range commands {
		log.WriteString(fmt.Sprintf("? I have started execution of the command '%s'\n", c.Name))
		r := dm.DbPlugin().RunCommand(c.ToString())
		result := NewParameterFromJSON(r)
		if result.HasError() {
			log.WriteString(fmt.Sprintf("!!! the execution of the command '%s' has failed: %s\n", c.Name, result.Error()))
			return log, result.Error()
		}
		log.WriteString(fmt.Sprintf("? the execution of the command '%s' has succeeded\n", c.Name))
	}
	return log, err
}

func (dm *DbMan) get(key string) string {
	return dm.Cfg.GetString(key)
}

func (dm *DbMan) DbPlugin() DatabaseProvider {
	return dm.db.Provider()
}

func (dm *DbMan) GetDbInfo() (*DbInfo, error) {
	// query the plugin for serialised information
	infoString := dm.DbPlugin().GetInfo()
	// unmarshal the parameter
	info := NewParameterFromJSON(infoString)
	// is there an error?
	if info.HasError() {
		// return the error
		return nil, info.Error()
	}
	// extract the result from the parameter
	result := info.Get("result")
	// is the result nil?
	if result == nil {
		// something went wrong, the plugin did not set the result value
		return nil, errors.New("!!! The database plugin did not return a result\n")
	}
	// there is a result but is it a map?
	if r, ok := result.(map[string]interface{}); ok {
		return NewDbInfoFromMap(r)
	}
	// if not
	return nil, errors.New("!!! The database plugin did not return a result of the correct type (i.e. map[string]interface{})\n")
}
