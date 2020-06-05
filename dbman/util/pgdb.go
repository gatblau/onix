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
	"time"
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

// check a connection to the server can be established
func (db *PgSQLProvider) CanConnectToServer() (bool, error) {
	_, err := db.newConn(true, false)
	return err != nil, err
}

// checks if the database exists
func (db *PgSQLProvider) DbExists() (bool, error) {
	conn, err := db.newConn(false, true)
	if err != nil {
		fmt.Printf("!!! I cannot connect to database: %v\n", err)
		return false, err
	}
	defer conn.Close()
	var count int
	sql := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s';", db.get(DbName))
	err = conn.QueryRow(context.Background(), sql).Scan(&count)
	return err == nil, err
}

// initialises the database from db init info
func (db *PgSQLProvider) InitialiseDb(init *DbInit) error {
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
		// connect to the server or database as admin or user
		conn, err := db.newConn(item.Admin, item.Db)
		if err != nil {
			fmt.Printf("!!! I am unable to connect to database: %v\n", err)
			return err
		}
		// execute the script
		_, err = conn.Exec(context.Background(), item.Script)
		// if an error is encountered then exit
		if err != nil {
			return err
		}
		// closes the connection
		conn.Close()
	}
	return nil
}

// gets the current app and db version
func (db *PgSQLProvider) GetVersion() (appVersion string, dbVersion string, err error) {
	conn, err := db.newConn(false, true)
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
func (db *PgSQLProvider) DeployDb(release *Release) error {
	conn, err := db.newConn(false, true)
	if err != nil {
		return err
	}
	// deploy the schemas
	for _, schema := range release.Schemas {
		_, err := conn.Exec(context.Background(), schema)
		if err != nil {
			return err
		}
	}
	// deploy the functions
	for _, function := range release.Functions {
		_, err := conn.Exec(context.Background(), function)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *PgSQLProvider) UpgradeDb() error {
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
// database:
//  - if true, adds the database name to the connection string
func (db *PgSQLProvider) connString(admin bool, database bool) string {
	connStr := ""
	if admin {
		connStr = fmt.Sprintf("postgresql://%v:%v@%v:%v",
			"postgres",
			db.get(DbAdminPwd),
			db.get(DbHost),
			db.get(DbPort))
	} else {
		connStr = fmt.Sprintf("postgresql://%v:%v@%v:%v",
			db.get(DbUsername),
			db.get(DbPassword),
			db.get(DbHost),
			db.get(DbPort))
	}
	if database {
		connStr = fmt.Sprintf("%v/%v", connStr, db.get(DbName))
	}
	return connStr
}

// create the version tracking table in the target database
func (db *PgSQLProvider) CreateVersionTable() error {
	conn, err := db.newConn(false, true)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(context.Background(), pgSQLVersionTable)
	return err
}

// this type carries either a connection or an error
// used in a channel used by the connection go routine
type conn struct {
	conn *pgxpool.Pool
	err  error
}

// create a new database connection
// if it cannot connect within 5 seconds, it returns an error
func (db *PgSQLProvider) newConn(admin bool, database bool) (*pgxpool.Pool, error) {
	// this channel receives an connection
	connect := make(chan conn, 1)
	// this channel receives a timeout flag
	timeout := make(chan bool, 1)
	// launch a go routine to try the database connection
	go func() {
		// connects to the database
		c, e := pgxpool.Connect(context.Background(), db.connString(admin, database))
		// sends connection over channel
		connect <- conn{conn: c, err: e}
	}()
	// launch a go routine
	go func() {
		// timeout period is 3.5 secs
		time.Sleep(3.5e9)
		timeout <- true
	}()

	select {
	// the connection has been established before the timeout
	case connection := <-connect:
		{
			if connection.err != nil {
				return nil, connection.err
			}
			return connection.conn, nil
		}
	// the connection has not yet returned when the timeout happens
	case <-timeout:
		{
			return nil, errors.New("!!! I cannot connect to the database, the timed out period has elapsed\n")
		}
	}
}

const pgSQLVersionTable = `
DO
  $$
    BEGIN
      ---------------------------------------------------------------------------
      -- VERSION - version of releases (not only database)
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'version')
      THEN
        CREATE TABLE version
        (
          application_version CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
          database_version    CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
          description         TEXT COLLATE pg_catalog."default",
          time                timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
          scripts_source      character varying(250),
          CONSTRAINT version_app_version_db_release_pk PRIMARY KEY (application_version, database_version)
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE version
          OWNER to onix;
      END IF;
    END;
    $$
`
