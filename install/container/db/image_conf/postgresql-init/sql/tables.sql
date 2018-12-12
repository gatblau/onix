/*
    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
    -- ITEM TYPE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item_type')
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
            id INTEGER NOT NULL DEFAULT nextval('item_type_id_seq'::regclass),
            key CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            name CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description TEXT COLLATE pg_catalog."default",
            attr_valid HSTORE,
            system boolean DEFAULT FALSE,
            version bigint NOT NULL DEFAULT 1,
            created timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
            updated timestamp(6) with time zone,
            changedby CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default",
            CONSTRAINT item_type_id_pk PRIMARY KEY (id),
            CONSTRAINT item_type_key_uc UNIQUE (key),
            CONSTRAINT item_type_name_uc UNIQUE (name)
        )
        WITH (
            OIDS = FALSE
        )
        TABLESPACE pg_default;

        ALTER TABLE item_type
            OWNER to onix;

        INSERT INTO item_type(key, name, description, system, changedby) VALUES ('INVENTORY', 'Ansible Inventory', 'An Ansible inventory.', TRUE, 'onix');
        INSERT INTO item_type(key, name, description, system, changedby) VALUES ('HOST-GROUP', 'Host Group', 'An Ansible host group.', TRUE, 'onix');
        INSERT INTO item_type(key, name, description, system, changedby) VALUES ('HOST', 'Host', 'An Operating System Host.', TRUE, 'onix');
        INSERT INTO item_type(key, name, description, system, changedby) VALUES ('LICENCE', 'A software licence.', 'Describes the information pertaining to a software licence.', TRUE, 'onix');

	  END IF;

    ---------------------------------------------------------------------------
    -- ITEM
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item')
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
            id bigint NOT NULL DEFAULT nextval('item_id_seq'::regclass),
            key character varying(100) COLLATE pg_catalog."default" NOT NULL,
            name character varying(200) COLLATE pg_catalog."default",
            description text COLLATE pg_catalog."default",
            meta jsonb,
            tag text[] COLLATE pg_catalog."default",
            attribute hstore,
            status smallint DEFAULT 0,
            item_type_id integer,
            version bigint NOT NULL DEFAULT 1,
            created timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
            updated timestamp(6) with time zone,
            changedby CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            CONSTRAINT item_id_pk PRIMARY KEY (id),
            CONSTRAINT item_key_uc UNIQUE (key),
            CONSTRAINT item_item_type_id_fk FOREIGN KEY (item_type_id)
                REFERENCES item_type (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE NO ACTION
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
    -- LINK TYPE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_type')
    THEN
        CREATE SEQUENCE link_type_id_seq
        INCREMENT 1
        START 1
        MINVALUE 1
        MAXVALUE 9223372036854775807
        CACHE 1;

        ALTER SEQUENCE link_type_id_seq
        OWNER TO onix;

        CREATE TABLE link_type (
           id INTEGER NOT NULL DEFAULT nextval('link_type_id_seq'::regclass),
           key CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
           name CHARACTER VARYING(200) COLLATE pg_catalog."default",
           description TEXT COLLATE pg_catalog."default",
           attr_valid HSTORE,
           system boolean DEFAULT FALSE,
           version bigint NOT NULL DEFAULT 1,
           created timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
           updated timestamp(6) with time zone,
           changedby CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default",
           CONSTRAINT link_type_id_pk PRIMARY KEY (id),
           CONSTRAINT link_type_key_uc UNIQUE (key),
           CONSTRAINT link_type_name_uc UNIQUE (name)
        )
        WITH (OIDS = FALSE) TABLESPACE pg_default;

        ALTER TABLE link_type OWNER to onix;

        INSERT INTO link_type(key, name, description, system, changedby) VALUES ('INVENTORY', 'Inventory Link Type.', 'Links items describing an inventory.', TRUE, 'onix');
        INSERT INTO link_type(key, name, description, system, changedby) VALUES ('LICENSE', 'Licence Link Type.', 'Links items related by a licence.', TRUE, 'onix');
    END IF;

    ---------------------------------------------------------------------------
    -- LINK
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link')
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
            id bigint NOT NULL DEFAULT nextval('link_id_seq'::regclass),
            key CHARACTER VARYING(200) COLLATE pg_catalog."default" NOT NULL,
            link_type_id integer,
            start_item_id bigint NOT NULL,
            end_item_id bigint NOT NULL,
            description text COLLATE pg_catalog."default",
            meta jsonb,
            tag text[] COLLATE pg_catalog."default",
            attribute hstore,
            version bigint NOT NULL DEFAULT 1,
            created TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6),
            updated timestamp(6) WITH TIME ZONE,
            changedby CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
            CONSTRAINT link_id_pk PRIMARY KEY (id),
            CONSTRAINT link_key_uc UNIQUE (key),
            CONSTRAINT link_link_type_id_fk FOREIGN KEY (link_type_id)
                REFERENCES link_type (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE NO ACTION,
            CONSTRAINT link_end_item_id_fk FOREIGN KEY (end_item_id)
                REFERENCES item (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE NO ACTION,
            CONSTRAINT link_start_item_id_fk FOREIGN KEY (start_item_id)
                REFERENCES item (id) MATCH SIMPLE
                ON UPDATE NO ACTION
                ON DELETE NO ACTION
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
    -- ITEM AUDIT
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item_audit')
    THEN
        CREATE TABLE item_audit
        (
            operation CHAR(1) NOT NULL,
            change_date TIMESTAMP NOT NULL,
            id bigint,
            key CHARACTER VARYING(100) COLLATE pg_catalog."default",
            name CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description CHARACTER VARYING(500) COLLATE pg_catalog."default",
            meta jsonb,
            tag text[] COLLATE pg_catalog."default",
            attribute hstore,
            status SMALLINT,
            item_type_id INTEGER,
            version bigint,
            created TIMESTAMP(6) with time zone,
            updated TIMESTAMP(6) with time zone,
            changedby CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE item_audit
            OWNER to onix;

        CREATE OR REPLACE FUNCTION audit_item() RETURNS TRIGGER AS $item_audit$
        BEGIN
            IF (TG_OP = 'DELETE') THEN
                INSERT INTO item_audit SELECT 'D', now(), OLD.*;
                RETURN OLD;
            ELSIF (TG_OP = 'UPDATE') THEN
                INSERT INTO item_audit SELECT 'U', now(), NEW.*;
                RETURN NEW;
            ELSIF (TG_OP = 'INSERT') THEN
                INSERT INTO item_audit SELECT 'I', now(), NEW.*;
                RETURN NEW;
            END IF;
            RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $item_audit$ LANGUAGE plpgsql;

        CREATE TRIGGER item_audit
        AFTER INSERT OR UPDATE OR DELETE ON item
            FOR EACH ROW EXECUTE PROCEDURE audit_item();
    END IF;

    ---------------------------------------------------------------------------
    -- LINK AUDIT
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_audit')
    THEN
        CREATE TABLE link_audit
        (
            operation CHAR(1) NOT NULL,
            change_date TIMESTAMP NOT NULL,
            id bigint,
            key CHARACTER VARYING(200) COLLATE pg_catalog."default",
            link_type_id integer,
            start_item_id bigint,
            end_item_id bigint,
            description CHARACTER VARYING(500) COLLATE pg_catalog."default",
            meta json,
            tag text[] COLLATE pg_catalog."default",
            attribute hstore,
            version bigint,
            created TIMESTAMP(6) with time zone,
            updated TIMESTAMP(6) with time zone,
            changedby CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE link_audit
            OWNER to onix;

        CREATE OR REPLACE FUNCTION audit_link() RETURNS TRIGGER AS $link_audit$
        BEGIN
            IF (TG_OP = 'DELETE') THEN
                INSERT INTO link_audit SELECT 'D', now(), OLD.*;
                RETURN OLD;
            ELSIF (TG_OP = 'UPDATE') THEN
                INSERT INTO link_audit SELECT 'U', now(), NEW.*;
                RETURN NEW;
            ELSIF (TG_OP = 'INSERT') THEN
                INSERT INTO link_audit SELECT 'I', now(), NEW.*;
                RETURN NEW;
            END IF;
            RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $link_audit$ LANGUAGE plpgsql;

        CREATE TRIGGER link_audit
        AFTER INSERT OR UPDATE OR DELETE ON link
          FOR EACH ROW EXECUTE PROCEDURE audit_link();
    END IF;

    ---------------------------------------------------------------------------
    -- ITEM_TYPE AUDIT
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item_type_audit')
    THEN
        CREATE TABLE item_type_audit
        (
            operation CHAR(1) NOT NULL,
            changed TIMESTAMP NOT NULL,
            id INTEGER,
            key CHARACTER VARYING(100) COLLATE pg_catalog."default",
            name CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description TEXT COLLATE pg_catalog."default",
            attr_valid HSTORE,
            system boolean,
            version bigint,
            created timestamp(6) with time zone,
            updated timestamp(6) with time zone,
            changedby CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE item_type_audit
            OWNER to onix;

        CREATE OR REPLACE FUNCTION audit_item_type() RETURNS TRIGGER AS $item_type_audit$
        BEGIN
            IF (TG_OP = 'DELETE') THEN
                INSERT INTO item_type_audit SELECT 'D', now(), OLD.*;
                RETURN OLD;
            ELSIF (TG_OP = 'UPDATE') THEN
                INSERT INTO item_type_audit SELECT 'U', now(), NEW.*;
                RETURN NEW;
            ELSIF (TG_OP = 'INSERT') THEN
                INSERT INTO item_type_audit SELECT 'I', now(), NEW.*;
                RETURN NEW;
            END IF;
            RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $item_type_audit$ LANGUAGE plpgsql;

        CREATE TRIGGER item_type_audit
        AFTER INSERT OR UPDATE OR DELETE ON item_type
          FOR EACH ROW EXECUTE PROCEDURE audit_item_type();
    END IF;

    ---------------------------------------------------------------------------
    -- ITEM_TYPE AUDIT
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_type_audit')
    THEN
        CREATE TABLE link_type_audit
        (
            operation CHAR(1) NOT NULL,
            changed TIMESTAMP NOT NULL,
            id INTEGER,
            key CHARACTER VARYING(100) COLLATE pg_catalog."default",
            name CHARACTER VARYING(200) COLLATE pg_catalog."default",
            description TEXT COLLATE pg_catalog."default",
            attr_valid HSTORE,
            system boolean,
            version bigint,
            created timestamp(6) with time zone,
            updated timestamp(6) with time zone,
            changedby CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE link_type_audit
            OWNER to onix;

        CREATE OR REPLACE FUNCTION audit_link_type() RETURNS TRIGGER AS $link_type_audit$
        BEGIN
            IF (TG_OP = 'DELETE') THEN
                INSERT INTO item_type_audit SELECT 'D', now(), OLD.*;
                RETURN OLD;
            ELSIF (TG_OP = 'UPDATE') THEN
                INSERT INTO item_type_audit SELECT 'U', now(), NEW.*;
                RETURN NEW;
            ELSIF (TG_OP = 'INSERT') THEN
                INSERT INTO item_type_audit SELECT 'I', now(), NEW.*;
                RETURN NEW;
            END IF;
                RETURN NULL; -- result is ignored since this is an AFTER trigger
        END;
        $link_type_audit$ LANGUAGE plpgsql;

        CREATE TRIGGER link_type_audit
            AFTER INSERT OR UPDATE OR DELETE ON link_type
            FOR EACH ROW EXECUTE PROCEDURE audit_link_type();
    END IF;
END;
$$
