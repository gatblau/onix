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
            START 1000
            MINVALUE 1000
            MAXVALUE 9223372036854775807
            CACHE 1;

        ALTER SEQUENCE item_type_id_seq
            OWNER TO onix;

        CREATE TABLE item_type
        (
          id          INTEGER                NOT NULL DEFAULT nextval('item_type_id_seq'::regclass),
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          attr_valid  HSTORE,
          filter      jsonb,
          system      boolean                         DEFAULT FALSE,
          version     bigint                 NOT NULL DEFAULT 1,
          created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(50)  NOT NULL COLLATE pg_catalog."default",
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

    END IF;

    ---------------------------------------------------------------------------
    -- ITEM_TYPE CHANGE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item_type_change')
    THEN
        CREATE TABLE item_type_change
        (
          operation   CHAR(1)               NOT NULL,
          changed     TIMESTAMP             NOT NULL,
          id          INTEGER,
          key         CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          attr_valid  HSTORE,
          filter      jsonb,
          system      boolean,
          version     bigint,
          created     timestamp(6) with time zone,
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE item_type_change
            OWNER to onix;

        CREATE OR REPLACE FUNCTION change_item_type() RETURNS TRIGGER AS $item_type_change$
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
            AFTER INSERT OR UPDATE OR DELETE ON item_type
            FOR EACH ROW EXECUTE PROCEDURE change_item_type();

        INSERT INTO item_type(id, key, name, description, system, changed_by) VALUES (50, 'ANSIBLE_INVENTORY', 'Ansible Inventory', 'An Ansible inventory.', TRUE, 'onix');
        INSERT INTO item_type(id, key, name, description, system, changed_by) VALUES (51, 'ANSIBLE_HOST_GROUP_SET', 'Host Group Set', 'An Ansible Set of Host Groups.', TRUE, 'onix');
        INSERT INTO item_type(id, key, name, description, system, changed_by) VALUES (52, 'ANSIBLE_HOST_GROUP', 'Ansible Host Group', 'An Ansible Inventory Host Group.', TRUE, 'onix');
        INSERT INTO item_type(id, key, name, description, system, changed_by) VALUES (53, 'ANSIBLE_HOST', 'Ansible Host', 'An Ansible Inventory Host', TRUE, 'onix');

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
          id           bigint                                              NOT NULL DEFAULT nextval('item_id_seq'::regclass),
          key          character varying(100) COLLATE pg_catalog."default" NOT NULL,
          name         character varying(200) COLLATE pg_catalog."default",
          description  text COLLATE pg_catalog."default",
          meta         jsonb,
          tag          text[] COLLATE pg_catalog."default",
          attribute    hstore,
          status       smallint                                                     DEFAULT 0,
          item_type_id integer,
          version      bigint                                              NOT NULL DEFAULT 1,
          created      timestamp(6) with time zone                                  DEFAULT CURRENT_TIMESTAMP(6),
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
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
    -- ITEM CHANGE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='item_change')
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
          tag          text[] COLLATE pg_catalog."default",
          attribute    hstore,
          status       SMALLINT,
          item_type_id INTEGER,
          version      bigint,
          created      timestamp(6) with time zone,
          updated      timestamp(6) with time zone,
          changed_by   CHARACTER VARYING(100)      NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE item_change
            OWNER to onix;

        CREATE OR REPLACE FUNCTION change_item() RETURNS TRIGGER AS $item_change$
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
            AFTER INSERT OR UPDATE OR DELETE ON item
            FOR EACH ROW EXECUTE PROCEDURE change_item();

    END IF;

    ---------------------------------------------------------------------------
    -- LINK TYPE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_type')
    THEN
        CREATE SEQUENCE link_type_id_seq
        INCREMENT 1
        START 1000
        MINVALUE 1000
        MAXVALUE 9223372036854775807
        CACHE 1;

        ALTER SEQUENCE link_type_id_seq
        OWNER TO onix;

        CREATE TABLE link_type
        (
          id          INTEGER                NOT NULL DEFAULT nextval('link_type_id_seq'::regclass),
          key         CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          attr_valid  HSTORE,
          system      boolean                NOT NULL DEFAULT FALSE,
          version     bigint                 NOT NULL DEFAULT 1,
          created     timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(50)  NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT link_type_id_pk PRIMARY KEY (id),
          CONSTRAINT link_type_key_uc UNIQUE (key),
          CONSTRAINT link_type_name_uc UNIQUE (name)
        )
        WITH (OIDS = FALSE) TABLESPACE pg_default;

        ALTER TABLE link_type OWNER to onix;

    END IF;

    ---------------------------------------------------------------------------
    -- LINK_TYPE CHANGE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_type_change')
    THEN
        CREATE TABLE link_type_change
        (
          operation   CHAR(1)               NOT NULL,
          changed     TIMESTAMP             NOT NULL,
          id          INTEGER,
          key         CHARACTER VARYING(100) COLLATE pg_catalog."default",
          name        CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description TEXT COLLATE pg_catalog."default",
          attr_valid  HSTORE,
          system      boolean,
          version     bigint,
          created     timestamp(6) with time zone,
          updated     timestamp(6) with time zone,
          changed_by  CHARACTER VARYING(50) NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE link_type_change
            OWNER to onix;

        CREATE OR REPLACE FUNCTION change_link_type() RETURNS TRIGGER AS $link_type_change$
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
            AFTER INSERT OR UPDATE OR DELETE ON link_type
            FOR EACH ROW EXECUTE PROCEDURE change_link_type();

        INSERT INTO link_type(id, key, name, description, system, changed_by) VALUES (1, 'APPLICATION', 'Application Link', 'Links items describing application components.', TRUE, 'onix');
        INSERT INTO link_type(id, key, name, description, system, changed_by) VALUES (2, 'NETWORK', 'Network Link', 'Links items describing network connections.', TRUE, 'onix');
        INSERT INTO link_type(id, key, name, description, system, changed_by) VALUES (3, 'WEB-CONTENT', 'Web Content Link', 'Links items describing web content.', TRUE, 'onix');
        INSERT INTO link_type(id, key, name, description, system, changed_by) VALUES (50, 'ANSIBLE_INVENTORY', 'Ansible Inventory Link', 'Links items describing an Ansible inventory.', TRUE, 'onix');

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
          id            bigint                                              NOT NULL DEFAULT nextval('link_id_seq'::regclass),
          key           CHARACTER VARYING(200) COLLATE pg_catalog."default" NOT NULL,
          link_type_id  integer,
          start_item_id bigint                                              NOT NULL,
          end_item_id   bigint                                              NOT NULL,
          description   text COLLATE pg_catalog."default",
          meta          jsonb,
          tag           text[] COLLATE pg_catalog."default",
          attribute     hstore,
          version       bigint                                              NOT NULL DEFAULT 1,
          created       TIMESTAMP(6) WITH TIME ZONE                                  DEFAULT CURRENT_TIMESTAMP(6),
          updated       timestamp(6) WITH TIME ZONE,
          changed_by    CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
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
    -- LINK CHANGE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_change')
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
          tag           text[] COLLATE pg_catalog."default",
          attribute     hstore,
          version       bigint,
          created       TIMESTAMP(6) with time zone,
          updated       TIMESTAMP(6) with time zone,
          changed_by    CHARACTER VARYING(100)      NOT NULL COLLATE pg_catalog."default"
        );

        ALTER TABLE link_change
            OWNER to onix;

        CREATE OR REPLACE FUNCTION change_link() RETURNS TRIGGER AS $link_change$
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
            AFTER INSERT OR UPDATE OR DELETE ON link
            FOR EACH ROW EXECUTE PROCEDURE change_link();

    END IF;

    ---------------------------------------------------------------------------
    -- LINK_RULE
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_rule')
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
          system             boolean                                             NOT NULL DEFAULT FALSE,
          version            bigint                                              NOT NULL DEFAULT 1,
          created            timestamp(6) with time zone                                  DEFAULT CURRENT_TIMESTAMP(6),
          updated            timestamp(6) with time zone,
          changed_by         CHARACTER VARYING(100)                              NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT link_rule_id_pk PRIMARY KEY (id),
          CONSTRAINT link_rule_key_uc UNIQUE (key),
          CONSTRAINT link_rule_start_item_type_id_fk FOREIGN KEY (start_item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION,
          CONSTRAINT link_rule_end_item_type_id_fk FOREIGN KEY (end_item_type_id)
            REFERENCES item_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION,
          CONSTRAINT link_rule_link_type_id_fk FOREIGN KEY (link_type_id)
            REFERENCES link_type (id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
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
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='link_rule_change')
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
          system             boolean,
          version            bigint,
          created            timestamp(6) with time zone,
          updated            timestamp(6) with time zone,
          changed_by         CHARACTER VARYING(100)
        );

        ALTER TABLE link_rule_change
            OWNER to onix;

        CREATE OR REPLACE FUNCTION change_link_rule() RETURNS TRIGGER AS $link_rule_change$
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
            AFTER INSERT OR UPDATE OR DELETE ON link_rule
            FOR EACH ROW EXECUTE PROCEDURE change_link_rule();

        INSERT INTO link_rule (id, key, name, description, link_type_id, start_item_type_id, end_item_type_id, changed_by, system) VALUES (1, 'ANSIBLE-INVENTORY->ANSIBLE-HOST-GROUP-SET', 'Inventory to Group of Host Groups link rule.', 'Allows to link an inventory with a group of host groups.', 50, 50, 51, 'onix', TRUE);
        INSERT INTO link_rule (id, key, name, description, link_type_id, start_item_type_id, end_item_type_id, changed_by, system) VALUES (2, 'ANSIBLE-INVENTORY->ANSIBLE-HOST-GROUP', 'Inventory to ANSIBLE-HOST-GROUP link rule.', 'Allows to link an inventory item with a host group item.', 50, 50, 52, 'onix', TRUE);
        INSERT INTO link_rule (id, key, name, description, link_type_id, start_item_type_id, end_item_type_id, changed_by, system) VALUES (3, 'ANSIBLE-HOST-GROUP-SET->ANSIBLE-HOST-GROUP', 'Group of Host Groups to Groups link rule.', 'Allows to link a group of host groups with a host group.', 50, 51, 52, 'onix', TRUE);
        INSERT INTO link_rule (id, key, name, description, link_type_id, start_item_type_id, end_item_type_id, changed_by, system) VALUES (4, 'ANSIBLE-HOST-GROUP->HOST', 'Host Group to Host link rule.', 'Allows to link a host group item with a host item.', 50, 52, 53, 'onix', TRUE);

    END IF;

    ---------------------------------------------------------------------------
    -- SNAPSHOT
    ---------------------------------------------------------------------------
    IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='snapshot')
    THEN
        CREATE SEQUENCE snapshot_id_seq
        INCREMENT 1
        START 1
        MINVALUE 1
        MAXVALUE 9223372036854775807
        CACHE 1;

        ALTER SEQUENCE snapshot_id_seq
            OWNER TO onix;

        CREATE TABLE snapshot
        (
          id            INTEGER                NOT NULL DEFAULT nextval('snapshot_id_seq'::regclass),
          label         CHARACTER VARYING(50)  NOT NULL COLLATE pg_catalog."default",
          root_item_key CHARACTER VARYING(100) NOT NULL COLLATE pg_catalog."default",
          name          CHARACTER VARYING(200) COLLATE pg_catalog."default",
          description   TEXT COLLATE pg_catalog."default",
          item_data     HSTORE,
          link_data     HSTORE,
          version       BIGINT,
          created       timestamp(6) with time zone     DEFAULT CURRENT_TIMESTAMP(6),
          updated       timestamp(6) with time zone,
          changed_by    CHARACTER VARYING(50)  NOT NULL COLLATE pg_catalog."default",
          CONSTRAINT snapshot_id_pk PRIMARY KEY (id),
          CONSTRAINT label_root_item_key_uc UNIQUE (label, root_item_key),
          CONSTRAINT root_item_key_item_data_link_data_uc UNIQUE (root_item_key, item_data, link_data),
          CONSTRAINT snapshot_root_item_key_fk FOREIGN KEY (root_item_key)
            REFERENCES item (key) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
        )
        WITH (
            OIDS = FALSE
        )
        TABLESPACE pg_default;

        ALTER TABLE snapshot
            OWNER to onix;

        CREATE INDEX fki_snapshot_root_item_key_fk
            ON snapshot USING btree (root_item_key)
            TABLESPACE pg_default;
    END IF;

  ---------------------------------------------------------------------------
  -- SNAPSHOT CHANGE
  ---------------------------------------------------------------------------
  IF NOT EXISTS (SELECT relname FROM pg_class WHERE relname='snapshot_change')
  THEN
    CREATE TABLE snapshot_change
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
      changed_by    CHARACTER VARYING(50)
    );

    ALTER TABLE snapshot_change
      OWNER to onix;

    CREATE OR REPLACE FUNCTION change_snapshot() RETURNS TRIGGER AS $snapshot_change$
    BEGIN
      IF (TG_OP = 'DELETE') THEN
        INSERT INTO snapshot_change SELECT 'D', now(), OLD.*;
        RETURN OLD;
      ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO snapshot_change SELECT 'U', now(), NEW.*;
        RETURN NEW;
      ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO snapshot_change SELECT 'I', now(), NEW.*;
        RETURN NEW;
      END IF;
      RETURN NULL; -- result is ignored since this is an AFTER trigger
      END;
    $snapshot_change$
    LANGUAGE plpgsql;

    CREATE TRIGGER snapshot_change
      AFTER INSERT OR UPDATE OR DELETE ON snapshot
      FOR EACH ROW EXECUTE PROCEDURE change_snapshot();

END IF;

END;
$$
