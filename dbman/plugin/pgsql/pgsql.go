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
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Implementation of DbMan's database provider for PostgreSQL
// NOTE:
//  - PgSQLProvider implicitly implements the DatabaseProvider interface
//  - the provider uses connection pooling so should not call conn.Close()
type PgSQLProvider struct {
	cfg *Conf
}

// pass DbMan configuration to the database provider
// config: configuration passed-in by DbMan to the plugin
func (db *PgSQLProvider) Setup(config *Conf) error {
	if config != nil {
		// allocate the parsed object to cfg
		db.cfg = config
		// return without error
		return nil
	}
	return errors.New("!!! the PostgreSQL Database plugin was not provided with a valid configuration\n")
}

// this function retrieves database version information in a Version struct
// Version: the current database version
// error: if failed to retrieve the version from the database
func (db *PgSQLProvider) GetVersion() (*Version, error) {
	// connect to the database server
	conn, err := db.newConn(true, true)
	// if the connection failed return the error
	if err != nil {
		return nil, err
	}
	// query the database version table
	rows, err := conn.Query(context.Background(), `
		SELECT appVersion, dbVersion, description, time, source
		FROM "version"
		ORDER BY time DESC
		LIMIT 1`)
	// if the query failed return the error
	if err != nil {
		return nil, err
	}
	var (
		appVersion, dbVersion, description, source string
		time                                       time.Time
	)
	// reads the query return values
	if rows.Next() {
		rows.Scan(&appVersion, &dbVersion, &description, &time, &source)
		rows.Close()
	} else {
		// no results
		rows.Close()
		return nil, err
	}
	// populate the Version struct with the query returned values
	v := &Version{
		AppVersion:  appVersion,
		DbVersion:   dbVersion,
		Description: description,
		Source:      source,
		Time:        time,
	}
	// return the version
	return v, err
}

// this function runs a database command
// command: the struct containing the information required to run the command
func (db *PgSQLProvider) RunCommand(command *Command) (bytes.Buffer, error) {
	// create a buffer to write execution output to be passed back to DbMan
	// use this instead of writing to stdout
	log := bytes.Buffer{}
	// acquires a database connection
	conn, err := db.newConn(command.AsAdmin, command.UseDb)
	// if cannot connect to the server return with the error
	if err != nil {
		return log, err
	}
	// if the command is to be run within a database transaction
	if command.Transactional {
		// log the db connection creation step
		log.WriteString(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", command.Transactional),
			db.label("as an admin", "as a user", command.AsAdmin),
			db.label("to the db", "to the server", command.UseDb)))
		// acquires a db transaction
		tx, err := conn.Begin(context.Background())
		// if error then return
		if err != nil {
			return log, err
		}
		// for each database script in the command
		for _, script := range command.Scripts {
			// execute the content of the script
			_, err := tx.Exec(context.Background(), script.Content)
			// log the execution step
			log.WriteString(fmt.Sprintf("? I have executed the script '%s'\n", script.Name))
			// if we have an error return it
			if isNull, err := db.error(err); !isNull {
				// rollback the transaction
				tx.Rollback(context.Background())
				// return the error
				return log, err
			}
		}
		// all good so commit the transaction
		tx.Commit(context.Background())
	} else {
		// log the db connection creation step
		log.WriteString(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", command.Transactional),
			db.label("as an admin", "as a user", command.AsAdmin),
			db.label("to the db", "to the server", command.UseDb)))
		// for each database script in the command
		for _, script := range command.Scripts {
			// execute the content of the script
			_, err := conn.Exec(context.Background(), script.Content)
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				return log, err
			}
		}
	}
	// return the execution log
	return log, err
}

