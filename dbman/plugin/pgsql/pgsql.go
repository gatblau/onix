package main

import (
	"context"
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// DbMan's database provider for PostgreSQL
// note: providers implicitly implement the DatabaseProvider interface
type PgSQLProvider struct {
	cfg *Conf
}

// pass DbMan configuration to the database provider
// config: map[string]interface{} serialised as a json string
func (db *PgSQLProvider) Setup(config string) string {
	// parse the configuration
	c, output := NewConf(config)
	// allocate the parsed object to cfg
	db.cfg = c
	// return the output
	return output
}

// retrieve database information
// return: map[string]interface{} serialised as a json string
func (db *PgSQLProvider) GetVersion() string {
	output := NewParameter()
	conn, err := db.newConn(true, true)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(), `
		SELECT appVersion, dbVersion, description, time, source
		FROM "version"
		ORDER BY time DESC
		LIMIT 1`)
	if err != nil {
		return output.ToError(err)
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
		output.SetErrorFromMessage("!!! query for version returned no results")
		return output.ToError(err)
	}
	// we have version information so populate result
	output.Set("appVersion", appVersion)
	output.Set("dbVersion", dbVersion)
	output.Set("description", description)
	output.Set("time", time)
	output.Set("source", source)
	// return the result as a JSON string
	return output.ToString()
}

func (db *PgSQLProvider) RunCommand(command string) string {
	output := NewParameter()
	cmd, err := NewCommand(command)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	conn, err := db.newConn(cmd.AsAdmin, cmd.UseDb)
	defer conn.Close()
	if err != nil {
		return output.ToError(err)
	}
	if cmd.Transactional {
		output.Log(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", cmd.Transactional),
			db.label("as an admin", "as a user", cmd.AsAdmin),
			db.label("to the db", "to the server", cmd.UseDb)))
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return output.ToError(err)
		}
		for _, script := range cmd.Scripts {
			_, err := tx.Exec(context.Background(), script.Content)
			output.Log(fmt.Sprintf("? I have executed the script '%s'\n", script.Name))
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				tx.Rollback(context.Background())
				return output.ToError(err)
			}
		}
		tx.Commit(context.Background())
	} else {
		output.Log(fmt.Sprintf("? I am creating a db connection that is %v, %v and %v\n",
			db.label("transactional", "non-transactional", cmd.Transactional),
			db.label("as an admin", "as a user", cmd.AsAdmin),
			db.label("to the db", "to the server", cmd.UseDb)))
		for _, script := range cmd.Scripts {
			_, err := conn.Exec(context.Background(), script.Content)
			// if we have an error
			if isNull, err := db.error(err); !isNull {
				return output.ToError(err)
			}
		}
	}
	return output.ToString()
}

func (db *PgSQLProvider) RunQuery(queryInfo string) string {
	output := NewParameter()
	query, err := NewQuery(queryInfo)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	conn, err := db.newConn(false, true)
	if err != nil {
		return output.ToError(err)
	}
	defer conn.Close()
	result, err := conn.Query(context.Background(), query.Content)
	if err != nil {
		return output.ToError(err)
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
	output.Set("table", Table{
		Header: header,
		Rows:   rows,
	})
	return output.ToString()
}

// set the release version
// versionInfo: a json serialised map[string]interface{} containing version information
func (db *PgSQLProvider) SetVersion(versionInfo string) string {
	output := NewParameter()
	input := NewParameterFromJSON(versionInfo)
	appVersion := input.GetString("appVersion")
	dbVersion := input.GetString("dbVersion")
	description := input.GetString("description")
	source := input.GetString("source")

	// create a db connection
	conn, err := db.newConn(false, true)
	if err != nil {
		return output.ToError(err)
	}
	defer conn.Close()
	// find out if version table exists
	_, err = conn.Query(context.Background(), "SELECT * FROM version")
	if err != nil {
		dbUser, found := db.get("Db.Username")
		if !found {
			return output.ToError(errors.New(fmt.Sprint("!!! could not find Db.Username config value")))
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
			return output.ToError(errors.New(fmt.Sprintf("!!! I cannot create the version table: %v", err)))
		}
	}
	// ready to insert a new version
	_, err = conn.Exec(context.Background(),
		fmt.Sprintf(`INSERT INTO "version"(appVersion, dbVersion, description, source) VALUES('%s', '%s', '%s', '%s');`,
			appVersion, dbVersion, description, source))
	if err != nil {
		return output.ToError(errors.New(fmt.Sprintf("!!! I cannot update the version table: %v", err)))
	}
	return output.ToString()
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
