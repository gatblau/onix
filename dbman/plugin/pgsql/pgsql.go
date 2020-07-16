//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
)

// DbMan's database provider for PostgreSQL
// PgSQLProvider implicitly implements the DatabaseProvider interface
// NOTE: the provider uses connection pooling so should not call conn.Close()
type PgSQLProvider struct {
	cfg *Conf
}

// pass DbMan configuration to the database provider
// config: map[string]interface{} serialised as a json string
func (db *PgSQLProvider) Setup(config *Conf) error {
	// allocate the parsed object to cfg
	db.cfg = config
	return nil
}

// retrieve database information
// return: map[string]interface{} serialised as a json string
func (db *PgSQLProvider) GetVersion() (*Version, error) {
	conn, err := db.newConn(true, true)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), `
		SELECT appVersion, dbVersion, description, time, source
		FROM "version"
		ORDER BY time DESC
		LIMIT 1`)
	if err != nil {
		return nil, err
	}
	var (
		appVersion, dbVersion, description, source string
		time                                       time.Time
	)
	if rows.Next() {
		rows.Scan(&appVersion, &dbVersion, &description, &time, &source)
		rows.Close()
	} else {
		// no results
		rows.Close()
		return nil, err
	}
	// we have version information so populate result
	v := &Version{
		AppVersion:  appVersion,
		DbVersion:   dbVersion,
		Description: description,
		Source:      source,
		Time:        time,
	}
	return v, err
}

func (db *PgSQLProvider) RunCommand(command *Command) (bytes.Buffer, error) {
	log := bytes.Buffer{}
	conn, err := db.newConn(command.AsAdmin, command.UseDb)
	if err != nil {
		return log, err
	}
	if command.Transactional {
		log.WriteString(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", command.Transactional),
			db.label("as an admin", "as a user", command.AsAdmin),
			db.label("to the db", "to the server", command.UseDb)))
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return log, err
		}
		for _, script := range command.Scripts {
			_, err := tx.Exec(context.Background(), script.Content)
			log.WriteString(fmt.Sprintf("? I have executed the script '%s'\n", script.Name))
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				tx.Rollback(context.Background())
				return log, err
			}
		}
		tx.Commit(context.Background())
	} else {
		log.WriteString(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", command.Transactional),
			db.label("as an admin", "as a user", command.AsAdmin),
			db.label("to the db", "to the server", command.UseDb)))
		for _, script := range command.Scripts {
			_, err := conn.Exec(context.Background(), script.Content)
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				return log, err
			}
		}
	}
	return log, err
}

