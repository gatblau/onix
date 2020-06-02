//   Onix Config Db - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type PgDb struct {
	cfg *AppCfg
}

// creates a new PostgreSql db instance
func NewPgDb(appCfg *AppCfg) Db {
	return &PgDb{
		cfg: appCfg,
	}
}

// check a connection to the database can be established
func (db *PgDb) CanConnect() (bool, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString())
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}

// checks if the database exists
func (db *PgDb) Exists() (bool, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString())
	if err != nil {
		fmt.Printf("oops! I cannot connect to database: %v\n", err)
		return false, err
	}
	defer conn.Close()
	var count int
	err = conn.QueryRow(context.Background(), "SELECT 1 from pg_database WHERE datname='$1';", db.get(DbName)).Scan(&count)
	if err != nil {
		fmt.Printf("oops! I cannot check if the database exists: %v\n", err)
	}
	return count == 1, err
}

// create the Onix database
func (db *PgDb) Initialise() error {
	conn, err := pgxpool.Connect(context.Background(), db.connString())
	if err != nil {
		fmt.Printf("oops! I am unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close()
	fmt.Printf("creating database: %s", db.get(DbName))
	_, err = conn.Exec(context.Background(), "CREATE DATABASE $1;", db.get(DbName))
	if err != nil {
		return err
	}
	fmt.Printf("creating user: %s", db.get(DbUsername))
	_, err = conn.Exec(context.Background(), "CREATE USER $1 WITH PASSWORD '$2';", db.get(DbUsername), db.get(DbPassword))
	if err != nil {
		return err
	}
	fmt.Printf("installing database extensions")
	_, err = conn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS hstore;")
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS intarray;")
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), VersionTable)
	if err != nil {
		return err
	}
	return nil
}

// gets the current app and db version
func (db *PgDb) GetVersion() (string, string, error) {
	conn, err := pgxpool.Connect(context.Background(), db.connString())
	if err != nil {
		log.Error().Msgf("unable to connect to database: %v\n", err)
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
func (db *PgDb) Deploy() error {
	// check if the database exists
	exist, err := db.Exists()
	if err != nil {
		return err
	}
	// if the database does not exists, then create it
	if !exist {
		err = db.Initialise()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("I have found database %s, skipping creation", db.get(DbName))
	}
	_, dbVer, err := db.GetVersion()
	if err != nil {
		return err
	}
	// only install the schemas if there is none
	if len(dbVer) == 0 {

	}
	return nil
}

func (db *PgDb) Upgrade() error {
	return nil
}

// return a configuration item
func (db *PgDb) get(key string) string {
	return db.cfg.Get(key)
}

// return the connection string
func (db *PgDb) connString() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v",
		db.get(DbUsername),
		db.get(DbPassword),
		db.get(DbHost),
		db.get(DbPort))
}
