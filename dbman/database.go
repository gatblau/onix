//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

type Database struct {
	cfg *Config
}

// checks if the database exists
func (db *Database) exists() (bool, error) {
	conn, err := pgx.Connect(context.Background(), db.cfg.DbConnString)
	if err != nil {
		log.Error().Msgf("unable to connect to database: %v\n", err)
		return false, err
	}
	defer conn.Close(context.Background())
	var count int
	err = conn.QueryRow(context.Background(), "SELECT 1 from pg_database WHERE datname='$1';", db.cfg.DbName).Scan(&count)
	if err != nil {
		log.Error().Msgf("cannot check if database exists: %v\n", err)
	}
	return count == 1, err
}

// createDb the Onix database
func (db *Database) createDb() error {
	conn, err := pgx.Connect(context.Background(), db.cfg.DbConnString)
	if err != nil {
		log.Error().Msgf("unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())
	log.Info().Msgf("creating database: %s", db.cfg.DbName)
	_, err = conn.Exec(context.Background(), "CREATE DATABASE $1;", db.cfg.DbName)
	if err != nil {
		return err
	}
	log.Info().Msgf("creating user: %s", db.cfg.DbUsername)
	_, err = conn.Exec(context.Background(), "CREATE USER $1 WITH PASSWORD '$2';", db.cfg.DbUsername, db.cfg.Password)
	if err != nil {
		return err
	}
	log.Info().Msg("installing database extensions")
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
func (db *Database) getVersion() (string, string, error) {
	conn, err := pgx.Connect(context.Background(), db.cfg.DbConnString)
	if err != nil {
		log.Error().Msgf("unable to connect to database: %v\n", err)
		return "", "", err
	}
	defer conn.Close(context.Background())
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

// installDb the database schemas
func (db *Database) installDb() error {
	// check if the database exists
	exist, err := db.exists()
	if err != nil {
		return err
	}
	// if the database does not exists, then createDb it
	if !exist {
		err = db.createDb()
		if err != nil {
			return err
		}
	} else {
		log.Info().Msgf("database %s found, skipping creation", db.cfg.DbName)
	}
	_, dbVer, err := db.getVersion()
	if err != nil {
		return err
	}
	// only install the schemas if there is none
	if len(dbVer) == 0 {

	}
	return nil
}
