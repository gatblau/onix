/*
    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software distributed under
    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
    either express or implied.
    See the License for the specific language governing permissions and limitations under the License.

    Contributors to this project, hereby assign copyright in this code to the project,
    to be licensed under the same terms as the rest of the code.
*/
DO
  $$
    BEGIN
      ---------------------------------------------------------------------------
      -- USER
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'user')
      THEN
        CREATE SEQUENCE user_id_seq
          INCREMENT 1
          START 100
          MINVALUE 100
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE user_id_seq
          OWNER TO onix;

        CREATE TABLE "user"
        (
          id          BIGINT                 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) NOT NULL COLLATE pg_catalog."default",
          email       CHARACTER VARYING(200) NOT NULL COLLATE pg_catalog."default",
          pwd         CHARACTER VARYING(300) COLLATE pg_catalog."default",
          salt        CHARACTER VARYING(300) COLLATE pg_catalog."default",
          expires     TIMESTAMP(6) WITH TIME ZONE,
          version     BIGINT                 NOT NULL DEFAULT 1,
          created     TIMESTAMP(6) WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP(6),
          updated     TIMESTAMP(6) WITH TIME ZONE,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT user_id_pk PRIMARY KEY (id),
          CONSTRAINT user_key_uc UNIQUE (key),
          CONSTRAINT user_name_uc UNIQUE (name),
          CONSTRAINT user_email_uc UNIQUE (email),
          CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE "user"
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- USER CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'user_change')
      THEN
        CREATE TABLE user_change
        (
          operation   CHAR(1)                NOT NULL,
          changed     TIMESTAMP              NOT NULL,
          id          BIGINT,
          key         CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          email       CHARACTER VARYING(200) COLLATE pg_catalog."default",
          pwd         CHARACTER VARYING(300) COLLATE pg_catalog."default",
          salt        CHARACTER VARYING(300) COLLATE pg_catalog."default",
          expires     TIMESTAMP(6) WITH TIME ZONE,
          version     BIGINT,
          created     TIMESTAMP(6) WITH TIME ZONE,
          updated     TIMESTAMP(6) WITH TIME ZONE,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        CREATE OR REPLACE FUNCTION ox_change_user() RETURNS TRIGGER AS
        $user_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO user_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO user_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO user_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $user_change$ LANGUAGE plpgsql;

        CREATE TRIGGER user_change
          AFTER INSERT OR UPDATE OR DELETE
          ON "user"
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_user();

        ALTER TABLE user_change
          OWNER to onix;
      END IF;

      -- generic users - should change passwords after first login
      INSERT INTO "user"(id, key, name, email, pwd, salt, version, changed_by)
        VALUES (1, 'admin', 'Administrator', 'admin@onix.com', 'E2BgmQs4vH4rYvj5Fe0p9DbZUKU=', '8DZMiAR+XGA=', 1, 'onix');
      INSERT INTO "user"(id, key, name, email, pwd, salt, version, changed_by)
        VALUES (2, 'reader', 'Reader', 'reader@onix.com', '/EvDpP8kHkfd30mXk+Ne9aA4h5o=', 'B0zo+y0Keiw=', 1, 'onix');
      INSERT INTO "user"(id, key, name, email, pwd, salt, version, changed_by)
        VALUES (3, 'writer', 'Writer', 'writer@onix.com', 'DkV3uMWjAjHSTZnW9TkJNI6XOzU=', 'yWJm38+RPtc=', 1, 'onix');

      ---------------------------------------------------------------------------
      -- PARTITION
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'partition')
      THEN
        CREATE SEQUENCE partition_id_seq
          INCREMENT 1
          START 100
          MINVALUE 100
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE partition_id_seq
          OWNER TO onix;

        CREATE TABLE partition
        (
          id          bigint                 NOT NULL DEFAULT nextval('partition_id_seq'::regclass),
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          version     bigint                 NOT NULL DEFAULT 1,
          created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          owner       CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default" DEFAULT 'ADMIN',
          CONSTRAINT partition_id_pk PRIMARY KEY (id),
          CONSTRAINT partition_key_uc UNIQUE (key),
          CONSTRAINT partition_name_uc UNIQUE (name)
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE partition
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- PARTITION CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'partition_change')
      THEN
        CREATE TABLE partition_change
        (
          operation   CHAR(1)                NOT NULL,
          changed     TIMESTAMP              NOT NULL,
          id          bigint,
          key         CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          version     bigint,
          created     timestamp(6) with time zone,
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          owner       CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        CREATE OR REPLACE FUNCTION ox_change_partition() RETURNS TRIGGER AS
        $partition_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO partition_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO partition_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO partition_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $partition_change$ LANGUAGE plpgsql;

        CREATE TRIGGER partition_change
          AFTER INSERT OR UPDATE OR DELETE
          ON partition
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_partition();

        ALTER TABLE partition_change
          OWNER to onix;
      END IF;

      INSERT INTO partition(id, key, name, description, version, changed_by)
      VALUES (0, 'REF', 'Default Reference Partition', 'Default partition for reference data.', 1, 'onix');
      INSERT INTO partition(id, key, name, description, version, changed_by)
      VALUES (1, 'INS', 'Default Instance Partition', 'Default partition for instance data.', 1, 'onix');

      ---------------------------------------------------------------------------
      -- ROLE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'role')
      THEN
        CREATE SEQUENCE role_id_seq
          INCREMENT 1
          START 100
          MINVALUE 100
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE role_id_seq
          OWNER TO onix;

        CREATE TABLE role
        (
          id          bigint                 NOT NULL DEFAULT nextval('role_id_seq'::regclass),
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          version     bigint                 NOT NULL DEFAULT 1,
          created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          level       integer                         default 0,
          owner       CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default" DEFAULT 'ADMIN',
          CONSTRAINT role_id_pk PRIMARY KEY (id),
          CONSTRAINT role_key_uc UNIQUE (key),
          CONSTRAINT role_name_uc UNIQUE (name)
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE role
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- ROLE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'role_change')
      THEN
        CREATE TABLE role_change
        (
          operation   CHAR(1)                NOT NULL,
          changed     TIMESTAMP              NOT NULL,
          id          bigint,
          key         CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          version     bigint,
          created     timestamp(6) with time zone,
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          level       integer,
          owner       CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        CREATE OR REPLACE FUNCTION ox_change_role() RETURNS TRIGGER AS
        $role_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO role_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO role_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO role_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $role_change$ LANGUAGE plpgsql;

        CREATE TRIGGER role_change
          AFTER INSERT OR UPDATE OR DELETE
          ON role
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_role();

        ALTER TABLE role_change
          OWNER to onix;
      END IF;

      INSERT INTO role(id, key, name, description, version, changed_by, level)
      VALUES (1, 'ADMIN', 'System Administrator', 'Can read and write configuration data models.', 1, 'onix', 2);
      INSERT INTO role(id, key, name, description, version, changed_by)
      VALUES (2, 'READER', 'System Reader', 'Can only read configuration data and models.', 1, 'onix');
      INSERT INTO role(id, key, name, description, version, changed_by)
      VALUES (3, 'WRITER', 'System Writer', 'Can read and write configuration data and read models.', 1, 'onix');

      ---------------------------------------------------------------------------
      -- PRIVILEGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'privilege')
      THEN
        CREATE SEQUENCE privilege_id_seq
          INCREMENT 1
          START 10
          MINVALUE 10
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE privilege_id_seq OWNER TO onix;

        CREATE TABLE privilege
        (
          id           bigint                 NOT NULL DEFAULT nextval('privilege_id_seq'::regclass),
          key          CHARACTER VARYING(100) NOT NULL,
          role_id      bigint,
          partition_id bigint,
          can_create   boolean,
          can_read     boolean,
          can_delete   boolean,
          version      bigint,
          created      timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT privilege_id_pk PRIMARY KEY (id, role_id, partition_id),
          CONSTRAINT privilege_role_id_fk FOREIGN KEY (role_id)
            REFERENCES role (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT privilege_partition_id_fk FOREIGN KEY (partition_id)
            REFERENCES partition (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (OIDS = FALSE)
          TABLESPACE pg_default;

        ALTER TABLE privilege
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- PRIVILEGE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'privilege_change')
      THEN
        CREATE TABLE privilege_change
        (
          operation    CHAR(1)   NOT NULL,
          changed      TIMESTAMP NOT NULL,
          id           INTEGER   NOT NULL,
          key          CHARACTER VARYING(100),
          role_id      bigint,
          partition_id bigint,
          can_create   boolean,
          can_read     boolean,
          can_delete   boolean,
          version      bigint,
          created      timestamp(6) with time zone,
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)
        );

        CREATE OR REPLACE FUNCTION ox_change_privilege() RETURNS TRIGGER AS
        $privilege_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO privilege_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO privilege_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO privilege_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $privilege_change$ LANGUAGE plpgsql;

        CREATE TRIGGER privilege_change
          AFTER INSERT OR UPDATE OR DELETE
          ON privilege
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_privilege();

        ALTER TABLE privilege_change
          OWNER to onix;
      END IF;

      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (1, 'ADMIN-REF', 1, 0, true, true, true, 'onix', 1); -- admin privilege on part 0
      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (2, 'ADMIN-INS', 1, 1, true, true, true, 'onix', 1); -- admin privilege on part 1
      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (3, 'READER-REF', 2, 0, false, true, false, 'onix', 1); -- reader privilege on part 0
      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (4, 'READER-INS', 2, 1, false, true, false, 'onix', 1); -- reader privilege on part 1
      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (5, 'WRITER-REF', 3, 0, false, true, false, 'onix', 1); -- writer privilege on part 0
      INSERT INTO privilege(id, key, role_id, partition_id, can_create, can_read, can_delete, changed_by, version)
      VALUES (6, 'WRITER-INS', 3, 1, true, true, true, 'onix', 1);

      ---------------------------------------------------------------------------
      -- MEMBERSHIP
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'membership')
      THEN
        CREATE SEQUENCE membership_id_seq
          INCREMENT 1
          START 10
          MINVALUE 10
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE membership_id_seq OWNER TO onix;

        CREATE TABLE membership
        (
          id           bigint                 NOT NULL DEFAULT nextval('membership_id_seq'::regclass),
          key          CHARACTER VARYING(100) NOT NULL,
          user_id      bigint,
          role_id      bigint,
          version      bigint,
          created      timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT membership_id_pk PRIMARY KEY (id, user_id, role_id),
          CONSTRAINT membership_role_id_fk FOREIGN KEY (role_id)
            REFERENCES role (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT membership_user_id_fk FOREIGN KEY (user_id)
            REFERENCES "user" (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (OIDS = FALSE)
          TABLESPACE pg_default;

        ALTER TABLE membership
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- MEMBERSHIP CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'membership_change')
      THEN
        CREATE TABLE membership_change
        (
          operation    CHAR(1)   NOT NULL,
          changed      TIMESTAMP NOT NULL,
          id           INTEGER   NOT NULL,
          key          CHARACTER VARYING(100),
          user_id      bigint,
          role_id      bigint,
          version      bigint,
          created      timestamp(6) with time zone,
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)
        );

        CREATE OR REPLACE FUNCTION ox_change_membership() RETURNS TRIGGER AS
        $membership_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO membership_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO membership_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO membership_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $membership_change$ LANGUAGE plpgsql;

        CREATE TRIGGER membership_change
          AFTER INSERT OR UPDATE OR DELETE
          ON membership
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_membership();

        ALTER TABLE membership_change
          OWNER to onix;
      END IF;

      INSERT INTO membership(id, key, user_id, role_id, changed_by, version)
        VALUES (1, 'ADMIN-MEMBER', 1, 1, 'onix', 1);
      INSERT INTO membership(id, key, user_id, role_id, changed_by, version)
        VALUES (1, 'READER-MEMBER', 2, 2, 'onix', 1);
      INSERT INTO membership(id, key, user_id, role_id, changed_by, version)
        VALUES (1, 'WRITER-MEMBER', 3, 3, 'onix', 1);

      ---------------------------------------------------------------------------
      -- MODEL
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'model')
      THEN
        CREATE SEQUENCE model_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE model_id_seq
          OWNER TO onix;

        CREATE TABLE model
        (
          id           INTEGER                NOT NULL DEFAULT nextval('model_id_seq'::regclass),
          key          CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name         CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description  TEXT COLLATE pg_catalog."default",
          managed      BOOLEAN NOT NULL DEFAULT FALSE,
          version      bigint                 NOT NULL DEFAULT 1,
          created      timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          partition_id bigint                 NOT NULL DEFAULT 0,
          CONSTRAINT model_id_pk PRIMARY KEY (id),
          CONSTRAINT model_key_uc UNIQUE (key),
          CONSTRAINT model_name_uc UNIQUE (name),
          CONSTRAINT model_partition_id_fk FOREIGN KEY (partition_id)
            REFERENCES partition (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE model
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- MODEL CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'model_change')
      THEN
        CREATE TABLE model_change
        (
          operation    CHAR(1)                NOT NULL,
          changed      TIMESTAMP              NOT NULL,
          id           INTEGER,
          key          CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name         CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description  TEXT COLLATE pg_catalog."default",
          managed      BOOLEAN,
          version      BIGINT,
          created      timestamp(6) with time zone,
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          partition_id bigint default 0
        );

        CREATE OR REPLACE FUNCTION ox_change_model() RETURNS TRIGGER AS
        $model_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO model_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO model_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO model_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $model_change$ LANGUAGE plpgsql;

        CREATE TRIGGER model_change
          AFTER INSERT OR UPDATE OR DELETE
          ON model
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_model();

        ALTER TABLE model_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- ITEM TYPE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'item_type')
      THEN
        CREATE SEQUENCE item_type_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE item_type_id_seq
          OWNER TO onix;

        CREATE TABLE item_type
        (
            id            INTEGER                NOT NULL DEFAULT nextval('item_type_id_seq'::regclass),
            key           CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            name          CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description   TEXT COLLATE pg_catalog."default",
            filter        JSONB,
            meta_schema   JSONB,
            version       BIGINT                 NOT NULL DEFAULT 1,
            created       TIMESTAMP(6) WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP(6),
            updated       TIMESTAMP(6) WITH TIME ZONE,
            changed_by    CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            model_id      INT                    NOT NULL,
            notify_change CHAR(1)                NOT NULL DEFAULT 'N' CHECK (notify_change IN ('N', 'T', 'I')), -- N: no, T: yes to a common topic by type, I:yes to a dedicated topic by instance
            tag           TEXT[] COLLATE pg_catalog."default",
            encrypt_meta  BOOLEAN                NOT NULL DEFAULT FALSE,
            encrypt_txt   BOOLEAN                NOT NULL DEFAULT FALSE,
            style         JSONB,
            CONSTRAINT item_type_id_pk PRIMARY KEY (id),
            CONSTRAINT item_type_key_uc UNIQUE (key),
            CONSTRAINT item_type_name_uc UNIQUE (name),
            CONSTRAINT item_type_model_id_fk FOREIGN KEY (model_id)
                REFERENCES model (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE CASCADE
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        CREATE INDEX fki_item_type_model_id_fk
          ON item_type USING btree (model_id)
          TABLESPACE pg_default;

        CREATE INDEX item_type_tag_ix
          ON item_type USING gin (tag COLLATE pg_catalog."default")
          TABLESPACE pg_default;

        ALTER TABLE item_type
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- ITEM_TYPE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'item_type_change')
      THEN
        CREATE TABLE item_type_change
        (
            operation     CHAR(1)                NOT NULL,
            changed       TIMESTAMP              NOT NULL,
            id            INTEGER,
            key           CHARACTER VARYING(100) COLLATE pg_catalog."default",
            name          CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description   TEXT COLLATE pg_catalog."default",
            filter        jsonb,
            meta_schema   jsonb,
            version       bigint,
            created       timestamp(6) with time zone,
            updated       timestamp(6) with time zone,
            changed_by    CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            model_id      int,
            notify_change CHAR(1),
            tag           text[] COLLATE pg_catalog."default",
            encrypt_meta  boolean,
            encrypt_txt   boolean,
            style         JSONB
        );

        CREATE OR REPLACE FUNCTION ox_change_item_type() RETURNS TRIGGER AS
        $item_type_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO item_type_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO item_type_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO item_type_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $item_type_change$ LANGUAGE plpgsql;

        CREATE TRIGGER item_type_change
          AFTER INSERT OR UPDATE OR DELETE
          ON item_type
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_item_type();

        ALTER TABLE item_type_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- ITEM
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'item')
      THEN
        CREATE SEQUENCE item_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE item_id_seq
          OWNER TO onix;

        CREATE TABLE item
        (
          id           bigint                                              NOT NULL DEFAULT nextval('item_id_seq'::regclass),
          key          character varying(100) COLLATE pg_catalog."default" NOT NULL,
          name         character varying(200) COLLATE pg_catalog."default",
          description  text COLLATE pg_catalog."default",
          meta         jsonb,
          meta_enc     boolean,
          txt          text,
          txt_enc      boolean,
          enc_key_ix   smallint,
          tag          text[] COLLATE pg_catalog."default",
          attribute    hstore,
          status       smallint                                                     DEFAULT 0,
          item_type_id integer,
          version      bigint                                              NOT NULL DEFAULT 1,
          created      timestamp(6) with time zone                                  DEFAULT CURRENT_TIMESTAMP(6),
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
          partition_id bigint                                              NOT NULL DEFAULT 1,
          CONSTRAINT item_id_pk PRIMARY KEY (id),
          CONSTRAINT item_key_uc UNIQUE (key),
          CONSTRAINT item_item_type_id_fk FOREIGN KEY (item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT item_partition_id_fk FOREIGN KEY (partition_id)
            REFERENCES partition (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE item
          OWNER to onix;

        CREATE UNIQUE INDEX item_id_uix
          ON item
            (id)
          TABLESPACE pg_default;

        CREATE INDEX item_tag_ix
          ON item USING gin
            (tag COLLATE pg_catalog."default")
          TABLESPACE pg_default;

        CREATE INDEX item_attribute_ix
          ON item USING gin
            (attribute)
          TABLESPACE pg_default;

        CREATE INDEX fki_item_item_type_id_fk
          ON item USING btree (item_type_id)
          TABLESPACE pg_default;
      END IF;

      ---------------------------------------------------------------------------
      -- ITEM CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'item_change')
      THEN
        CREATE TABLE item_change
        (
          operation    CHAR(1)                     NOT NULL,
          changed      timestamp(6) with time zone NOT NULL,
          id           bigint,
          key          CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name         CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description  text COLLATE pg_catalog."default",
          meta         jsonb,
          meta_enc     boolean,
          txt          text,
          txt_enc      boolean,
          enc_key_ix   smallint,
          tag          text[] COLLATE pg_catalog."default",
          attribute    hstore,
          status       SMALLINT,
          item_type_id INTEGER,
          version      bigint,
          created      timestamp(6) with time zone,
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)      NOT NULL COLLATE pg_catalog."default",
          partition_id bigint
        );

        CREATE OR REPLACE FUNCTION ox_change_item() RETURNS TRIGGER AS
        $item_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO item_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO item_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO item_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $item_change$ LANGUAGE plpgsql;

        CREATE TRIGGER item_change
          AFTER INSERT OR UPDATE OR DELETE
          ON item
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_item();

        ALTER TABLE item_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- LINK TYPE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link_type')
      THEN
        CREATE SEQUENCE link_type_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE link_type_id_seq
          OWNER TO onix;

        CREATE TABLE link_type
        (
            id           INTEGER                NOT NULL DEFAULT nextval('link_type_id_seq'::regclass),
            key          CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            name         CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description  TEXT COLLATE pg_catalog."default",
            meta_schema  jsonb,
            tag          TEXT[] COLLATE pg_catalog."default",
            encrypt_meta BOOLEAN                NOT NULL DEFAULT FALSE,
            encrypt_txt  BOOLEAN                NOT NULL DEFAULT FALSE,
            style        JSONB,
            version      bigint                 NOT NULL DEFAULT 1,
            created      timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
            updated      timestamp(6) with time zone,
            changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            model_id     int                    NOT NULL,
            CONSTRAINT link_type_id_pk PRIMARY KEY (id),
            CONSTRAINT link_type_key_uc UNIQUE (key),
            CONSTRAINT link_type_name_uc UNIQUE (name),
            CONSTRAINT link_type_model_id_fk FOREIGN KEY (model_id)
                REFERENCES model (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE CASCADE
        )
          WITH (OIDS = FALSE)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_type_model_id_fk
          ON link_type USING btree (model_id)
          TABLESPACE pg_default;

        CREATE INDEX link_type_tag_ix
          ON link_type USING gin (tag COLLATE pg_catalog."default")
          TABLESPACE pg_default;

        ALTER TABLE link_type
          OWNER to onix;

      END IF;

      ---------------------------------------------------------------------------
      -- LINK_TYPE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link_type_change')
      THEN
        CREATE TABLE link_type_change
        (
            operation    CHAR(1)                NOT NULL,
            changed      TIMESTAMP              NOT NULL,
            id           INTEGER,
            key          CHARACTER VARYING(100) COLLATE pg_catalog."default",
            name         CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description  TEXT COLLATE pg_catalog."default",
            meta_schema  JSONB,
            tag          TEXT[],
            encrypt_meta BOOLEAN,
            encrypt_txt  BOOLEAN,
            style        JSONB,
            version      BIGINT,
            created      TIMESTAMP(6) WITH TIME ZONE,
            updated      TIMESTAMP(6) WITH TIME ZONE,
            changed_by   CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            model_id     INT
        );

        CREATE OR REPLACE FUNCTION ox_change_link_type() RETURNS TRIGGER AS
        $link_type_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO link_type_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO link_type_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO link_type_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $link_type_change$ LANGUAGE plpgsql;

        CREATE TRIGGER link_type_change
          AFTER INSERT OR UPDATE OR DELETE
          ON link_type
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_link_type();

        ALTER TABLE link_type_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- LINK
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link')
      THEN
        CREATE SEQUENCE link_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE link_id_seq
          OWNER TO onix;

        CREATE TABLE link
        (
          id            bigint                                              NOT NULL DEFAULT nextval('link_id_seq'::regclass),
          key           CHARACTER VARYING(200) COLLATE pg_catalog."default" NOT NULL,
          link_type_id  integer,
          start_item_id bigint                                              NOT NULL,
          end_item_id   bigint                                              NOT NULL,
          description   text COLLATE pg_catalog."default",
          meta          jsonb,
          meta_enc      boolean NOT NULL DEFAULT FALSE,
          txt           text,
          txt_enc       boolean NOT NULL DEFAULT FALSE,
          enc_key_ix    smallint,
          tag           text[] COLLATE pg_catalog."default",
          attribute     hstore,
          version       bigint                                              NOT NULL DEFAULT 1,
          created       TIMESTAMP(6) WITH TIME ZONE                         DEFAULT CURRENT_TIMESTAMP(6),
          updated       timestamp(6) WITH TIME ZONE,
          changed_by    CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT link_id_pk PRIMARY KEY (id),
          CONSTRAINT link_key_uc UNIQUE (key),
          CONSTRAINT link_link_type_id_fk FOREIGN KEY (link_type_id)
            REFERENCES link_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT link_end_item_id_fk FOREIGN KEY (end_item_id)
            REFERENCES item (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT link_start_item_id_fk FOREIGN KEY (start_item_id)
            REFERENCES item (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE link
          OWNER to onix;

        CREATE INDEX fki_link_link_type_id_fk
          ON link USING btree (link_type_id)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_start_item_id_fk
          ON link USING btree
            (start_item_id)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_end_item_id_fk
          ON link USING btree
            (end_item_id)
          TABLESPACE pg_default;

        CREATE INDEX link_tag_ix
          ON link USING gin
            (tag COLLATE pg_catalog."default")
          TABLESPACE pg_default;

        CREATE INDEX link_attribute_ix
          ON link USING gin
            (attribute)
          TABLESPACE pg_default;

      END IF;

      ---------------------------------------------------------------------------
      -- LINK CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link_change')
      THEN
        CREATE TABLE link_change
        (
          operation     CHAR(1)                     NOT NULL,
          changed       timestamp(6) with time zone NOT NULL,
          id            bigint,
          key           CHARACTER VARYING(200) COLLATE pg_catalog."default",
          link_type_id  integer,
          start_item_id bigint,
          end_item_id   bigint,
          description   text COLLATE pg_catalog."default",
          meta          jsonb,
          meta_enc      boolean,
          txt           text,
          txt_enc       boolean,
          enc_key_ix    smallint,
          tag           text[] COLLATE pg_catalog."default",
          attribute     hstore,
          version       bigint,
          created       TIMESTAMP(6) with time zone,
          updated       TIMESTAMP(6) with time zone,
          changed_by    CHARACTER VARYING(100)      NOT NULL COLLATE pg_catalog."default"
        );

        CREATE OR REPLACE FUNCTION ox_change_link() RETURNS TRIGGER AS
        $link_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO link_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO link_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO link_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $link_change$ LANGUAGE plpgsql;

        CREATE TRIGGER link_change
          AFTER INSERT OR UPDATE OR DELETE
          ON link
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_link();

        ALTER TABLE link_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- LINK_RULE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link_rule')
      THEN
        CREATE SEQUENCE link_rule_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE link_rule_id_seq
          OWNER TO onix;

        CREATE TABLE link_rule
        (
          id                 bigint                                              NOT NULL DEFAULT nextval('link_rule_id_seq'::regclass),
          key                character varying(300) COLLATE pg_catalog."default" NOT NULL,
          name               character varying(200) COLLATE pg_catalog."default",
          description        text COLLATE pg_catalog."default",
          link_type_id       integer                                             NOT NULL,
          start_item_type_id integer                                             NOT NULL,
          end_item_type_id   integer                                             NOT NULL,
          version            bigint                                              NOT NULL DEFAULT 1,
          created            timestamp(6) with time zone                                  DEFAULT CURRENT_TIMESTAMP(6),
          updated            timestamp(6) with time zone,
          changed_by         CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT link_rule_id_pk PRIMARY KEY (id),
          CONSTRAINT link_rule_key_uc UNIQUE (key),
          CONSTRAINT link_rule_start_item_type_id_fk FOREIGN KEY (start_item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT link_rule_end_item_type_id_fk FOREIGN KEY (end_item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
          CONSTRAINT link_rule_link_type_id_fk FOREIGN KEY (link_type_id)
            REFERENCES link_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (OIDS = FALSE)
          TABLESPACE pg_default;

        ALTER TABLE link_rule
          OWNER to onix;

        CREATE UNIQUE INDEX link_rule_id_uix
          ON link_rule
            (id)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_rule_link_type_id_fk
          ON link_rule USING btree (link_type_id)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_rule_start_item_type_id_fk
          ON link_rule USING btree (start_item_type_id)
          TABLESPACE pg_default;

        CREATE INDEX fki_link_rule_end_item_type_id_fk
          ON link_rule USING btree (end_item_type_id)
          TABLESPACE pg_default;

      END IF;

      ---------------------------------------------------------------------------
      -- LINK_RULE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'link_rule_change')
      THEN
        CREATE TABLE link_rule_change
        (
          operation          CHAR(1),
          changed            timestamp(6) with time zone,
          id                 bigint,
          key                character varying(300),
          name               character varying(200),
          description        text,
          link_type_id       integer,
          start_item_type_id integer,
          end_item_type_id   integer,
          version            bigint,
          created            timestamp(6) with time zone,
          updated            timestamp(6) with time zone,
          changed_by         CHARACTER VARYING(100)
        );

        CREATE OR REPLACE FUNCTION ox_change_link_rule() RETURNS TRIGGER AS
        $link_rule_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO link_rule_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO link_rule_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO link_rule_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $link_rule_change$
          LANGUAGE plpgsql;

        CREATE TRIGGER link_rule_change
          AFTER INSERT OR UPDATE OR DELETE
          ON link_rule
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_link_rule();

        ALTER TABLE link_rule_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- TAG
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'tag')
      THEN
        CREATE SEQUENCE tag_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE tag_id_seq
          OWNER TO onix;

        CREATE TABLE tag
        (
          id            INTEGER                NOT NULL DEFAULT nextval('tag_id_seq'::regclass),
          label         CHARACTER VARYING(50)  NOT NULL COLLATE pg_catalog."default",
          root_item_key CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name          CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description   TEXT COLLATE pg_catalog."default",
          item_data     HSTORE,
          link_data     HSTORE,
          version       BIGINT,
          created       timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated       timestamp(6) with time zone,
          changed_by    CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT tag_id_pk PRIMARY KEY (id),
          CONSTRAINT label_root_item_key_uc UNIQUE (label, root_item_key),
          CONSTRAINT root_item_key_item_data_link_data_uc UNIQUE (root_item_key, item_data, link_data),
          CONSTRAINT tag_root_item_key_fk FOREIGN KEY (root_item_key)
            REFERENCES item (key) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
        )
          WITH (
            OIDS = FALSE
          )
          TABLESPACE pg_default;

        ALTER TABLE tag
          OWNER to onix;

        CREATE INDEX fki_tag_root_item_key_fk
          ON tag USING btree (root_item_key)
          TABLESPACE pg_default;
      END IF;

      ---------------------------------------------------------------------------
      -- TAG CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'tag_change')
      THEN
        CREATE TABLE tag_change
        (
          operation     CHAR(1),
          changed       timestamp(6) with time zone,
          id            INTEGER,
          label         CHARACTER VARYING(50),
          root_item_key CHARACTER VARYING(100),
          name          CHARACTER VARYING(200),
          description   TEXT,
          item_data     HSTORE,
          link_data     HSTORE,
          version       BIGINT,
          created       timestamp(6) with time zone,
          updated       timestamp(6) with time zone,
          changed_by    CHARACTER VARYING(100)
        );

        CREATE OR REPLACE FUNCTION ox_change_tag() RETURNS TRIGGER AS
        $tag_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO tag_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO tag_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO tag_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $tag_change$
          LANGUAGE plpgsql;

        CREATE TRIGGER tag_change
          AFTER INSERT OR UPDATE OR DELETE
          ON tag
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_tag();

        ALTER TABLE tag_change
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- TYPE ATTRIBUTE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'type_attribute')
      THEN
        CREATE SEQUENCE type_attribute_id_seq
          INCREMENT 1
          START 1
          MINVALUE 1
          MAXVALUE 9223372036854775807
          CACHE 1;

        ALTER SEQUENCE type_attribute_id_seq
          OWNER TO onix;

      CREATE TABLE type_attribute
      (
        id          INTEGER NOT NULL DEFAULT nextval('type_attribute_id_seq'::regclass), -- a surrogate key for referential integrity
        key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default", -- a natural key for managing CRUD operations
        name        CHARACTER VARYING(200) COLLATE pg_catalog."default", -- the name of the attribute
        description TEXT COLLATE pg_catalog."default", -- an explanation of the attribute for clients to see
        type        CHARACTER VARYING(100) NOT NULL, -- is this a number, string, etc?
        def_value   CHARACTER VARYING(200), -- zero or more default values separated by commas
        required    BOOLEAN NOT NULL DEFAULT FALSE, -- is this a required attribute?
        regex       VARCHAR(300), -- tell client how to validate value
        item_type_id INTEGER NULL, -- the item type this attribute belongs to
        link_type_id INTEGER NULL, -- the link type this attribute belongs to
        version     bigint                 NOT NULL DEFAULT 1,
        created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
        updated     timestamp(6) with time zone,
        changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
        CONSTRAINT item_type_attribute_id_fk FOREIGN KEY (item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE,
        CONSTRAINT link_type_attribute_id_fk FOREIGN KEY (link_type_id)
            REFERENCES link_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE CASCADE
      )
      WITH (OIDS = FALSE) TABLESPACE pg_default;

      CREATE INDEX fki_type_attribute_item_type_id_fk
          ON type_attribute USING btree (item_type_id)
          TABLESPACE pg_default;

      CREATE INDEX fki_type_attribute_link_type_id_fk
          ON type_attribute USING btree (link_type_id)
          TABLESPACE pg_default;

      ALTER TABLE type_attribute
          OWNER to onix;
      END IF;

      ---------------------------------------------------------------------------
      -- TYPE ATTRIBUTE CHANGE
      ---------------------------------------------------------------------------
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = 'type_attribute_change')
      THEN
        CREATE TABLE type_attribute_change
        (
          operation   CHAR(1)                     NOT NULL,
          changed     timestamp(6) with time zone NOT NULL,
          id          INTEGER NOT NULL, -- a surrogate key for referential integrity
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default", -- a natural key for managing CRUD operations
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default", -- the name of the attribute
          description TEXT COLLATE pg_catalog."default", -- an explanation of the attribute for clients to see
          type        CHARACTER VARYING(100) NOT NULL, -- is this a number, string, etc?
          def_value   CHARACTER VARYING(300), -- zero or more default values separated by commas
          required    BOOLEAN NOT NULL DEFAULT FALSE, -- is this a required attribute?
          regex       VARCHAR(300), -- tell client how to validate value
          item_type_id INTEGER NULL, -- the item type this attribute belongs to
          link_type_id INTEGER NULL, -- the link type this attribute belongs to
          version     bigint                 NOT NULL DEFAULT 1,
          created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        CREATE OR REPLACE FUNCTION ox_change_type_attribute() RETURNS TRIGGER AS
        $type_attribute_change$
        BEGIN
          IF (TG_OP = 'DELETE') THEN
            INSERT INTO type_attribute_change SELECT 'D', now(), OLD.*;
            RETURN OLD;
          ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO type_attribute_change SELECT 'U', now(), NEW.*;
            RETURN NEW;
          ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO type_attribute_change SELECT 'I', now(), NEW.*;
            RETURN NEW;
          END IF;
          RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $type_attribute_change$ LANGUAGE plpgsql;

        CREATE TRIGGER type_attribute_change
          AFTER INSERT OR UPDATE OR DELETE
          ON type_attribute
          FOR EACH ROW
        EXECUTE PROCEDURE ox_change_type_attribute();

        ALTER TABLE type_attribute_change
          OWNER to onix;
      END IF;

    END;
    $$
