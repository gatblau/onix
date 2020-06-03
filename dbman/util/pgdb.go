//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

// the database provider for PostgreSQL
// NOTE: database providers implicitly implement the DatabaseProvider interface
type PgSQLProvider struct {
	cfg *AppCfg
}

// creates a new db instance
func NewDb(appCfg *AppCfg) DatabaseProvider {
	switch strings.ToLower(appCfg.Get(DbProvider)) {
	case "pgsql":
		// load the default native postgres provider
		return &PgSQLProvider{
			cfg: appCfg,
		}
	default:
		// only supports connections to postgres at the moment
		// in time, a plugin approach for database providers could be implemented
		panic(errors.New(fmt.Sprintf("!!! the database provider '%v' is not supported.", appCfg.Get(DbProvider))))
	}
}

// check a connection to the database can be established
func (db *PgSQLProvider) CanConnect() (bool, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString(false))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}

// checks if the database exists
func (db *PgSQLProvider) Exists() (bool, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString(true))
	if err != nil {
		fmt.Printf("!!! I cannot connect to database: %v\n", err)
		return false, err
	}
	defer conn.Close()
	var count int
	err = conn.QueryRow(context.Background(), "SELECT 1 from pg_database WHERE datname='$1';", db.get(DbName)).Scan(&count)
	if err != nil {
		fmt.Printf("!!! I cannot check if the database exists: %v\n", err)
	}
	return count == 1, err
}

// initialises the database from db init info
func (db *PgSQLProvider) Initialise(init *DbInit) error {
	// connect to the database
	conn, err := pgxpool.Connect(context.Background(), db.connString(true))
	if err != nil {
		fmt.Printf("!!! I am unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close()

	for _, item := range init.Items {
		// prepares to execute script
		// print the action to be carried out
		fmt.Println(item.Action)
		// merge any script variables
		for _, value := range item.Vars {
			// get the value of the variable from the configuration
			confValue := db.cfg.Get(value.From)
			// replace all occurrences of the placeholder (value.Name) with the confValue
			result := strings.Replace(item.Script, value.Name, confValue, -1)
			// update the value of script with the merged result
			item.Script = result
		}
		// execute the script
		_, err = conn.Exec(context.Background(), item.Script)
		// if an error is encountered then exit
		if err != nil {
			return err
		}
	}
	return nil
}

// gets the current app and db version
func (db *PgSQLProvider) GetVersion() (appVersion string, dbVersion string, err error) {
	conn, err := db.newConn(false)
	if err != nil {
		return "", "", err
	}
	defer conn.Close()
	rows, _ := conn.Query(context.Background(), "SELECT application_version, database_version from version ORDER BY time DESC LIMIT 1;")
	if rows.Next() {
		var appVer string
		var dbVer string
		err := rows.Scan(&appVer, &dbVer)
		if err != nil {
			return "", "", err
		}
		return appVer, dbVer, err
	}
	return "", "", rows.Err()
}

// deploy the database schemas
func (db *PgSQLProvider) Deploy() error {
	return nil
}

func (db *PgSQLProvider) Upgrade() error {
	return nil
}

// return a configuration item
func (db *PgSQLProvider) get(key string) string {
	return db.cfg.Get(key)
}

// return the connection string
// admin:
//  - if true, a connection using the postgres user is returned
//  - if false, a connection using the database user is returned
func (db *PgSQLProvider) connString(admin bool) string {
	if admin {
		return fmt.Sprintf("postgresql://%v:%v@%v:%v",
			"postgres",
			db.get(DbAdminPwd),
			db.get(DbHost),
			db.get(DbPort))
	}
	return fmt.Sprintf("postgresql://%v:%v@%v:%v",
		db.get(DbUsername),
		db.get(DbPassword),
		db.get(DbHost),
		db.get(DbPort))
}

// create the version tracking table in the target database
func (db *PgSQLProvider) CreateVersionTable() error {
	conn, err := db.newConn(false)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(context.Background(), PgSQLVersionTable)
	return err
}

// create a new connection
func (db *PgSQLProvider) newConn(admin bool) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString(admin))
	if err != nil {
		fmt.Printf("!!! I cannot connect to the database: %v\n", err)
		return nil, err
	}
	return conn, nil
}
