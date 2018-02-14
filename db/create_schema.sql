DROP INDEX fki_link_end_item_id_fk;
DROP INDEX fki_link_start_item_id_fk;
DROP TABLE link;
DROP INDEX fki_item_item_type_id_fk;
DROP TABLE item;
DROP TABLE item_type;

CREATE TABLE item_type
(
    id integer NOT NULL,
    name character varying(200) COLLATE pg_catalog."default",
    description character varying(500) COLLATE pg_catalog."default",
    CONSTRAINT item_type_id_pk PRIMARY KEY (id),
    CONSTRAINT item_type_name_uc UNIQUE (name)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE item_type
    OWNER to admin;

CREATE TABLE item
(
    id bigint NOT NULL DEFAULT nextval('item_id_seq'::regclass),
    name character varying(200) COLLATE pg_catalog."default",
    description character varying(500) COLLATE pg_catalog."default",
    item_type_id integer,
    meta json,
    version bigint NOT NULL DEFAULT 1,
    created timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated timestamp(6) with time zone,
    tag character varying(300) COLLATE pg_catalog."default",
    key character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT item_id_pk PRIMARY KEY (id),
    CONSTRAINT item_key_uc UNIQUE (key)﻿NOT NULL,
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
    OWNER to admin;

CREATE INDEX fki_item_item_type_id_fk
    ON item USING btree
    (item_type_id)
    TABLESPACE pg_default;

CREATE TABLE link
(
    id bigint NOT NULL DEFAULT nextval('link_id_seq'::regclass),
    meta json,
    description character varying(500) COLLATE pg_catalog."default",
    version bigint,
    created timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated timestamp(6) with time zone,
    start_item_id bigint,
    end_item_id bigint,
    role character varying(200) COLLATE pg_catalog."default" NOT NULL,
    key character varying(50) COLLATE pg_catalog."default"﻿NOT NULL,
    tag character varying(300) COLLATE pg_catalog."default",
    CONSTRAINT link_id_pk PRIMARY KEY (id),
    CONSTRAINT link_key_uc UNIQUE (key),
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

ALTER TABLE public.link
    OWNER to admin;

CREATE INDEX fki_link_end_item_id_fk
    ON link USING btree
    (end_item_id)
    TABLESPACE pg_default;

CREATE INDEX fki_link_start_item_id_fk
    ON link USING btree
    (start_item_id)
    TABLESPACE pg_default;

CREATE TABLE item_audit(
    operation char(1)   NOT NULL,
    stamp timestamp NOT NULL,
    userid text NOT NULL,
    id bigint,
    name character varying(200) COLLATE pg_catalog."default",
    description character varying(500) COLLATE pg_catalog."default",
    item_type_id integer,
    meta json,
    version bigint,
    created timestamp(6) with time zone,
    updated timestamp(6) with time zone,
    tag character varying(300) COLLATE pg_catalog."default",
    key character varying(50) COLLATE pg_catalog."default"
);

CREATE OR REPLACE FUNCTION audit_item() RETURNS TRIGGER AS $item_audit$
    BEGIN
        IF (TG_OP = 'DELETE') THEN
            INSERT INTO item_audit SELECT 'D', now(), user, OLD.*;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO item_audit SELECT 'U', now(), user, NEW.*;
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO item_audit SELECT 'I', now(), user, NEW.*;
            RETURN NEW;
        END IF;
        RETURN NULL; -- result is ignored since this is an AFTER trigger
    END;
$item_audit$ LANGUAGE plpgsql;

CREATE TRIGGER item_audit
AFTER INSERT OR UPDATE OR DELETE ON item
    FOR EACH ROW EXECUTE PROCEDURE audit_item();

CREATE TABLE link_audit(
    operation char(1) NOT NULL,
    stamp timestamp NOT NULL,
    userid text NOT NULL,
    id bigint,
    meta json,
    description character varying(500) COLLATE pg_catalog."default",
    version bigint,
    created timestamp(6) with time zone,
    updated timestamp(6) with time zone,
    start_item_id bigint,
    end_item_id bigint,
    role character varying(200) COLLATE pg_catalog."default",
    key character varying(50) COLLATE pg_catalog."default",
    tag character varying(300) COLLATE pg_catalog."default"
);

CREATE OR REPLACE FUNCTION audit_link() RETURNS TRIGGER AS $link_audit$
    BEGIN
        IF (TG_OP = 'DELETE') THEN
            INSERT INTO link_audit SELECT 'D', now(), user, OLD.*;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO link_audit SELECT 'U', now(), user, NEW.*;
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO link_audit SELECT 'I', now(), user, NEW.*;
            RETURN NEW;
        END IF;
        RETURN NULL; -- result is ignored since this is an AFTER trigger
    END;
$link_audit$ LANGUAGE plpgsql;

CREATE TRIGGER link_audit
AFTER INSERT OR UPDATE OR DELETE ON link
    FOR EACH ROW EXECUTE PROCEDURE audit_link();