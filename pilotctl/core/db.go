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
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewDb(host, port, db, uname, pwd string) (*Db, error) {
	d := &Db{
		db:    db,
		host:  host,
		uname: uname,
		pwd:   pwd,
		port:  port,
	}
	pool, err := newPool(connStr(uname, pwd, host, port, db))
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
type conn struct {
	conn *pgxpool.Pool
	err  error
}

// create a new database connection pool
// if it cannot connect within 5 seconds, it returns an error
func newPool(connStr string) (*pgxpool.Pool, error) {
	// this channel receives an connection
	connect := make(chan conn, 1)
	// this channel receives a timeout flag
	timeout := make(chan bool, 1)
	// launch a go routine to try the database connection
	go func() {
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
func connStr(uname, pwd, host, port, db string) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", uname, pwd, host, port, db)
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

func (db *Db) RunQuery(query string) (*Table, error) {
	// acquires a database connection
	conn, err := db.pool.Acquire(context.Background())
	// release the connection
	defer conn.Release()
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
				t := toTime(v.Microseconds)
				row = append(row, t)
			} else if v, ok := value.(int64); ok {
				row = append(row, strconv.Itoa(int(v)))
			} else if v, ok := value.(int32); ok {
				n := strconv.Itoa(int(v))
				row = append(row, n)
			} else if v, ok := value.(bool); ok {
				row = append(row, strconv.FormatBool(v))
			} else if v, ok := value.(pgtype.TextArray); ok {
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
