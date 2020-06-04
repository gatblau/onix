//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

const PgSQLVersionTable = `
DO
  $$
    BEGIN
      ---------------------------------------------------------------------------
      -- VERSION - version of releases (not only database)
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'version')
      THEN
        CREATE TABLE version
        (
          application_version CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
          database_version    CHARACTER VARYING(25) NOT NULL COLLATE pg_catalog."default",
          description         TEXT COLLATE pg_catalog."default",
          time                timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
          scripts_source      character varying(250),
          CONSTRAINT version_app_version_db_release_pk PRIMARY KEY (application_version, database_version)
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE version
          OWNER to onix;
      END IF;
    END;
    $$
`

// default config file content
const cfgFile = `
# verbosity of logging (Trace, Debug, Warning, Info, Error, Fatal, Panic)
LogLevel = "Warning"

# configuration for running DbMan in http mode
[Http]
	Metrics = "true"
	AuthMode    = "basic"
	Port        = "8085"
	Username    = "admin"
	Password    = "0n1x"

# configuration for the Onix Web API integration
[Db]
    Provider    = "pgsql"
    Name        = "onix"
    Host        = "localhost"
    Port        = "5432"
    Username    = "onix"
    Password    = "onix"
    AdminPwd    = "onix"

# configuration of database scripts remote repository
[Schema]
    URI         = "https://raw.githubusercontent.com/gatblau/ox-db/master"
    Username    = ""
    Token       = ""
`