// this function runs a database query
// query: the struct containing the information required to run the query
func (db *PgSQLProvider) RunQuery(query *Query) (*Table, error) {
	// acquires a database connection
	conn, err := db.newConn(false, true)
	// if cannot connect to the server return with the error
	if err != nil {
		return nil, err
	}
	// execute the query content
	result, err := conn.Query(context.Background(), query.Content)
	// if error then return it
	if err != nil {
		return nil, err
	}
	// puts together a generic table result
	header := make(Row, 0) // the table header
	rows := make([]Row, 0) // a slice of table rows
	// for each row in the result
	for result.Next() {
		// only the first time round populate the header
		if len(header) == 0 {
			// for each field in the result set
			for _, desc := range result.FieldDescriptions() {
				headerName := string(desc.Name)
				if headerName == "?column?" { // the query jas not defined a column name
					headerName = "undefined"
				}
				// add a new header
				header = append(header, headerName)
			}
		}
		// create a new row
		row := make(Row, 0)
		// populate the row with returned values from the query
		values, err := result.Values()
		// if error return it
		if err != nil {
			return nil, err
		}
		// for each record in the result set
		for _, value := range values {
			// if the value is convertible to string
			if v, ok := value.(string); ok {
				// append the value to the row slice
				row = append(row, v)
			} else
			// in the case of time values
			if v, ok := value.(time.Time); ok {
				// append the string representation of the value to the row slice
				row = append(row, v.String())
			} else if v, ok := value.(pgtype.Interval); ok {
				t := db.toTime(v.Microseconds)
				row = append(row, t)
			} else if v, ok := value.(int32); ok {
				n := strconv.Itoa(int(v))
				row = append(row, n)
			} else {
				valueType := reflect.TypeOf(value)
				if valueType != nil {
					row = append(row, fmt.Sprintf("unsupported type '%s.%s'", valueType.PkgPath(), valueType.Name()))
				}
			}
		}
		// add the row to the row set
		rows = append(rows, row)
	}
	// closes the result set
	result.Close()
	// return an instance of the generic table populated with the header and rows
	return &Table{
		Header: header,
		Rows:   rows,
	}, err
}

// this function sets the version in the database
// version: struct containing version information to persist in the database
func (db *PgSQLProvider) SetVersion(version *Version) error {
	// create a db connection
	conn, err := db.newConn(false, true)
	// if error then return it
	if err != nil {
		return err
	}
	// find out if the version table exists in the database
	_, err = conn.Query(context.Background(), "SELECT * FROM version")
	// an error indicates the table does not exist, and therefore attempts to create it
	if err != nil {
		// create the version table
		err2 := db.createVersionTable(conn)
		// if error return it
		if err2 != nil {
			return err2
		}
	}
	// insert a entry in the version table
	_, err = conn.Exec(context.Background(),
		fmt.Sprintf(`INSERT INTO "version"(appVersion, dbVersion, description, source) VALUES('%s', '%s', '%s', '%s');`,
			version.AppVersion, version.DbVersion, version.Description, version.Source))
	// if error return it
	if err != nil {
		return errors.New(fmt.Sprintf("!!! I cannot update the version table: %v\n", err))
	}
	return err
}

// query information about the database server and returns it as a DbInfo struct
func (db *PgSQLProvider) GetInfo() (*DbInfo, error) {
	// acquires a database connection
	conn, err := db.newConn(true, true)
	// if error returns it
	if err != nil {
		return nil, err
	}
	// query database server information
	rows, err := conn.Query(context.Background(), `SELECT version()`)
	// if error returns it
	if err != nil {
		return nil, err
	}
	// retrieve the query result into the info var
	var info string
	if rows.Next() {
		rows.Scan(&info)
		rows.Close()
	} else {
		// no results
		rows.Close()
		return nil, err
	}
	// parses the result
	parts := strings.Split(info, ",")
	if len(parts) != 3 {
		return nil, errors.New("!!! I cannot parse database server information: retrieved information not in understandable format")
	}
	o := strings.Split(parts[0], "on")
	// return the dbinfo struct
	return &DbInfo{
		Database:        strings.Trim(o[0], " "),
		OperatingSystem: strings.Trim(o[1], " "),
		Compiler:        strings.Trim(parts[1], " "),
		ProcessorBits:   strings.Trim(parts[2], " "),
	}, nil
}

// =========================================================================
// UTILITY FUNCTIONS
// =========================================================================

// creates a table to hold version information in the database
func (db *PgSQLProvider) createVersionTable(conn *pgxpool.Pool) error {
	// need to get the database user from DbMan's configuration so that can grant ownership
	// to the created table below
	dbUser, found := db.get("Db.Username")
	// if not found, return error
	if !found {
		return errors.New(fmt.Sprint("!!! could not find Db.Username config value\n"))
	}
	// try and create the version table in the database
	// the table structure is defined by the specific database provider so that it contains the information
	// in the Version struct
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
	// if error return it
	if err != nil {
		return errors.New(fmt.Sprintf("!!! I cannot create the version table: %v\n", err))
	}
	return nil
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

// converts microseconds into HH:mm:SS.ms
func (db *PgSQLProvider) toTime(microseconds int64) string {
	milliseconds := (microseconds / 1000) % 1000
	seconds := (((microseconds / 1000) - milliseconds) / 1000) % 60
	minutes := (((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) % 60
	hours := ((((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) - minutes) / 60
	return fmt.Sprintf("%02v:%02v:%02v.%03v", hours, minutes, seconds, milliseconds)
}
