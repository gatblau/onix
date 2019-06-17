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