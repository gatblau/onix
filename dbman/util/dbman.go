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
)

var DM *DbMan

type DbMan struct {
	// configuration
	Cfg *AppCfg
	// scrips manager
	script *ScriptManager
	// database manager
	db DatabaseProvider
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
	if err != nil {
		return err
	}
	return nil
}

func (dm *DbMan) Deploy(targetAppVersion string) error {
	var (
		newDb bool = false
	)
	// check if the database exists
	exist, _ := dm.db.DbExists()
	// if the database does not exists, then create it
	if !exist {
		fmt.Printf("! I could not find the database '%v': proceeding to create it.\n", dm.Cfg.Get(DbName))
		err := dm.InitialiseDb()
		if err != nil {
			return err
		}
		// a new database has been created
		newDb = true
	}
	// if the database already exists, the it needs to check what is the current version
	if !newDb {
		fmt.Printf("? I am checking database version compatibility for requested version '%v'\n", targetAppVersion)
		// get database version
		currentAppVersion, _, _ := dm.db.GetVersion()
		// if the currently deployed db version exists and it does not match the target version
		if len(currentAppVersion) > 0 && currentAppVersion != targetAppVersion {
			// should not deploy the schemas, more likely an upgrade is needed?
			return errors.New(fmt.Sprintf("!!! I cannot deploy the database schemas for application version '%v' as it differs from the existing application version '%v'\n", currentAppVersion, targetAppVersion))
		}
	}
	fmt.Printf("? I am fetching database release info for application version '%v'.\n", targetAppVersion)
	release, err := dm.script.fetchRelease(targetAppVersion)
	if err != nil {
		return err
	}
	// deploys the release
	return dm.db.DeployDb(release)
}
