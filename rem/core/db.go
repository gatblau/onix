package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"strconv"
	"time"
)

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
func (db *Db) newConn() (*pgxpool.Pool, error) {
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

func (db *Db) RunCommand(scripts []string) (bytes.Buffer, error) {
	// create a buffer to write execution output to be passed back to client
	// use this instead of writing to stdout
	log := bytes.Buffer{}
	// acquires a database connection
	conn, err := db.newConn()
	// if cannot connect to the server return with the error
	if err != nil {
		return log, err
	}
	// if the command is to be run within a database transaction
	// acquires a db transaction
	tx, err := conn.Begin(context.Background())
	// if error then return
	if err != nil {
		return log, err
	}
	// for each database script in the command
	for _, script := range scripts {
		// execute the content of the script
		_, err := tx.Exec(context.Background(), script)
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
	// return the execution log
	return log, err
}

func (db *Db) RunQuery(query string) (*Table, error) {
	// acquires a database connection
	conn, err := db.newConn()
	// if cannot connect to the server return with the error
	if err != nil {
		return nil, err
	}
	// execute the query content
	result, err := conn.Query(context.Background(), query)
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
				if headerName == "?column?" { // the query has not defined a column name
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
			} else if v, ok := value.(bool); ok {
				row = append(row, strconv.FormatBool(v))
			} else if v, ok := value.([]string); ok {
				row = append(row, toCSV(v))
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

func toCSV(v []string) string {
	str := bytes.Buffer{}
	for i, s := range v {
		str.WriteString(s)
		if i < len(v) {
			str.WriteString(",")
		}
	}
	return str.String()
}

// converts microseconds into HH:mm:SS.ms
func (db *Db) toTime(microseconds int64) string {
	milliseconds := (microseconds / 1000) % 1000
	seconds := (((microseconds / 1000) - milliseconds) / 1000) % 60
	minutes := (((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) % 60
	hours := ((((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) - minutes) / 60
	return fmt.Sprintf("%02v:%02v:%02v.%03v", hours, minutes, seconds, milliseconds)
}

// Table generic table used as a serializable result set for queries
type Table struct {
	Header Row   `json:"header,omitempty"`
	Rows   []Row `json:"row,omitempty"`
}

// Row a row in the table
type Row []string