func (db *PgSQLProvider) RunQuery(query *Query) (*Table, error) {
	conn, err := db.newConn(false, true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.Query(context.Background(), query.Content)
	if err != nil {
		return nil, err
	}
	header := make(Row, 0)
	rows := make([]Row, 0)
	for result.Next() {
		// only the first time round populate the header
		if len(header) == 0 {
			for _, desc := range result.FieldDescriptions() {
				header = append(header, string(desc.Name))
			}
		}
		// create a new row
		row := make(Row, 0)
		// populate the row with returned values from the query
		values, err := result.Values()
		if err != nil {
		}
		for _, value := range values {
			if v, ok := value.(string); ok {
				row = append(row, v)
			}
			if v, ok := value.(time.Time); ok {
				row = append(row, v.String())
			}
		}
		// add the row to the row set
		rows = append(rows, row)
	}
	result.Close()
	return &Table{
		Header: header,
		Rows:   rows,
	}, err
}

// set the release version
// versionInfo: a json serialised map[string]interface{} containing version information
func (db *PgSQLProvider) SetVersion(version *Version) error {
	// create a db connection
	conn, err := db.newConn(false, true)
	if err != nil {
		return err
	}
	// find out if version table exists
	_, err = conn.Query(context.Background(), "SELECT * FROM version")
	if err != nil {
		dbUser, found := db.get("Db.Username")
		if !found {
			return errors.New(fmt.Sprint("!!! could not find Db.Username config value"))
		}
		// could not find version table, so try and create it
		_, err := conn.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE "version"
            (
                appVersion  CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
                dbVersion   CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
                description TEXT COLLATE pg_catalog."default",
                time        TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6),
                source      CHARACTER VARYING(250),
                CONSTRAINT version_app_version_db_release_pk PRIMARY KEY (appVersion, dbVersion)
            ) WITH (OIDS = FALSE) TABLESPACE pg_default;
            ALTER TABLE version OWNER to %v;`, dbUser))
		if err != nil {
			return errors.New(fmt.Sprintf("!!! I cannot create the version table: %v", err))
		}
	}
	// ready to insert a new version
	_, err = conn.Exec(context.Background(),
		fmt.Sprintf(`INSERT INTO "version"(appVersion, dbVersion, description, source) VALUES('%s', '%s', '%s', '%s');`,
			version.AppVersion, version.DbVersion, version.Description, version.Source))
	if err != nil {
		return errors.New(fmt.Sprintf("!!! I cannot update the version table: %v", err))
	}
	return err
}

// get database server information
func (db *PgSQLProvider) GetInfo() (*DbInfo, error) {
	conn, err := db.newConn(true, true)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), `SELECT version()`)
	if err != nil {
		return nil, err
	}
	var info string
	if rows.Next() {
		rows.Scan(&info)
		rows.Close()
	} else {
		// no results
		rows.Close()
		return nil, err
	}
	fmt.Println(info)
	parts := strings.Split(info, ",")
	if len(parts) != 3 {
		return nil, errors.New("!!! I cannot parse database server information: retrieved information not in understandable format")
	}
	o := strings.Split(parts[0], "on")
	return &DbInfo{
		Database:        strings.Trim(o[0], " "),
		OperatingSystem: strings.Trim(o[1], " "),
		Compiler:        strings.Trim(parts[1], " "),
		ProcessorBits:   strings.Trim(parts[2], " "),
	}, nil
}

// return the connection string
// admin:
//  - if true, a connection using the postgres user is returned
//  - if false, a connection using the database user is returned
// database:
//  - if true, adds the database name to the connection string
func (db *PgSQLProvider) connString(admin bool, database bool) (string, error) {
	connStr := ""
	host, found := db.get("Db.Host")
	if !found {
		return "", errors.New(fmt.Sprint("!!! could not find Db.Host config value"))
	}
	port, found := db.get("Db.Port")
	if !found {
		return "", errors.New(fmt.Sprint("!!! could not find Db.Port config value"))
	}
	if admin {
		adminUsername, found := db.get("Db.AdminUsername")
		if !found {
			return "", errors.New(fmt.Sprint("!!! could not find Db.AdminUsername config value"))
		}
		adminPwd, found := db.get("Db.AdminPassword")
		if !found {
			return "", errors.New(fmt.Sprint("!!! could not find Db.AdminPassword config value"))
		}
		connStr = fmt.Sprintf("postgresql://%v:%v@%v:%v", adminUsername, adminPwd, host, port)
	} else {
		username, found := db.get("Db.Username")
		if !found {
			return "", errors.New(fmt.Sprint("!!! could not find Db.Username config value"))
		}
		pwd, found := db.get("Db.Password")
		if !found {
			return "", errors.New(fmt.Sprint("!!! could not find Db.Password config value"))
		}
		connStr = fmt.Sprintf("postgresql://%v:%v@%v:%v", username, pwd, host, port)
	}
	if database {
		dbname, found := db.get("Db.Name")
		if !found {
			return "", errors.New(fmt.Sprint("!!! could not find Db.Name config value"))
		}
		connStr = fmt.Sprintf("%v/%v", connStr, dbname)
	}
	return connStr, nil
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
		// gets the connection string to use
		connStr, e := db.connString(admin, database)
		// if we could not construct a valid connection string
		if e != nil {
			// send the error through the channel
			connect <- conn{conn: nil, err: e}
		}
		// connects to the database
		c, e := pgxpool.Connect(context.Background(), connStr)
		// sends connection through the channel
		connect <- conn{conn: c, err: e}
	}()
	// launch a go routine
	go func() {
		// timeout period is 4 secs
		time.Sleep(4e9)
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

// return an enhanced error
func (db *PgSQLProvider) error(err error) (bool, error) {
	isNull := err == nil
	if !isNull {
		// if the error is a postgres error
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// return detailed info
			var where string
			if len(pgErr.Where) > 0 {
				where = fmt.Sprintf(" - where: %s", pgErr.Where)
			}
			return isNull, errors.New(fmt.Sprintf("%s: %s (SQLSTATE %s)%s", pgErr.Severity, pgErr.Message, pgErr.Code, where))
		} else {
			return isNull, err
		}
	}
	return isNull, err
}

func (db *PgSQLProvider) label(textTrue string, textFalse string, use bool) string {
	if use {
		return textTrue
	} else {
		return textFalse
	}
}

func (db *PgSQLProvider) get(key string) (string, bool) {
	return db.cfg.GetString(key)
}
