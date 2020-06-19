//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// the database provider for PostgreSQL
// NOTE: database providers implicitly implement the DatabaseProvider interface
type PgSQLProvider struct {
	cfg *Config
}

func (db *PgSQLProvider) Setup(config *Config) {
	db.cfg = config
}

// execute the specified getAction
func (db *PgSQLProvider) RunCommand(command *Command) (string, error) {
	log := bytes.Buffer{}
	conn, err := db.newConn(command.AsAdmin, command.UseDb)
	defer conn.Close()
	if err != nil {
		return log.String(), err
	}
	if command.Transactional {
		log.WriteString(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", command.Transactional),
			db.label("as an admin", "as a user", command.AsAdmin),
			db.label("to the db", "to the server", command.UseDb)))
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return log.String(), err
		}
		for _, script := range command.Scripts {
			_, err := tx.Exec(context.Background(), script.Content)
			log.WriteString(fmt.Sprintf("? I have executed the script '%s'\n", script.Name))
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				tx.Rollback(context.Background())
				return log.String(), err
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
				return log.String(), err
			}
		}
	}
	return log.String(), nil
}

func (db *PgSQLProvider) RunQuery(query *Query, params ...interface{}) (Table, error) {
	conn, err := db.newConn(false, true)
	if err != nil {
		return Table{}, nil
	}
	defer conn.Close()
	result, err := conn.Query(context.Background(), query.Content)
	if err != nil {
		return Table{}, nil
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
				row = append(row, string(v))
			}
			if v, ok := value.(time.Time); ok {
				row = append(row, v.String())
			}
		}
		// add the row to the row set
		rows = append(rows, row)
	}
	result.Close()
	return Table{
		Header: header,
		Rows:   rows,
	}, nil
}

func (db *PgSQLProvider) GetVersion() (appVersion string, dbVersion string, err error) {
	conn, err := db.newConn(true, true)
	if err != nil {
		return appVersion, dbVersion, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT appVersion, dbVersion 
		FROM "version"
		ORDER BY time DESC
		LIMIT 1`)
	if err != nil {
		return appVersion, dbVersion, err
	}
	if rows.Next() {
		rows.Scan(&appVersion, &dbVersion)
	}
	rows.Close()
	return appVersion, dbVersion, err
}

// add version
func (db *PgSQLProvider) SetVersion(appVersion string, dbVersion string, description string, source string) error {
	// create a db connection
	conn, err := db.newConn(false, true)
	if err != nil {
		return errors.New(fmt.Sprintf("!!! I cannot connect to the database: %v", err))
	}
	defer conn.Close()
	// find out if version table exists
	_, err = conn.Query(context.Background(), "SELECT * FROM version")
	if err != nil {
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
            ALTER TABLE version OWNER to %v;`, db.cfg.Get(DbUsername)))
		if err != nil {
			return errors.New(fmt.Sprintf("!!! I cannot create the version table: %v", err))
		}
	}
	// ready to insert a new version
	_, err = conn.Exec(context.Background(),
		fmt.Sprintf(`INSERT INTO "version"(appVersion, dbVersion, description, source) VALUES('%s', '%s', '%s', '%s');`,
			appVersion, dbVersion, description, source))
	if err != nil {
		return errors.New(fmt.Sprintf("!!! I cannot update the version table: %v", err))
	}
	return err
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
			db.get(DbAdminUser),
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
