//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"strings"
	"time"
)

var DM *DbMan

type DbMan struct {
	// configuration
	Cfg *AppCfg
	// scrips manager
	script *ScriptManager
	// database manager
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
	dbProvider := NewDb(cfg)
	return &DbMan{
		Cfg:    cfg,
		script: scriptManager,
		db:     dbProvider,
	}, nil
}

func (dm *DbMan) GetReleasePlan() (*Plan, error) {
	return dm.script.fetchPlan()
}

func (dm *DbMan) GetReleaseInfo(appVersion string) (*Release, error) {
	return dm.script.fetchRelease(appVersion)
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

// print the current configuration set to stdout
func (dm *DbMan) PrintConfigSet() {
	dm.Cfg.print()
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
	// try and fetch the getReleaseInfo plan
	_, err := dm.script.fetchPlan()
	if err != nil {
		fmt.Printf("!!! check failed: %v\n", err)
		results["scripts uri"] = err.Error()
	} else {
		results["scripts uri"] = "OK"
	}
	// try and connect to the database
	_, err = dm.db.CanConnectToServer()
	if err != nil {
		results["db connection"] = fmt.Sprintf("FAILED: %v", err)
	} else {
		results["db connection"] = "OK"
	}
	return results
}

// initialises the database (i.e. create database, user, extensions, etc)
// it does not include schema deployment or upgrades
func (dm *DbMan) InitialiseDb() error {
	start := time.Now()

	// check if the database exists
	exist, err := dm.db.DbExists()
	// if there is an error, could not connect to the database
	if err != nil {
		return err
	}
	if !exist {
		fmt.Printf("? I am fetching database initialisation info.\n")
		init, err := dm.script.fetchInit()
		if err != nil {
			return err
		}
		fmt.Printf("? I am applying the database initialisation scripts.\n")
		err = dm.db.InitialiseDb(init)
		if err != nil {
			return err
		}
		fmt.Printf("? I am creating the database version tracking table.\n")
		// NOTE: its schema is enforced by DbMan
		err = dm.db.CreateVersionTable()
		// logs the time taken
		fmt.Printf("? I have initialised the database in %v", time.Since(start))
		return err
	}
	fmt.Printf("? I cannot execute the initialisation because the database already exist\n")
	return nil
}

func (dm *DbMan) Deploy(targetAppVersion string) error {
	start := time.Now()
	// check if the database exists
	exist, err := dm.db.DbExists()
	// if there is an error, could not connect to the database
	if err != nil {
		return err
	}
	// if the database does not exists, then exit (init should be run first)
	if !exist {
		return errors.New(fmt.Sprintf("! I could not find the database '%v': call 'dbman db init' before attemting to deploy.\n", dm.Cfg.Get(DbName)))
	}
	// get database version
	appVer, dbVer, err := dm.db.GetVersion()
	// if the version cannot be retrieved return
	if err != nil {
		return err
	}
	// if there is a previous version, deploy should not be called, returns
	if len(appVer) > 0 && len(dbVer) > 0 {
		return errors.New(fmt.Sprintf("!!! I have found a previous deployment for application version '%v', I cannot continue.\n", appVer))
	}
	// we have an empty version table so we are ready to deploy
	fmt.Printf("? I am fetching database release info for application version '%v'.\n", targetAppVersion)
	release, err := dm.script.fetchRelease(targetAppVersion)
	if err != nil {
		return err
	}
	// deploys the release
	err = dm.db.DeployDb(release)
	if err != nil {
		return err
	}
	// add the deployed version in the tracking table
	err = dm.db.InsertVersion(targetAppVersion, release.Release, "deployed by DbMan", dm.Cfg.Get(SchemaURI))
	// logs the time taken
	fmt.Printf("? I have deployed the database in %v", time.Since(start))
	return err
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
