package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strings"
	"time"
)

func NewDb(host, port, db, uname, pwd string, maxConn int) (*Db, error) {
	d := &Db{
		db:    db,
		host:  host,
		uname: uname,
		pwd:   pwd,
		port:  port,
	}
	pool, err := newPool(connStr(uname, pwd, host, port, db, maxConn))
	if err != nil {
		return nil, err
	}
	d.pool = pool
	return d, nil
}

type Db struct {
	db    string
	host  string
	uname string
	pwd   string
	port  string
	pool  *pgxpool.Pool
}

// this type carries either a connection or an error
// used in a channel used by the connection go routine
type connection struct {
	conn *pgxpool.Pool
	err  error
}

// create a new database connection pool
// if it cannot connect within a given period, it returns an error
func newPool(connStr string) (*pgxpool.Pool, error) {
	// this channel receives a connection
	connect := make(chan connection, 1)
	// this channel receives a timeout flag
	timeout := make(chan bool, 1)

	// launch a go routine to try the database connection
	go func() {
		// connects to the database
		c, e := pgxpool.Connect(context.Background(), connStr)
		// sends connection through the channel
		connect <- connection{conn: c, err: e}
	}()
	// launch a go routine
	go func() {
		// timeout period is 2 minutes
		time.Sleep(120 * time.Second)
		timeout <- true
	}()

	select {
	// the connection has been established before the timeout
	case c := <-connect:
		{
			if c.err != nil {
				return nil, c.err
			}
			return c.conn, nil
		}
	// the connection has not yet returned when the timeout happens
	case <-timeout:
		{
			return nil, errors.New("cannot connect to pilotctl database, the timed out period has elapsed\n")
		}
	}
}

// return the connection string
func connStr(uname, pwd, host, port, db string, maxConn int) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?pool_max_conns=%d", uname, pwd, host, port, db, maxConn)
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

func (db *Db) RunCommand(script string, arguments ...interface{}) error {
	// acquires a database connection
	conn, err := db.pool.Acquire(context.Background())
	// if cannot connect to the server return with the error
	if err != nil {
		return err
	}
	// release the connection
	defer conn.Release()
	// if the command is to be run within a database transaction
	// acquires a db transaction
	tx, err := conn.Begin(context.Background())
	// if error then return
	if err != nil {
		return err
	}
	// execute the content of the script
	_, err = tx.Exec(context.Background(), script, arguments...)
	// if we have an error return it
	if isNull, err := db.error(err); !isNull {
		// rollback the transaction
		tx.Rollback(context.Background())
		// return the error
		return err
	}
	// all good so commit the transaction
	tx.Commit(context.Background())
	// return the execution log
	return err
}

func (db *Db) Query(query string, args ...interface{}) (pgx.Rows, error) {
	// acquires a database connection
	conn, err := db.pool.Acquire(context.Background())
	// release the connection
	defer conn.Release()
	// if error then return it
	if err != nil {
		return nil, err
	}
	// execute the query content
	return conn.Query(context.Background(), query, args...)
}

// Table generic table used as a serializable result set for queries
type Table struct {
	Header Row   `json:"header,omitempty"`
	Rows   []Row `json:"row,omitempty"`
}

// Row a row in the table
type Row []string

func toHStore(m map[string]string) pgtype.Hstore {
	hstore := pgtype.Hstore{
		Status: pgtype.Present,
	}
	text := func(s string) pgtype.Text {
		return pgtype.Text{String: s, Status: pgtype.Present}
	}
	for key, value := range m {
		hstore.Map[key] = text(value)
	}
	return hstore
}

func toHStoreString(m map[string]string) string {
	sb := new(strings.Builder)
	if len(m) == 0 {
		return ""
	} else {
		for key, value := range m {
			appendEscaped(sb, key)
			sb.WriteString("=>")
			appendEscaped(sb, value)
			sb.WriteString(", ")
		}
	}
	str := sb.String()
	return str[:len(str)-2]
}

func appendEscaped(sb *strings.Builder, val interface{}) {
	if val != nil {
		sb.WriteString("\"")
		str, ok := val.(string)
		if !ok {
			log.Printf("WARNING: HStore value could not be cast as string, value might not be persisted correctly in the data source\n")
		}
		for pos := 0; pos < len(str); pos++ {
			if str[pos] == '"' || str[pos] == '\\' {
				sb.WriteString("\\")
			}
			sb.Write([]byte{str[pos]})
		}
		sb.WriteString("\"")
	} else {
		sb.WriteString("NULL")
	}
}

func fromHStoreString(value string) map[string]string {
	m := make(map[string]string)
	parts := strings.Split(value, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=>")
		key := strings.TrimPrefix(kv[0], " ")
		key = strings.TrimSuffix(key, " ")
		key = key[1 : len(key)-1]
		value := strings.TrimPrefix(kv[1], " ")
		value = strings.TrimSuffix(value, " ")
		value = value[1 : len(value)-1]
		m[key] = value
	}
	return m
}
