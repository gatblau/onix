package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"time"

	// "github.com/jackc/pgconn"
	// "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

func NewDb(host, port, db, uname, pwd string) *Db {
	return &Db{
		db:    db,
		host:  host,
		uname: uname,
		pwd:   pwd,
		port:  port,
	}
}

type Db struct {
	db    string
	host  string
	uname string
	pwd   string
	port  string
}

// this type carries either a connection or an error
// used in a channel used by the connection go routine
type conn struct {
	conn *pgxpool.Pool
	err  error
}

// create a new database connection
// if it cannot connect within 5 seconds, it returns an error
func (db *Db) newConn(admin bool, database bool) (*pgxpool.Pool, error) {
	// this channel receives an connection
	connect := make(chan conn, 1)
	// this channel receives a timeout flag
	timeout := make(chan bool, 1)
	// launch a go routine to try the database connection
	go func() {
		// gets the connection string to use
		connStr := db.connString()
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

// return the connection string
func (db *Db) connString() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", db.uname, db.pwd, db.host, db.port, db.db)
}

// return an enhanced error
func (db *Db) error(err error) (bool, error) {
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
